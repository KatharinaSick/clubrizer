import notificationsService, { type PushSubscriptionPayload } from '@/service/notifications'

// What this device can do with push notifications right now:
// - unsupported:     the browser can't receive push notifications at all
// - needsHomeScreen: iPhone/iPad in a Safari tab — push only works once the app is added to
//                    the home screen and opened from there
// - blocked:         the user denied the permission; only the browser/phone settings can undo it
// - off / on:        whether this device is registered for the signed-in user
export type PushState = 'unsupported' | 'needsHomeScreen' | 'blocked' | 'off' | 'on'

export function isIOS(): boolean {
  // iPadOS reports itself as a Mac, but Macs have no touch screen.
  return /iPad|iPhone|iPod/.test(navigator.userAgent) || (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
}

export function isStandalone(): boolean {
  return window.matchMedia('(display-mode: standalone)').matches
    || (navigator as Navigator & { standalone?: boolean }).standalone === true
}

function isSupported(): boolean {
  return 'serviceWorker' in navigator && 'PushManager' in window && 'Notification' in window
}

// getBrowserSubscription returns this browser's subscription, or null if there is none or no
// service worker is registered. Unlike navigator.serviceWorker.ready it never waits forever.
async function getBrowserSubscription(): Promise<PushSubscription | null> {
  const registration = await navigator.serviceWorker.getRegistration()
  return (await registration?.pushManager.getSubscription()) ?? null
}

// readyRegistration waits for the active service worker, which subscribing needs. ready never
// settles if the worker failed to register, so give up after a while instead of hanging.
const SERVICE_WORKER_TIMEOUT_MS = 10_000

function readyRegistration(): Promise<ServiceWorkerRegistration> {
  return Promise.race([
    navigator.serviceWorker.ready,
    new Promise<never>((_, reject) =>
      setTimeout(() => reject(new Error('service worker not ready')), SERVICE_WORKER_TIMEOUT_MS),
    ),
  ])
}

export async function getPushState(): Promise<PushState> {
  if (isIOS() && !isStandalone()) {
    return 'needsHomeScreen'
  }
  if (!isSupported()) {
    return 'unsupported'
  }
  if (Notification.permission === 'denied') {
    return 'blocked'
  }
  if (Notification.permission !== 'granted') {
    return 'off'
  }

  const subscription = await getBrowserSubscription()
  if (!subscription) {
    return 'off'
  }
  // The browser alone doesn't know whether the backend still has this device for the current
  // user (e.g. someone else signed in here since).
  return (await notificationsService.isSubscribed(subscription.endpoint)) ? 'on' : 'off'
}

// turnOnPush asks for permission and registers this device. Must be called directly from a
// tap/click — browsers (especially Safari) only show the permission prompt for user actions.
// The email of whoever turned notifications on in this browser. Kept across logouts, so the
// same person signing in again gets their notifications back without turning them on again,
// while a different person signing in releases the device (see restorePushOnSignIn).
const OWNER_KEY = 'push.owner'

function getOwner(): string | null {
  try {
    return localStorage.getItem(OWNER_KEY)
  } catch {
    return null
  }
}

function setOwner(email: string | null) {
  try {
    if (email === null) {
      localStorage.removeItem(OWNER_KEY)
    } else {
      localStorage.setItem(OWNER_KEY, email)
    }
  } catch {
    // storage unavailable — notifications then just need turning on again after a logout
  }
}

export async function turnOnPush(ownerEmail: string): Promise<PushState> {
  const permission = await Notification.requestPermission()
  if (permission === 'denied') {
    return 'blocked'
  }
  if (permission !== 'granted') {
    return 'off'
  }

  const registration = await readyRegistration()
  let subscription = await registration.pushManager.getSubscription()
  if (!subscription) {
    const publicKey = await notificationsService.getVapidPublicKey()
    subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: base64UrlToUint8Array(publicKey),
    })
  }

  await notificationsService.subscribe(subscription.toJSON() as PushSubscriptionPayload)
  setOwner(ownerEmail)
  return 'on'
}

export async function turnOffPushOnThisDevice(): Promise<void> {
  const subscription = await getBrowserSubscription()
  if (!subscription) {
    return
  }
  await notificationsService.unsubscribe(subscription.endpoint)
  await subscription.unsubscribe()
  setOwner(null)
}

// releasePushOnLogout stops notifications to this device when signing out, so nothing arrives
// here while nobody is signed in. Works even if the session has already expired (e.g. a failed
// token refresh). The browser keeps its subscription and owner, so the same person gets
// notifications back on their next sign-in (see restorePushOnSignIn). Never blocks the logout.
export async function releasePushOnLogout(): Promise<void> {
  if (!isSupported()) {
    return
  }
  try {
    const subscription = await getBrowserSubscription()
    if (subscription) {
      const json = subscription.toJSON() as PushSubscriptionPayload
      await notificationsService.releaseDevice(json.endpoint, json.keys.auth)
    }
  } catch {
    // ignore — logging out matters more
  }
}

// restorePushOnSignIn runs after every sign-in. If the same person turned notifications on in
// this browser before, they are quietly registered again. If it was someone else, their
// registration is removed from this device — even if their session expired without logging
// out — so their notifications stop arriving here; the new person can turn notifications on
// for themselves. Best-effort: never blocks or fails the sign-in.
export async function restorePushOnSignIn(email: string): Promise<void> {
  if (!isSupported() || Notification.permission !== 'granted') {
    return
  }
  try {
    const subscription = await getBrowserSubscription()
    if (!subscription) {
      return
    }

    const json = subscription.toJSON() as PushSubscriptionPayload
    if (getOwner() === email) {
      await notificationsService.subscribe(json, true)
      return
    }

    await notificationsService.releaseDevice(json.endpoint, json.keys.auth)
    await subscription.unsubscribe()
    setOwner(null)
  } catch {
    // ignore — notifications can always be turned on again on the notifications screen
  }
}

function base64UrlToUint8Array(base64Url: string): Uint8Array<ArrayBuffer> {
  const padding = '='.repeat((4 - (base64Url.length % 4)) % 4)
  const base64 = (base64Url + padding).replace(/-/g, '+').replace(/_/g, '/')
  const raw = atob(base64)
  const bytes = new Uint8Array(raw.length)
  for (let i = 0; i < raw.length; i++) {
    bytes[i] = raw.charCodeAt(i)
  }
  return bytes
}
