import { defineStore } from 'pinia'
import authService from '@/service/auth'
import usersService from '@/service/users'
import { useStorage } from '@vueuse/core'

export type UserStatus = 'pending' | 'approved' | 'rejected'

export interface User {
  email: string
  givenName: string | null
  familyName: string | null
  nickName: string | null
  picture?: string | null
  status: UserStatus
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
  },
  actions: {
    async requestOTP(email: string) {
      await authService.requestOTP(email)
    },
    async verifyOTP(email: string, code: string) {
      const { user, accessToken } = await authService.verifyOTP(email, code)
      this.user = user
      this.accessToken = accessToken
      this.isAuthenticated = true
    },
    async updateProfile(firstName: string, lastName: string, nickName?: string) {
      const { user, accessToken } = await usersService.updateProfile(firstName, lastName, nickName)
      this.user = user
      this.accessToken = accessToken
    },
    async updateProfilePicture(file: File) {
      const { user, accessToken } = await usersService.updateProfilePicture(file)
      this.user = user
      this.accessToken = accessToken
    },
    async refreshTokens() {
      this.isLoading = true
      this.isAuthenticated = false
      this.accessToken = ''
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
      } catch {
        // ignore error
      } finally {
        this.user = emptyUser
        this.isAuthenticated = false
        this.accessToken = ''
        this.isLoading = false
      }
    },
  },
})

const emptyUser: User = {
  email: '',
  givenName: null,
  familyName: null,
  nickName: null,
  picture: undefined,
  status: 'pending',
}
