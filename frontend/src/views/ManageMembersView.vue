<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useApprovalsStore } from '@/stores/approvals'
import { useRequestStore } from '@/stores/request'
import Avatar from '@/components/Avatar.vue'
import Header from '@/components/Header.vue'
import IconCheck from '@/components/icons/IconCheck.vue'
import IconClose from '@/components/icons/IconClose.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'
import type { ApprovalRequest, ApprovalKid } from '@/service/approvals'
import membersService, { type Member, type MemberKid } from '@/service/members'

const approvals = useApprovalsStore()
const requestStore = useRequestStore()

// Which card + action is in flight, so we can show a spinner on the pressed button and lock
// the card against a double-submit.
const busy = ref<{ id: string; action: 'approve' | 'reject' } | null>(null)

const members = ref<Member[]>([])
const membersLoaded = ref(false)

onMounted(() => {
  approvals.refresh().catch(() => {
    // error shown globally via RequestError
  })
  membersService
    .listMembers()
    .then((m) => { members.value = m })
    .catch(() => {
      // error shown globally via RequestError
    })
    .finally(() => { membersLoaded.value = true })
})

function memberName(member: Member): string {
  const full = [member.givenName, member.familyName].filter(Boolean).join(' ')
  return full || member.nickName || member.email
}

function memberKidName(kid: MemberKid): string {
  return [kid.givenName, kid.familyName].filter(Boolean).join(' ')
}

function accountName(req: ApprovalRequest): string {
  const full = [req.givenName, req.familyName].filter(Boolean).join(' ')
  return req.nickName || full || req.email
}

function kidName(kid: ApprovalKid): string {
  return [kid.givenName, kid.familyName].filter(Boolean).join(' ')
}

function accountTypeLabel(req: ApprovalRequest): string {
  if (req.status === 'approved') return i18n.global.t('manageMembers.requests.newKids')
  if (!req.selfParticipates) return i18n.global.t('manageMembers.requests.guardian')
  return i18n.global.t('manageMembers.requests.newMember')
}

async function decide(req: ApprovalRequest, action: 'approve' | 'reject') {
  if (busy.value) return
  requestStore.clearError()
  busy.value = { id: req.userId, action }

  // An already-approved account (a later-added kid) must not be sent as a user id — only its
  // pending kids are decided. A pending account is sent together with its kids.
  const decision = {
    userIds: req.status === 'pending' ? [req.userId] : [],
    kidIds: req.pendingKids.map((k) => k.id),
  }

  try {
    if (action === 'approve') {
      await approvals.approve(decision)
    } else {
      await approvals.reject(decision)
    }
  } catch {
    // error shown globally via RequestError
  } finally {
    busy.value = null
  }
}

function isCardBusy(req: ApprovalRequest): boolean {
  return busy.value?.id === req.userId
}
</script>

