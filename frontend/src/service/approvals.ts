import axios from '@/plugins/axios'

export interface ApprovalKid {
  id: string
  givenName: string | null
  familyName: string | null
  picture: string | null
}

// An approval request is one account plus any of its pending kids. status is the account's
// own status: 'pending' when the account itself awaits approval, 'approved' when the account
// is already in and only newly-added kids need a decision.
export interface ApprovalRequest {
  userId: string
  email: string
  givenName: string | null
  familyName: string | null
  nickName: string | null
  picture: string | null
  status: 'onboarding' | 'pending' | 'approved' | 'rejected'
  selfParticipates: boolean
  pendingKids: ApprovalKid[]
}

export interface ApprovalDecision {
  userIds: string[]
  kidIds: string[]
}

const approvalsService = {
  async listApprovals(): Promise<ApprovalRequest[]> {
    const response = await axios.get('/admin/approvals')
    return response.data ?? []
  },

  async approve(decision: ApprovalDecision): Promise<void> {
    await axios.post('/admin/approvals/approve', decision)
  },

  async reject(decision: ApprovalDecision): Promise<void> {
    await axios.post('/admin/approvals/reject', decision)
  },
}

export default approvalsService
