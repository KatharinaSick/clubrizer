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

async function getBrowserSubscription(): Promise<PushSubscription | null> {
  const registration = await navigator.serviceWorker.ready
  return registration.pushManager.getSubscription()
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
  // user (turned off from another device, or someone else signed in here since).
  return (await notificationsService.isSubscribed(subscription.endpoint)) ? 'on' : 'off'
}

// turnOnPush asks for permission and registers this device. Must be called directly from a
// tap/click — browsers (especially Safari) only show the permission prompt for user actions.
export async function turnOnPush(): Promise<PushState> {
  const permission = await Notification.requestPermission()
  if (permission === 'denied') {
    return 'blocked'
  }
  if (permission !== 'granted') {
    return 'off'
  }

  const registration = await navigator.serviceWorker.ready
  let subscription = await registration.pushManager.getSubscription()
  if (!subscription) {
    const publicKey = await notificationsService.getVapidPublicKey()
    subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: base64UrlToUint8Array(publicKey),
    })
  }

  await notificationsService.subscribe(subscription.toJSON() as PushSubscriptionPayload)
  return 'on'
}

export async function turnOffPushOnThisDevice(): Promise<void> {
  const subscription = await getBrowserSubscription()
  if (!subscription) {
    return
  }
  await notificationsService.unsubscribe(subscription.endpoint)
  await subscription.unsubscribe()
}

export async function turnOffPushOnAllDevices(): Promise<void> {
  await notificationsService.unsubscribeAll()
  const subscription = await getBrowserSubscription()
  await subscription?.unsubscribe()
}

// unregisterPushOnLogout removes this device before signing out, so the next person using
// it doesn't get the previous user's notifications. Never blocks the logout.
export async function unregisterPushOnLogout(): Promise<void> {
  if (!isSupported() || Notification.permission !== 'granted') {
    return
  }
  try {
    await turnOffPushOnThisDevice()
  } catch {
    // ignore — logging out matters more
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