<template>
  <div class="manageMembers">
    <div class="manageMembersCard">
      <Header :title="i18n.global.t('manageMembers.header')" left-action="back" />

      <RequestError class="manageMembersError" />

      <section class="manageMembersSection">
        <h2 class="manageMembersSectionTitle">
          {{ $t('manageMembers.requests.title') }}
          <span v-if="approvals.hasPending" class="manageMembersCount">{{ approvals.pendingCount }}</span>
        </h2>

        <p v-if="approvals.loaded && !approvals.hasPending" class="manageMembersEmpty">
          {{ $t('manageMembers.requests.empty') }}
        </p>

        <div v-else class="manageMembersList">
          <div
            v-for="req in approvals.requests"
            :key="req.userId"
            class="manageMembersRequest"
            :class="{ manageMembersRequestBusy: isCardBusy(req) }"
          >
            <div class="manageMembersRequestInfo">
              <span class="manageMembersRequestName">{{ accountName(req) }}</span>
              <span class="manageMembersRequestBadge">{{ accountTypeLabel(req) }}</span>
              <template v-if="req.pendingKids.length">
                <span class="manageMembersKidsLabel">{{ $t('manageMembers.requests.kids') }}</span>
                <ul class="manageMembersKids">
                  <li v-for="kid in req.pendingKids" :key="kid.id" class="manageMembersKid">{{ kidName(kid) }}</li>
                </ul>
              </template>
            </div>

            <div class="manageMembersActions">
              <button
                type="button"
                class="manageMembersAction manageMembersActionApprove"
                :disabled="busy !== null"
                :aria-label="$t('manageMembers.requests.approve')"
                @click="decide(req, 'approve')"
              >
                <IconCheck class="manageMembersActionIcon" />
              </button>
              <button
                type="button"
                class="manageMembersAction manageMembersActionReject"
                :disabled="busy !== null"
                :aria-label="$t('manageMembers.requests.reject')"
                @click="decide(req, 'reject')"
              >
                <IconClose class="manageMembersActionIcon" />
              </button>
            </div>
          </div>
        </div>
      </section>

      <section class="manageMembersSection">
        <h2 class="manageMembersSectionTitle">
          {{ $t('manageMembers.members.title') }}
          <span v-if="members.length" class="manageMembersCount">{{ members.length }}</span>
        </h2>

        <p v-if="membersLoaded && members.length === 0" class="manageMembersEmpty">
          {{ $t('manageMembers.members.empty') }}
        </p>

        <div v-else class="manageMembersList">
          <div
            v-for="(member, index) in members"
            :key="index"
            class="manageMembersMember"
          >
            <Avatar
              interactive
              :picture="member.picture"
              :given-name="member.givenName"
              :family-name="member.familyName"
              :nick-name="member.nickName"
              size="md"
            />
            <div class="manageMembersMemberInfo">
              <span class="manageMembersMemberName">{{ memberName(member) }}</span>
              <span class="manageMembersMemberEmail">{{ member.email }}</span>
              <div v-if="!member.selfParticipates || member.roles.length" class="manageMembersBadges">
                <span v-if="!member.selfParticipates" class="manageMembersBadge">
                  {{ $t('manageMembers.members.guardian') }}
                </span>
                <span
                  v-for="role in member.roles"
                  :key="role.id"
                  class="manageMembersBadge manageMembersBadgeRole"
                >{{ role.name }}</span>
              </div>
              <template v-if="member.kids.length">
                <span class="manageMembersKidsLabel">{{ $t('manageMembers.members.kids') }}</span>
                <ul class="manageMembersKids">
                  <li v-for="(kid, kidIndex) in member.kids" :key="kidIndex" class="manageMembersKid">{{ memberKidName(kid) }}</li>
                </ul>
              </template>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.manageMembersError {
  margin-bottom: var(--padding);
}

.manageMembersSectionTitle {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--padding);
}

.manageMembersCount {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  box-sizing: border-box;
  border-radius: 10px;
  background: var(--blue);
  color: var(--white);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

.manageMembersEmpty {
  color: var(--text-gray);
  font-size: var(--font-size-medium);
  text-align: center;
  margin: calc(var(--padding) * 2) 0;
}

.manageMembersList {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.manageMembersRequest {
  padding: var(--padding);
  border: 1px solid var(--gray);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  gap: var(--padding);
}

.manageMembersRequestBusy {
  opacity: 0.5;
}

.manageMembersRequestInfo {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.manageMembersRequestName {
  font-weight: var(--font-weight-medium);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.manageMembersRequestBadge {
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.manageMembersKidsLabel {
  margin-top: 8px;
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  color: var(--text-gray);
}

.manageMembersKids {
  list-style: none;
  margin: 2px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.manageMembersKid {
  font-size: var(--font-size-small);
  color: var(--text-light);
}

.manageMembersSection + .manageMembersSection {
  margin-top: calc(var(--padding) * 2);
}

.manageMembersMember {
  padding: var(--padding);
  border: 1px solid var(--gray);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  gap: var(--padding);
}

.manageMembersMemberInfo {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
  flex: 1;
}

.manageMembersMemberName {
  font-weight: var(--font-weight-medium);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.manageMembersMemberEmail {
  font-size: var(--font-size-small);
  color: var(--text-gray);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.manageMembersBadges {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.manageMembersBadge {
  font-size: var(--font-size-small);
  color: var(--text-gray);
  background: var(--light-gray);
  padding: 2px var(--gap);
  border-radius: var(--border-radius);
}

.manageMembersBadgeRole {
  color: var(--blue);
  background: var(--light-blue);
  text-transform: capitalize;
}

.manageMembersActions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.manageMembersAction {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.manageMembersAction:disabled {
  cursor: default;
}

.manageMembersActionApprove {
  background: var(--light-green);
  color: var(--green);
}

.manageMembersActionReject {
  background: var(--light-red);
  color: var(--red);
}

.manageMembersActionIcon {
  width: 22px;
  height: 22px;
}

@media (min-width: 768px) {
  .manageMembers {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .manageMembersCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }
}
</style>
