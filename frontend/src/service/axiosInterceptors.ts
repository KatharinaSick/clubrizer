import axios from '@/plugins/axios'
import type { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

interface RetryableRequest extends AxiosRequestConfig {
  _retry?: boolean
}

let isRefreshingToken = false

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
      const { refreshTokens, logout, accessToken } = useAuthStore()

      const originalRequest = error.config as RetryableRequest

      if (error.response?.status === 401 && !originalRequest._retry && !isRefreshingToken) {
        originalRequest._retry = true
        isRefreshingToken = true

        try {
          await refreshTokens()
          originalRequest.headers = { ...originalRequest.headers };
          originalRequest.headers.Authorization = accessToken;
          return axios(originalRequest)
        } catch (error: unknown) {
          logout()
          router.push('/signin')
          return Promise.reject(error)
        } finally {
          isRefreshingToken = false
        }
      }

      return Promise.reject(error)
    }
  )
}
export { configureAxiosInterceptors }
