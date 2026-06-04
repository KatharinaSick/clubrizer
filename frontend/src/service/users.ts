import axios from '@/plugins/axios'
import { jwtDecode } from 'jwt-decode'
import type { User } from '@/stores/auth'

const usersService = {
  async updateProfile(
    firstName: string,
    lastName: string,
    nickName?: string,
  ): Promise<{ user: User; accessToken: string }> {
    const response = await axios.patch(
      '/users/me/profile',
      { firstName, lastName, nickName: nickName || null },
      { withCredentials: true },
    )
    const { accessToken } = response.data
    return { user: jwtDecode<User>(accessToken), accessToken }
  },

  async updateProfilePicture(file: File): Promise<{ user: User; accessToken: string }> {
    const formData = new FormData()
    formData.append('picture', file)
    const response = await axios.post('/users/me/picture', formData, { withCredentials: true })
    const { accessToken } = response.data
    return { user: jwtDecode<User>(accessToken), accessToken }
  },
}

export default usersService
