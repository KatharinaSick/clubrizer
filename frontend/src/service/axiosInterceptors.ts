import axios from '@/plugins/axios'
import type { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

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

      if (isAuthenticated) {
        config.headers.Authorization = accessToken
      }
      return config
    },
    async (error: AxiosError) => {
      return Promise.reject(error)
    }
  )

  axios.interceptors.response.use(
    (response: AxiosResponse) => {
      return response
    },
    async (error: AxiosError) => {
      const originalRequest = error.config as RetryableRequest

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      // A refresh is already in progress — queue this request to retry once done
      if (isRefreshingToken) {
        return new Promise((resolve) => {
          refreshSubscribers.push((token: string) => {
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
        return axios(originalRequest)
      } catch (refreshError: unknown) {
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
