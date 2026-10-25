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

// Some types are grouped: several in a row are combined into one notification instead of
// piling up (e.g. ten new events added at once). The notification keeps the items it shows
// in its data, so the next push can replace it with an updated summary without alerting
// again. Once it's dismissed, the count starts over. Browsers that can't look up shown
// notifications (getNotifications) simply show each one on its own.
interface Group {
  tag: string
  // item is this push's line in the summary
  item: string
  // summary turns all items so far into the grouped notification's params and link
  summary: (items: string[]) => { params: Record<string, string>; url: string }
}

function groupOf(message: PushMessage): Group | undefined {
  const params = message.params ?? {}
  switch (message.type) {
    case 'newEvent':
      return {
        tag: 'new-events',
        item: params.title ?? '',
        summary: (items) => ({
          params: { count: String(items.length), titles: items.join(', ') },
          url: '/events',
        }),
      }
    case 'eventResponse':
      return {
        tag: `rsvp-${params.eventId}`,
        item: `${params.name} ${params.response === 'yes' ? '✅' : '❌'}`,
        summary: (items) => ({
          params: {
            count: String(items.length),
            title: params.title ?? '',
            answers: items.join(', '),
          },
          url: inAppPath(message.url),
        }),
      }
  }
  return undefined
}

self.addEventListener('push', (event) => {
  // Browsers require every push to show a notification, so always show one.
  event.waitUntil(show(parseMessage(event.data)))
})

async function show(message: PushMessage) {
  const { t, te } = i18n.global
  const key = te(`pushNotifications.${message.type}.title`)
    ? message.type
    : 'fallback'
  const options: NotificationOptions = { icon: '/pwa-192x192.png' }
  let params = formatParams(message.params)
  let title = t(`pushNotifications.${key}.title`, params)
  let url = inAppPath(message.url)

  const group = key === 'fallback' ? undefined : groupOf(message)
  if (group) {
    const items = [...(await previousItems(group.tag)), group.item]
    if (items.length > 1) {
      const summary = group.summary(items)
      params = summary.params
      title = t(`pushNotifications.${key}.groupTitle`, params)
      options.body = t(`pushNotifications.${key}.groupBody`, params)
      url = summary.url
    }
    options.tag = group.tag
    options.data = { url, items }
  }

  options.body ??= t(`pushNotifications.${key}.body`, params)
  options.data ??= { url }
  await self.registration.showNotification(title, options)
}

// previousItems returns the items of the notification with this tag that is still shown, if any.
async function previousItems(tag: string): Promise<string[]> {
  try {
    const [previous] = await self.registration.getNotifications({ tag })
    const items = previous?.data?.items
    return Array.isArray(items) ? items : []
  } catch {
    return []
  }
}

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
// A "response" param ("yes"/"no") becomes the translated "answer".
function formatParams(params: Record<string, string> | undefined): Record<string, string> {
  const formatted = { ...params }
  if (formatted.response === 'yes' || formatted.response === 'no') {
    formatted.answer = i18n.global.t(`pushNotifications.eventResponse.${formatted.response}`)
  }
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
