import axios from '@/plugins/axios'
import { jwtDecode } from 'jwt-decode'
import type { User } from '@/stores/auth'

const authService = {
  async login(idToken: string): Promise<{user: User, accessToken: string}> {
    try {
      const response = await axios.post('/signin', { idToken }, { withCredentials: true })
      const { accessToken } = response.data

      return {user: jwtDecode<User>(accessToken), accessToken}
    } catch (error: any) {
      throw error
    }
  },

  async refreshTokens(): Promise<{user: User, accessToken: string}> {
    try {
      const response = await axios.post('/oauth/tokens', undefined, { withCredentials: true })
      const { accessToken } = response.data

      return {user: jwtDecode<User>(accessToken), accessToken}
    } catch (error: any) {
      throw error
    }
  },
}

export default authService
