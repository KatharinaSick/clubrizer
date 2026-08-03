import { defineStore } from 'pinia'
import authService from '@/service/auth'
import usersService from '@/service/users'
import { useStorage } from '@vueuse/core'

export type UserStatus = 'onboarding' | 'pending' | 'approved' | 'rejected'

// Permission key that gates the member-management screen. Mirrors the backend's
// rbac.PermissionUsersApprove — the backend is always the source of truth; the permission
// in the token is only used to show or hide UI.
export const PERMISSION_MANAGE_MEMBERS = 'users.approve'

export interface User {
  email: string
  givenName: string | null
  familyName: string | null
  nickName: string | null
  picture?: string | null
  status: UserStatus
  // Whether the account holder participates in events themselves. false = a guardian
  // account ("only my kids") that only manages kids and never RSVPs for itself.
  selfParticipates: boolean
  // Effective permission keys from the JWT — a UI hint only (see PERMISSION_MANAGE_MEMBERS).
  permissions: string[]
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
    // Whether the account may manage members (review the approval queue). Backend-enforced;
    // this only decides UI visibility.
    canManageMembers: (state) => (state.user.permissions ?? []).includes(PERMISSION_MANAGE_MEMBERS),
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
    async setAccountType(selfParticipates: boolean) {
      const { user, accessToken } = await usersService.setAccountType(selfParticipates)
      this.user = user
      this.accessToken = accessToken
    },
    async finishOnboarding() {
      const { user, accessToken } = await usersService.finishOnboarding()
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
  selfParticipates: true,
  permissions: [],
}
