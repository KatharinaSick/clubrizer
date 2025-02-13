import { defineStore } from 'pinia'
import authService from '@/service/auth'
import { useStorage } from '@vueuse/core'

export interface User {
  email: string,
  givenName: string,
  familyName: string,
  nickName: string,
  picture?: string,
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: useStorage<User>('auth.user', emptyUser),
    accessToken: useStorage<string>('auth.accessToken', ''),
    isAuthenticated: useStorage<boolean>('auth.isAuthenticated', false),
    isLoading: false,
    error: null
  }),
  getters: {
    isLoggedIn: (state) => state.isAuthenticated,
    //accessToken: (state) => state.accessToken
  },
  actions: {
    async login(idToken: string) {
      this.isLoading = true
      this.error = null
      this.isAuthenticated = false
      try {
        const { user, accessToken } = await authService.login(idToken)
        this.user = user
        this.accessToken = accessToken
        this.isAuthenticated = true
      } catch (error: any) {
        this.error = error.response?.data || 'Unknown reason'
      } finally {
        this.isLoading = false
      }
    },
    async refreshTokens() {
      this.isLoading = true
      this.error = null
      this.isAuthenticated = false
      this.accessToken = ""
      try {
        const { user, accessToken } = await authService.refreshTokens()
        this.user = user
        this.accessToken = accessToken
        this.isAuthenticated = true
      } catch (error: any) {
        this.error = error.response?.data || 'Unknown reason'
      } finally {
        this.isLoading = false
      }
    },
    logout() {
      this.user = emptyUser
      this.isAuthenticated = false
      this.isLoading = false
      this.error = null
    }
  }
})

const emptyUser: User = {
  email: '',
  givenName: '',
  familyName: '',
  nickName: '',
  picture: undefined
}
