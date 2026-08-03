import { defineStore } from 'pinia'
import approvalsService, { type ApprovalRequest, type ApprovalDecision } from '@/service/approvals'
import { useAuthStore } from '@/stores/auth'

// Holds the approval queue for members who can manage members. Kept in a store (not a view)
// because the pending state drives red badges across the app — the profile icon in the
// global navigation, the profile menu gear, and the "Manage Members" menu item.
export const useApprovalsStore = defineStore('approvals', {
  state: () => ({
    requests: [] as ApprovalRequest[],
    loaded: false,
  }),
  getters: {
    hasPending: (state) => state.requests.length > 0,
    pendingCount: (state) => state.requests.length,
  },
  actions: {
    // refresh reloads the queue. For an account that can't manage members it simply clears
    // the list (no request), so the badges are always off for them.
    async refresh() {
      const auth = useAuthStore()
      if (!auth.canManageMembers) {
        this.requests = []
        this.loaded = true
        return
      }
      this.requests = await approvalsService.listApprovals()
      this.loaded = true
    },

    async approve(decision: ApprovalDecision) {
      await approvalsService.approve(decision)
      await this.refresh()
    },

    async reject(decision: ApprovalDecision) {
      await approvalsService.reject(decision)
      await this.refresh()
    },
  },
})
