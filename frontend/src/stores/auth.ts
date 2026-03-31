import { defineStore } from 'pinia'
import authService from '@/service/auth'
import { useStorage } from '@vueuse/core'

export type UserStatus = 'pending' | 'approved' | 'rejected'

export interface User {
  email: string,
  givenName: string,
  familyName: string,
  nickName: string,
  picture?: string,
  status: UserStatus,
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: useStorage<User>('auth.user', emptyUser),
    accessToken: useStorage<string>('auth.accessToken', ''),
    isAuthenticated: useStorage<boolean>('auth.isAuthenticated', false),
    isLoading: false,
  }),
  getters: {
    isLoggedIn: (state) => state.isAuthenticated,
    //accessToken: (state) => state.accessToken
  },
  actions: {
    async login(idToken: string) {
      this.isLoading = true
      this.isAuthenticated = false
      try {
        const { user, accessToken } = await authService.login(idToken)
        this.user = user
        this.accessToken = accessToken
        this.isAuthenticated = true
      } catch {
        // error is handled globally by the axios interceptor
      } finally {
        this.isLoading = false
      }
    },
    async refreshTokens() {
      this.isLoading = true
      this.isAuthenticated = false
      this.accessToken = ""
      try {
        const { user, accessToken } = await authService.refreshTokens()
        this.user = user
        this.accessToken = accessToken
        this.isAuthenticated = true
      } catch {
        // error is handled globally by the axios interceptor
      } finally {
        this.isLoading = false
      }
    },
    async logout() {
      try {
        await authService.logout()
      } catch (error) {
        // ignore error
      } finally {
        this.user = emptyUser
        this.isAuthenticated = false
        this.accessToken = ""
        this.isLoading = false
      }
    }
  }
})

const emptyUser: User = {
  email: '',
  givenName: '',
  familyName: '',
  nickName: '',
  picture: undefined,
  status: 'pending',
}
