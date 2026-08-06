import axios from '@/plugins/axios'

export interface MemberRole {
  id: string
  name: string
}

export interface MemberKid {
  givenName: string | null
  familyName: string | null
  picture: string | null
}

// A member is one approved account in the club roster, with its roles (the implicit base
// membership role is already excluded by the backend) and its approved kids.
export interface Member {
  email: string
  givenName: string | null
  familyName: string | null
  nickName: string | null
  picture: string | null
  selfParticipates: boolean
  roles: MemberRole[]
  kids: MemberKid[]
}

const membersService = {
  async listMembers(): Promise<Member[]> {
    const response = await axios.get('/admin/members')
    return response.data ?? []
  },
}

export default membersService
