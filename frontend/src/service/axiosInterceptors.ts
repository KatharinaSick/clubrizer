import axios from '@/plugins/axios'
import type { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import router from '@/router'
import i18n from '@/plugins/i18n'

interface RetryableRequest extends AxiosRequestConfig {
  _retry?: boolean
}

let isRefreshingToken = false
let refreshSubscribers: ((token: string) => void)[] = []

const configureAxiosInterceptors = (): void => {
  setUpRefreshTokenInterceptors()
}

const setUpRefreshTokenInterceptors = (): void => {
  axios.interceptors.request.use(
    (config) => {
      const { accessToken, isAuthenticated } = useAuthStore()

      // The token refresh request (/oauth/tokens) is a silent background operation —
      // it's invisible to the user, so we don't want it to trigger the loading indicator.
      if (config.url !== '/oauth/tokens') {
        useRequestStore().startRequest()
      }

      if (isAuthenticated) {
        config.headers.Authorization = accessToken
      }
      return config
    },
    async (error: AxiosError) => {
      if (error.config?.url !== '/oauth/tokens') {
        useRequestStore().endRequest()
      }
      return Promise.reject(error)
    }
  )

  axios.interceptors.response.use(
    (response: AxiosResponse) => {
      if (response.config.url !== '/oauth/tokens') {
        useRequestStore().endRequest()
      }
      return response
    },
    async (error: AxiosError) => {
      const originalRequest = error.config as RetryableRequest

      if (error.response?.status !== 401 || originalRequest._retry) {
        if (originalRequest.url !== '/oauth/tokens') {
          useRequestStore().endRequest()
        }
        // Show a global error for everything except:
        // - 401: handled by the refresh flow above
        // - 422: form validation errors, handled locally by the component
        const status = error.response?.status
        if (status !== 422) {
          const rawData = error.response?.data
          const message = (rawData != null && String(rawData).trim()) || i18n.global.t('request.unexpectedError')
          useRequestStore().setError(message)
        }
        return Promise.reject(error)
      }

      // A refresh is already in progress — queue this request to retry once done.
      if (isRefreshingToken) {
        return new Promise((resolve) => {
          refreshSubscribers.push((token: string) => {
            // Balance the startRequest() that was called for the original request.
            // The retry below will go through the interceptors and call startRequest() again.
            useRequestStore().endRequest()
            originalRequest.headers = { ...originalRequest.headers }
            originalRequest.headers!.Authorization = token
            resolve(axios(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshingToken = true

      const { refreshTokens, logout, accessToken } = useAuthStore()

      try {
        await refreshTokens()
        refreshSubscribers.forEach((cb) => cb(accessToken))
        originalRequest.headers = { ...originalRequest.headers }
        originalRequest.headers!.Authorization = accessToken
        // Balance the startRequest() for the original request before retrying.
        // The retry goes through the interceptors and calls startRequest() again.
        useRequestStore().endRequest()
        return axios(originalRequest)
      } catch (refreshError: unknown) {
        // Balance the startRequest() for the original request — it will never be retried.
        useRequestStore().endRequest()
        refreshSubscribers.forEach((cb) => cb(''))
        logout()
        router.push('/signin')
        return Promise.reject(refreshError)
      } finally {
        refreshSubscribers = []
        isRefreshingToken = false
      }
    }
  )
}
export { configureAxiosInterceptors }
