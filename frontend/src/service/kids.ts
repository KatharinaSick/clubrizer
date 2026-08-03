import axios from '@/plugins/axios'
import type { UserStatus } from '@/stores/auth'

export interface Kid {
  id: string
  givenName: string | null
  familyName: string | null
  picture?: string | null
  status: UserStatus
}

const kidsService = {
  async listKids(): Promise<Kid[]> {
    const response = await axios.get('/users/me/kids')
    return response.data
  },

  async addKid(firstName: string, lastName: string): Promise<Kid> {
    const response = await axios.post('/users/me/kids', { firstName, lastName })
    return response.data
  },

  // replaceKids saves the whole onboarding list of first names in one request; the
  // backend makes the account's kids match it (adds new, drops removed). This is the
  // onboarding "add your kids" step, distinct from the post-approval kid management above.
  async replaceKids(firstNames: string[]): Promise<Kid[]> {
    const response = await axios.post('/users/me/onboarding/kids', { firstNames })
    return response.data
  },

  async updateKid(id: string, firstName: string, lastName: string): Promise<Kid> {
    const response = await axios.patch(`/users/me/kids/${id}`, { firstName, lastName })
    return response.data
  },

  async removeKid(id: string): Promise<void> {
    await axios.delete(`/users/me/kids/${id}`)
  },

  async updateKidPicture(id: string, file: File): Promise<Kid> {
    const formData = new FormData()
    formData.append('picture', file)
    const response = await axios.post(`/users/me/kids/${id}/picture`, formData)
    return response.data
  },
}

export default kidsService
