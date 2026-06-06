import axios from '@/plugins/axios'
import { jwtDecode } from 'jwt-decode'
import type { User } from '@/stores/auth'

const authService = {
  async requestOTP(email: string): Promise<void> {
    await axios.post('/auth/otp', { email })
  },

  async verifyOTP(email: string, code: string): Promise<{ user: User; accessToken: string }> {
    const response = await axios.post('/auth/verify', { email, code }, { withCredentials: true })
    const { accessToken } = response.data
    return { user: jwtDecode<User>(accessToken), accessToken }
  },

  async refreshTokens(): Promise<{ user: User; accessToken: string }> {
    const response = await axios.post('/auth/refresh', undefined, { withCredentials: true })
    const { accessToken } = response.data
    return { user: jwtDecode<User>(accessToken), accessToken }
  },

  async logout(): Promise<void> {
    await axios.post('/auth/logout', undefined, { withCredentials: true })
  },
}

export default authService
