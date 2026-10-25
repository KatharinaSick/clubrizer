import axios from '@/plugins/axios'
import type { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import router from '@/router'
import i18n from '@/plugins/i18n'

declare module 'axios' {
  interface AxiosRequestConfig {
    // silent marks a background request the user didn't ask for (e.g. re-registering this
    // device for notifications after sign-in). It shows no loading bar and no global error.
    silent?: boolean
  }
}

interface RetryableRequest extends AxiosRequestConfig {
  _retry?: boolean
}

// The token refresh request is a silent background operation too — it's invisible to the
// user, so it doesn't count towards the loading indicator.
function isTracked(config?: AxiosRequestConfig): boolean {
  return config?.url !== '/auth/refresh' && !config?.silent
}

let isRefreshingToken = false
// Requests waiting for a running refresh. They get the new token, or null if the refresh
// failed — then they are rejected instead of retried without a token.
let refreshSubscribers: ((token: string | null) => void)[] = []

// These endpoints are part of the auth flow itself — a 401 from them should
// never trigger a token refresh attempt.
const noRefreshUrls = ['/auth/refresh', '/auth/otp', '/auth/verify']

const configureAxiosInterceptors = (): void => {
  setUpRefreshTokenInterceptors()
}

const setUpRefreshTokenInterceptors = (): void => {
  axios.interceptors.request.use(
    (config) => {
      const { accessToken, isAuthenticated } = useAuthStore()

      if (isTracked(config)) {
        useRequestStore().startRequest()
      }

      if (isAuthenticated) {
        config.headers.Authorization = accessToken
      }
      return config
    },
    async (error: AxiosError) => {
      if (isTracked(error.config)) {
        useRequestStore().endRequest()
      }
      return Promise.reject(error)
    },
  )

  axios.interceptors.response.use(
    (response: AxiosResponse) => {
      if (isTracked(response.config)) {
        useRequestStore().endRequest()
      }
      return response
    },
    async (error: AxiosError) => {
      const originalRequest = error.config as RetryableRequest
      const isNoRefreshUrl = noRefreshUrls.includes(originalRequest.url ?? '')

      if (error.response?.status !== 401 || originalRequest._retry || isNoRefreshUrl) {
        if (isTracked(originalRequest)) {
          useRequestStore().endRequest()
        }
        // Show a global error for everything except:
        // - 401 (handled by the refresh flow, unless it's a no-refresh URL)
        // - 422: form validation errors, handled locally by the component
        // - silent requests: the user didn't ask for them, so they shouldn't see their errors
        const status = error.response?.status
        if (status !== 422 && !originalRequest.silent) {
          const rawData = error.response?.data
          const message =
            (rawData != null && String(rawData).trim()) ||
            i18n.global.t('request.unexpectedError')
          useRequestStore().setError(message)
        }
        return Promise.reject(error)
      }

      // A refresh is already in progress — queue this request to retry once done.
      if (isRefreshingToken) {
        return new Promise((resolve, reject) => {
          refreshSubscribers.push((token: string | null) => {
            // Balance the startRequest() that was called for the original request.
            // The retry below will go through the interceptors and call startRequest() again.
            if (isTracked(originalRequest)) {
              useRequestStore().endRequest()
            }
            if (token === null) {
              reject(error)
              return
            }
            originalRequest.headers = { ...originalRequest.headers }
            originalRequest.headers!.Authorization = token
            resolve(axios(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshingToken = true

      // Keep the store itself, not destructured values: the new access token only exists
      // after the refresh, so it has to be read from the store afterwards.
      const auth = useAuthStore()

      try {
        await auth.refreshTokens()
        const accessToken = auth.accessToken
        refreshSubscribers.forEach((cb) => cb(accessToken))
        originalRequest.headers = { ...originalRequest.headers }
        originalRequest.headers!.Authorization = accessToken
        // Balance the startRequest() for the original request before retrying.
        // The retry goes through the interceptors and calls startRequest() again.
        if (isTracked(originalRequest)) {
          useRequestStore().endRequest()
        }
        return axios(originalRequest)
      } catch (refreshError: unknown) {
        // Balance the startRequest() for the original request — it will never be retried.
        if (isTracked(originalRequest)) {
          useRequestStore().endRequest()
        }
        refreshSubscribers.forEach((cb) => cb(null))
        auth.logout()
        router.push('/signin')
        return Promise.reject(refreshError)
      } finally {
        refreshSubscribers = []
        isRefreshingToken = false
      }
    },
  )
}
export { configureAxiosInterceptors }
