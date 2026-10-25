import axios from '@/plugins/axios'

// Shape of PushSubscription.toJSON() that the backend stores.
export interface PushSubscriptionPayload {
  endpoint: string
  keys: {
    p256dh: string
    auth: string
  }
}

// One switch on the notification settings screen. The backend decides which categories a
// user can see (member requests only for people who can approve members).
export interface NotificationSetting {
  category: string
  enabled: boolean
}

const notificationsService = {
  async getVapidPublicKey(): Promise<string> {
    const response = await axios.get('/notifications/vapid-public-key')
    return response.data.publicKey
  },

  async subscribe(subscription: PushSubscriptionPayload, silent = false): Promise<void> {
    await axios.post('/notifications/subscriptions', subscription, { silent })
  },

  async isSubscribed(endpoint: string, silent = false): Promise<boolean> {
    const response = await axios.post('/notifications/subscriptions/status', { endpoint }, { silent })
    return response.data.subscribed
  },

  async unsubscribe(endpoint: string): Promise<void> {
    await axios.delete('/notifications/subscriptions', { data: { endpoint } })
  },

  // releaseDevice removes this browser's registration whoever it belongs to; the subscription's
  // secret (auth) proves we're on that device.
  // It runs in the background (sign-in, logout), so it never shows errors.
  async releaseDevice(endpoint: string, auth: string): Promise<void> {
    await axios.post('/notifications/subscriptions/release', { endpoint, auth }, { silent: true })
  },

  // sendTest returns how many of the user's devices accepted the test notification.
  async sendTest(): Promise<number> {
    const response = await axios.post('/notifications/test')
    return response.data.delivered
  },

  async getSettings(): Promise<NotificationSetting[]> {
    const response = await axios.get('/notifications/settings')
    return response.data
  },

  async updateSetting(category: string, enabled: boolean): Promise<void> {
    await axios.put('/notifications/settings', { category, enabled })
  },
}

export default notificationsService
