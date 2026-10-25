/// <reference lib="webworker" />
import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { clientsClaim } from 'workbox-core'
import i18n from './plugins/i18n'

declare let self: ServiceWorkerGlobalScope

// Same behaviour as the previously generated service worker: activate updates right away
// (registerType 'autoUpdate'), precache the build and serve the app for every navigation.
self.skipWaiting()
clientsClaim()
cleanupOutdatedCaches()
precacheAndRoute(self.__WB_MANIFEST)
// In dev nothing is precached, so there is no index.html to fall back to.
if (import.meta.env.PROD) {
  registerRoute(new NavigationRoute(createHandlerBoundToURL('index.html')))
}

// PushMessage mirrors the backend's notifications.message. The backend only says what
// happened; the wording comes from i18n (pushNotifications.<type>), in the device's language.
interface PushMessage {
  type: string
  params?: Record<string, string>
  url?: string
}

self.addEventListener('push', (event) => {
  const message = parseMessage(event.data)
  const { t, te } = i18n.global
  const key = te(`pushNotifications.${message.type}.title`) ? message.type : 'fallback'
  const params = formatParams(message.params)

  // Browsers require every push to show a notification, so always show one.
  event.waitUntil(
    self.registration.showNotification(t(`pushNotifications.${key}.title`, params), {
      body: t(`pushNotifications.${key}.body`, params),
      icon: '/pwa-192x192.png',
      data: { url: inAppPath(message.url) },
    }),
  )
})

// Tapping a notification brings the app to the front (or opens it) on the linked page.
self.addEventListener('notificationclick', (event) => {
  event.notification.close()
  const url: string = event.notification.data?.url ?? '/'

  event.waitUntil(
    (async () => {
      const windows = await self.clients.matchAll({ type: 'window', includeUncontrolled: true })
      const client = windows[0]
      if (client) {
        await client.focus()
        try {
          await client.navigate(url)
          return
        } catch {
          // the window isn't controlled by this worker yet — open a new one instead
        }
      }
      await self.clients.openWindow(url)
    })(),
  )
})

function parseMessage(data: PushMessageData | null): PushMessage {
  try {
    const message = data?.json()
    if (message && typeof message.type === 'string') {
      return message
    }
  } catch {
    // not JSON — fall through to the generic notification
  }
  return { type: 'fallback' }
}

// A "date" param arrives as a timestamp and is shown in the device's language and timezone.
function formatParams(params: Record<string, string> | undefined): Record<string, string> {
  const formatted = { ...params }
  if (formatted.date) {
    formatted.date = new Date(formatted.date).toLocaleString(navigator.language, {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    })
  }
  return formatted
}

// Only allow paths inside the app, never links to other sites.
function inAppPath(url: string | undefined): string {
  if (url && url.startsWith('/') && !url.startsWith('//')) {
    return url
  }
  return '/'
}
