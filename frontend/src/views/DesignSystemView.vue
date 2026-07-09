<script setup lang="ts">
// Design system reference page — strings are intentionally hardcoded (technical documentation, not app content)
import { ref } from 'vue'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'
import Select from '@/components/Select.vue'
import type { Option } from '@/components/Select.vue'
import Alert from '@/components/Alert.vue'
import Avatar from '@/components/Avatar.vue'
import Divider from '@/components/Divider.vue'
import Event from '@/components/Event.vue'
import type { EventProps } from '@/components/EventTitle.vue'
import Header from '@/components/Header.vue'
import MenuButton from '@/components/MenuButton.vue'
import type { MenuItem } from '@/components/MenuButton.vue'

const menuItems: MenuItem[] = [
  { label: 'Edit', onClick: () => {} },
  { label: 'Delete', danger: true, onClick: () => {} },
]

const textInput = ref('')
const errorInput = ref('')
const dateInput = ref('')
const timeInput = ref('')
const textareaInput = ref('')
const selectValue = ref<Option | null>(null)
const selectLoadingValue = ref<Option | null>(null)

const sampleEvent: EventProps = {
  id: 'ds-1',
  title: 'Skate Training',
  description: 'Weekly skate training session at the park',
  category: {
    name: 'Skatetraining',
    color: '#0762EF',
    picture: 'https://storage.googleapis.com/clubrizer-static-images/event-categories/skatetraining.png'
  },
  location: 'Skatepark Donauinsel',
  startTime: new Date(2026, 5, 28, 17, 0).toISOString(),
  creator: {
    givenName: 'Katharina',
    familyName: 'Sick',
    nickName: 'Kathi'
  }
}

const selectOptions: Option[] = [
  { id: '1', name: 'Skatetraining' },
  { id: '2', name: 'Meeting' },
  { id: '3', name: 'Workshop' }
]

const colors = [
  { name: '--blue', value: '#0E6EF0', label: 'Blue · Primary flat' },
  { name: '--green', value: '#28E8A2', label: 'Green · Success / confirm' },
  { name: '--red', value: '#F3675E', label: 'Red · Error / decline' },
  { name: '--light-red', value: '#FCD7D4', label: 'Light Red · Error bg' },
  { name: '--orange', value: '#FF9F43', label: 'Orange · Warning' },
  { name: '--light-orange', value: '#FFF0DE', label: 'Light Orange · Warning bg' },
  { name: '--gray', value: '#E2E0E0', label: 'Gray · Borders' },
  { name: '--light-gray', value: '#F5F5F5', label: 'Light Gray · Input bg' },
  { name: '--text-light', value: '#0C1D36', label: 'Text Primary' },
  { name: '--text-gray', value: '#B6BBC3', label: 'Text Secondary' },
  { name: '--white', value: '#FFFFFF', label: 'White · Surface / card bg' },
]

const spacingGeneral = [
  { name: '--gap', value: '8px', desc: 'Gaps between inline or adjacent elements (flex gap, grid gap)' },
  { name: '--padding', value: '12px', desc: 'Component and page padding — the base unit for all layouts' },
]
const spacingDerived = [
  { name: '2× --padding', value: '24px', desc: 'Section spacing, label vertical offset' },
  { name: '4× --padding', value: '48px', desc: 'Large page section gaps' },
]
const spacingLayout = [
  { name: '--content-max-width', value: '1100px', desc: 'Max width for all main content areas — center with margin: 0 auto' },
]
</script>

<template>
  <div class="dsPage">

    <!-- ─── PAGE HEADER ─── -->
    <header class="dsPageHeader">
      <h1 class="dsPageTitle">Clubrizer Design System</h1>
      <p class="dsPageSubtitle">
        Design tokens, components, and patterns that define Clubrizer's visual language.
        Use this as a reference when building or reviewing UI — mobile and desktop.
      </p>
    </header>

    <!-- ─── DESIGN TOKENS ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Design Tokens</h2>
      <p class="dsSectionDesc">Core values defined in <code>base.css</code>. Never hardcode colors, sizes, or spacing — always use the corresponding variable.</p>

      <!-- Colors -->
      <h3 class="dsSubsectionTitle">Colors</h3>
      <div class="dsColorGrid">
        <div class="dsColorSwatch">
          <div class="dsColorBlock" style="background: var(--horizontal-gradient)" />
          <code class="dsToken">--horizontal-gradient</code>
          <span class="dsColorHex">#0762EF → #52E3FE</span>
          <span class="dsColorLabel">Primary gradient (left→right)</span>
        </div>
        <div class="dsColorSwatch">
          <div class="dsColorBlock" style="background: var(--gradient)" />
          <code class="dsToken">--gradient</code>
          <span class="dsColorHex">#0762EF → #52E3FE</span>
          <span class="dsColorLabel">Gradient (bottom-left→top-right) · FAB, modal border, avatar</span>
        </div>
        <div
          v-for="color in colors"
          :key="color.name"
          class="dsColorSwatch"
        >
          <div
            class="dsColorBlock"
            :style="{
              background: color.value,
              border: color.value === '#FFFFFF' || color.value === '#F5F5F5' || color.value === '#E2E0E0'
                ? '1px solid var(--gray)'
                : 'none'
            }"
          />
          <code class="dsToken">{{ color.name }}</code>
          <span class="dsColorHex">{{ color.value }}</span>
          <span class="dsColorLabel">{{ color.label }}</span>
        </div>
      </div>
      <div class="dsNote">
        The gradient start color is <code>#0762EF</code> (designer spec) while <code>--blue</code> is <code>#0E6EF0</code>.
        Use <code>--horizontal-gradient</code> or <code>--gradient</code> for all primary actions — never use <code>--blue</code> as a fill.
        Reserve <code>--blue</code> for icon tints, active states, and borders.
      </div>

      <!-- Typography -->
      <h3 class="dsSubsectionTitle">Typography</h3>
      <div class="dsTypographyTable">
        <div class="dsTypographyRow">
          <span style="font-size: 20px; font-weight: 700">Heading — 20px Bold</span>
          <code class="dsToken">--font-size-large + --font-weight-bold</code>
        </div>
        <div class="dsTypographyRow">
          <span style="font-size: 20px; font-weight: 500">Subheading — 20px Medium</span>
          <code class="dsToken">--font-size-large + --font-weight-medium</code>
        </div>
        <div class="dsTypographyRow">
          <span style="font-size: 16px; font-weight: 700">Body Bold — 16px Bold</span>
          <code class="dsToken">--font-size-medium + --font-weight-bold</code>
        </div>
        <div class="dsTypographyRow">
          <span style="font-size: 16px; font-weight: 500">Body — 16px Medium</span>
          <code class="dsToken">--font-size-medium + --font-weight-medium</code>
        </div>
        <div class="dsTypographyRow">
          <span style="font-size: 14px; font-weight: 300">Small — 14px Regular</span>
          <code class="dsToken">--font-size-small + --font-weight-regular</code>
        </div>
        <div class="dsTypographyRow">
          <span style="font-size: 12px; font-weight: 300; color: var(--text-gray)">Micro — 12px (labels, timestamps)</span>
          <code class="dsToken">12px (no variable — use inline)</code>
        </div>
      </div>
      <div class="dsNote">
        Font family: <strong>DM Sans</strong>, loaded from Google Fonts and applied globally in <code>base.css</code>.
        Weights available: 300 (Regular), 500 (Medium), 700 (Bold).
      </div>

      <!-- Spacing -->
      <h3 class="dsSubsectionTitle">Spacing</h3>

      <p class="dsSpacingGroupLabel">General use</p>
      <div class="dsSpacingList">
        <div v-for="s in spacingGeneral" :key="s.name" class="dsSpacingRow">
          <code class="dsToken dsSpacingName">{{ s.name }}</code>
          <span class="dsSpacingValue">{{ s.value }}</span>
          <span class="dsColorLabel">{{ s.desc }}</span>
        </div>
      </div>

      <p class="dsSpacingGroupLabel">Derived from <code class="dsToken">--padding</code></p>
      <div class="dsSpacingList">
        <div v-for="s in spacingDerived" :key="s.name" class="dsSpacingRow">
          <code class="dsToken dsSpacingName">{{ s.name }}</code>
          <span class="dsSpacingValue">{{ s.value }}</span>
          <span class="dsColorLabel">{{ s.desc }}</span>
        </div>
      </div>

      <p class="dsSpacingGroupLabel">Component-specific</p>
      <div class="dsSpacingList">
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">--padding-input</code>
          <span class="dsSpacingValue">10px</span>
          <span class="dsColorLabel">Inner padding for <code class="dsToken">&lt;Input&gt;</code> and <code class="dsToken">&lt;Select&gt;</code> only — don't use elsewhere</span>
        </div>
      </div>

      <p class="dsSpacingGroupLabel">Layout</p>
      <div class="dsSpacingList">
        <div v-for="s in spacingLayout" :key="s.name" class="dsSpacingRow">
          <code class="dsToken dsSpacingName">{{ s.name }}</code>
          <span class="dsSpacingValue">{{ s.value }}</span>
          <span class="dsColorLabel">{{ s.desc }}</span>
        </div>
      </div>

      <!-- Border Radius + Shadow -->
      <div class="dsTokenRow">
        <div class="dsTokenCard">
          <h3 class="dsSubsectionTitle dsSubsectionTitleNoMargin">Border Radius</h3>
          <div class="dsRadiusRow">
            <div class="dsRadiusBox" />
            <div>
              <code class="dsToken">--border-radius: 12px</code>
              <p class="dsColorLabel">Cards, buttons, inputs, modals, navigation bar</p>
            </div>
          </div>
          <div class="dsRadiusRow">
            <div class="dsRadiusBox dsRadiusPill" />
            <div>
              <code class="dsToken">border-radius: 50%</code>
              <p class="dsColorLabel">Avatars, FAB, circular action buttons</p>
            </div>
          </div>
        </div>

        <div class="dsTokenCard">
          <h3 class="dsSubsectionTitle dsSubsectionTitleNoMargin">Shadow</h3>
          <div class="dsShadowDemo">
            <div class="dsShadowBox">Shadow preview</div>
            <div>
              <code class="dsToken">--box-shadow</code>
              <p class="dsColorLabel">0px 1px 8px rgba(12, 29, 54, 0.20)</p>
              <p class="dsColorLabel" style="margin-top: 4px">Applied to: cards, buttons (secondary), navigation, modals, FAB</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── BUTTONS ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Buttons</h2>
      <p class="dsSectionDesc">
        Always use <code>&lt;Button&gt;</code> — never raw <code>&lt;button&gt;</code>.
        Buttons are full-width by default; constrain the parent container on desktop.
      </p>

      <h3 class="dsSubsectionTitle">Variants</h3>
      <div class="dsButtonGrid">
        <div class="dsComponentItem">
          <Button title="Primary" theme="primary" />
          <code class="dsToken">theme="primary" (default)</code>
          <p class="dsColorLabel">Main CTA: create, save, submit, sign in</p>
        </div>
        <div class="dsComponentItem">
          <Button title="Secondary" theme="secondary" />
          <code class="dsToken">theme="secondary"</code>
          <p class="dsColorLabel">Alternate or cancel actions</p>
        </div>
        <div class="dsComponentItem">
          <Button title="Going" theme="green" />
          <code class="dsToken">theme="green"</code>
          <p class="dsColorLabel">Confirm, accept, RSVP yes</p>
        </div>
        <div class="dsComponentItem">
          <Button title="Won't Go" theme="red" />
          <code class="dsToken">theme="red"</code>
          <p class="dsColorLabel">Decline, delete, RSVP no</p>
        </div>
        <div class="dsComponentItem">
          <div class="dsGhostPreview">
            <Button title="Send Code" theme="ghost" />
          </div>
          <code class="dsToken">theme="ghost"</code>
          <p class="dsColorLabel">White bg + blue text — use on gradient/coloured backgrounds (e.g. sign-in screen)</p>
        </div>
        <div class="dsComponentItem">
          <Button title="Resend Code" theme="tertiary" />
          <code class="dsToken">theme="tertiary"</code>
          <p class="dsColorLabel">No bg, no border, small muted text — tertiary / least prominent actions (e.g. resend, skip)</p>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">States</h3>
      <div class="dsButtonGrid">
        <div class="dsComponentItem">
          <Button title="Disabled" theme="primary" :disabled="true" />
          <code class="dsToken">:disabled="true"</code>
        </div>
        <div class="dsComponentItem">
          <Button title="Loading" theme="primary" :loading="true" />
          <code class="dsToken">:loading="true"</code>
          <p class="dsColorLabel">Shows spinner, disables button — always use on async actions to prevent double-submit</p>
        </div>
        <div class="dsComponentItem">
          <Button title="Disabled Secondary" theme="secondary" :disabled="true" />
          <code class="dsToken">secondary + disabled</code>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">With icon slot</h3>
      <div class="dsButtonGrid">
        <div class="dsComponentItem">
          <Button title="Create Event" theme="primary">
            <template #icon>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" />
              </svg>
            </template>
          </Button>
          <code class="dsToken">#icon slot</code>
          <p class="dsColorLabel">18×18px slot, left of label</p>
        </div>
      </div>
    </section>

    <!-- ─── FORM INPUTS ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Form Inputs</h2>
      <p class="dsSectionDesc">
        Use <code>&lt;Input&gt;</code> for all text fields. The label floats up on focus/fill —
        don't add a separate <code>&lt;label&gt;</code> element.
      </p>

      <h3 class="dsSubsectionTitle">Text inputs</h3>
      <div class="dsFormGrid">
        <div class="dsComponentItem">
          <Input id="ds-text" type="text" placeholder="Title" v-model="textInput" />
          <code class="dsToken">type="text"</code>
        </div>
        <div class="dsComponentItem">
          <Input id="ds-required" type="text" placeholder="Email" :required="true" v-model="textInput" />
          <code class="dsToken">:required="true"</code>
          <p class="dsColorLabel">Appends * to the label</p>
        </div>
        <div class="dsComponentItem">
          <Input id="ds-error" type="text" placeholder="Email" v-model="errorInput" error="This email is already in use" />
          <code class="dsToken">:error="message"</code>
          <p class="dsColorLabel">Use for 422 validation responses</p>
        </div>
        <div class="dsComponentItem">
          <Input id="ds-multiline" type="text" placeholder="Description" :multiLine="true" v-model="textareaInput" />
          <code class="dsToken">:multiLine="true"</code>
          <p class="dsColorLabel">5-row textarea</p>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Date &amp; time</h3>
      <div class="dsFormGrid">
        <div class="dsComponentItem">
          <Input id="ds-date" type="date" placeholder="Date" v-model="dateInput" />
          <code class="dsToken">type="date"</code>
          <p class="dsColorLabel">Shows text until focused, then native date picker</p>
        </div>
        <div class="dsComponentItem">
          <Input id="ds-time" type="time" placeholder="Time" v-model="timeInput" />
          <code class="dsToken">type="time"</code>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Select</h3>
      <p class="dsSectionDesc">
        Use <code>&lt;Select&gt;</code> for dropdown choices. Pass <code>options</code> as <code>Array&lt;{'{'}id: string, name: string{'}'}&gt;</code>.
      </p>
      <div class="dsFormGrid">
        <div class="dsComponentItem">
          <Select id="ds-select" placeholder="Category" :options="selectOptions" v-model="selectValue" />
          <code class="dsToken">&lt;Select&gt;</code>
        </div>
        <div class="dsComponentItem">
          <Select id="ds-select-loading" placeholder="Loading categories…" :options="[]" :loading="true" v-model="selectLoadingValue" />
          <code class="dsToken">:loading="true"</code>
        </div>
      </div>
    </section>

    <!-- ─── EVENT CARD ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Event Card</h2>
      <p class="dsSectionDesc">
        Use <code>&lt;Event&gt;</code> to render a full event card — illustration top, date badge and info bottom.
        The illustration URL comes from GCS (<code>gs://clubrizer-static-images/event-categories/</code>).
        The date badge background uses <code>category.color</code>.
      </p>
      <div class="dsEventPreview">
        <Event :event="sampleEvent" />
      </div>
      <div class="dsNote">
        Always provide a <code>category.color</code> that contrasts with white text (it's used as the date badge background).
        If no <code>category.picture</code> is set, the image area will be empty — always supply one.
      </div>
    </section>

    <!-- ─── AVATAR ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Avatar</h2>
      <p class="dsSectionDesc">
        Use <code>&lt;Avatar&gt;</code> for all user representations.
        Renders a photo if <code>picture</code> is provided, otherwise falls back to initials on a gradient background.
      </p>

      <h3 class="dsSubsectionTitle">Sizes</h3>
      <div class="dsAvatarRow">
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="sm" />
          <code class="dsToken">size="sm" (32px)</code>
        </div>
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="md" />
          <code class="dsToken">size="md" · default (48px)</code>
        </div>
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="lg" />
          <code class="dsToken">size="lg" (64px)</code>
        </div>
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="xl" />
          <code class="dsToken">size="xl" (80px)</code>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Variants</h3>
      <div class="dsAvatarRow">
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar
            given-name="Kate"
            family-name="Poshuk"
            picture="https://i.pravatar.cc/150?img=47"
            size="lg"
          />
          <code class="dsToken">:picture="url"</code>
          <p class="dsColorLabel">Photo from user's auth provider</p>
        </div>
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="lg" :gradient="true" />
          <code class="dsToken">:gradient="true"</code>
          <p class="dsColorLabel">Gradient ring — used on the profile screen</p>
        </div>
        <div class="dsComponentItem dsComponentItemCenter">
          <Avatar given-name="Kate" family-name="Poshuk" size="lg" label="Kate Poshuk" />
          <code class="dsToken">label="name"</code>
          <p class="dsColorLabel">Tap to reveal — used in event attendee lists</p>
        </div>
      </div>
    </section>

    <!-- ─── NAVIGATION ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Navigation</h2>

      <h3 class="dsSubsectionTitle">Page header bar</h3>
      <p class="dsSectionDesc">
        Use <code>&lt;Header&gt;</code> as the standard top bar for every view that needs a title or back navigation.
        It handles title centering and optional back button — never build a custom title row or back button.
      </p>
      <div class="dsHeaderPreview">
        <Header title="Events" :show-divider="true" />
      </div>
      <div class="dsHeaderPreview" style="margin-top: var(--gap)">
        <Header title="Create Event" left-action="back" />
      </div>
      <div class="dsSpacingList" style="margin-top: 12px">
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">title</code>
          <span class="dsColorLabel">Centered page or step title</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">left-action="back"</code>
          <span class="dsColorLabel">Shows the back arrow — calls <code>router.back()</code> by default</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">:back-fn="fn"</code>
          <span class="dsColorLabel">Override the back action — use when back steps within a multi-step view rather than the browser history (e.g. sign-in OTP flow)</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">:show-divider="true"</code>
          <span class="dsColorLabel">Adds a separator below — use on list/overview screens</span>
        </div>
      </div>
      <div class="dsNote">
        The icon and title inherit <code>color</code> from the parent — wrap in a container with <code>color: var(--white)</code> to use on gradient backgrounds.
      </div>

      <h3 class="dsSubsectionTitle">Mobile — Bottom bar</h3>
      <p class="dsSectionDesc">
        The <code>&lt;Navigation&gt;</code> component is always mounted in <code>App.vue</code>.
        Visibility on mobile is controlled via <code>meta: {'{'} showNavigation: false {'}'}</code> in the router — this adds
        a CSS class that hides it, rather than a <code>v-if</code> (so it can remain visible on desktop even when hidden on mobile).
        Active tab is tracked via <code>meta.activeNav</code> (not Vue Router's built-in <code>activeClass</code>, which doesn't
        activate on child routes) — set <code>activeNav: 'events'</code> or <code>activeNav: 'profile'</code> in the router config
        for each route. On mobile, active uses <code>--blue</code>, inactive uses <code>--gray</code>.
      </p>
      <div class="dsNavPreview">
        <div class="dsMockNav">
          <div class="dsMockNavItem dsMockNavActive">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
            <span class="dsMockNavLabel">Events</span>
          </div>
          <div class="dsMockNavItem">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
            <span class="dsMockNavLabel">Home</span>
          </div>
          <div class="dsMockNavItem">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
            </svg>
            <span class="dsMockNavLabel">Profile</span>
          </div>
        </div>
      </div>
      <div class="dsNote">
        Add a notification dot (<code>position: absolute</code>, 6px circle in <code>--blue</code>) on a tab icon to indicate unread items.
        Set <code>meta: {'{'} showNavigation: false {'}'}</code> for full-screen flows (create event, sign in) — this hides the nav on mobile while keeping it visible on desktop.
      </div>

      <h3 class="dsSubsectionTitle">Desktop — Top bar</h3>
      <p class="dsSectionDesc">
        On desktop (<code>min-width: 768px</code>), the same <code>&lt;Navigation&gt;</code> component reflows into a full-width top bar.
        Icons are hidden, text labels appear, and the active link gets a <code>--blue</code> bottom border instead of a coloured icon.
        Inactive links use <code>--text-color</code> (not <code>--gray</code> as on mobile). The brand name is shown on the left.
        Padding aligns with the content area using <code>max(var(--padding), calc((100% - var(--content-max-width)) / 2))</code>.
      </p>
      <div class="dsDesktopNavPreview">
        <div class="dsMockTopNav">
          <span class="dsMockBrand">{{ $t('team') }}</span>
          <div class="dsMockTopNavLinks">
            <span class="dsMockTopNavLink dsMockNavActive">Events</span>
            <span class="dsMockTopNavLink">Profile</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── FLOATING ACTION BUTTON ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Floating Action Button</h2>
      <p class="dsSectionDesc">
        Use <code>&lt;FloatingActionButton&gt;</code> for primary creation actions.
        It anchors itself bottom-right of its nearest <code>position: relative</code> parent.
        When <code>:actions</code> is provided it expands into a labelled menu on tap.
      </p>
      <div class="dsFabPreview">
        <div class="dsFabItem">
          <p class="dsColorLabel">Simple — emits <code>click</code></p>
          <div class="dsFabBox">
            <div class="dsMockFab">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
            </div>
          </div>
          <code class="dsToken">&lt;FloatingActionButton @click="…" /&gt;</code>
        </div>
        <div class="dsFabItem">
          <p class="dsColorLabel">With actions — opens menu</p>
          <div class="dsFabBox">
            <div class="dsMockFabMenu">
              <span class="dsMockFabLabel">New Event</span>
              <span class="dsMockFabLabel">New Category</span>
            </div>
            <div class="dsMockFab">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
            </div>
          </div>
          <code class="dsToken">:actions="[{'{'}id, label, onClick{'}'}]"</code>
        </div>
        <div class="dsFabItem">
          <p class="dsColorLabel">Loading — blocks interaction</p>
          <div class="dsFabBox">
            <div class="dsMockFab">
              <div class="dsMockFabSpinner" />
            </div>
          </div>
          <code class="dsToken">:loading="true"</code>
        </div>
      </div>
      <div class="dsFabPreview" style="margin-top: 20px">
        <div class="dsFabItem">
          <p class="dsColorLabel">Desktop — extended pill with label</p>
          <div class="dsFabBox" style="width: 160px">
            <div class="dsMockFabPill">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              <span class="dsMockFabPillLabel">New Event</span>
            </div>
          </div>
          <code class="dsToken">:label="$t('events.fab')"</code>
          <p class="dsColorLabel">On <code>min-width: 768px</code>, becomes a pill shape — icon left, text right. Position right-aligns with content area.</p>
        </div>
      </div>
      <div class="dsNote">
        The FAB uses <code>position: absolute</code> so its container needs <code>position: relative</code>.
        Mobile offset: <code>right: var(--padding)</code>, <code>bottom: calc(64px + var(--padding))</code> — sitting above the nav bar.
        Desktop offset: <code>right</code> aligns with the content area edge; <code>bottom: calc(var(--padding) * 2)</code> — no nav bar to clear.
      </div>
    </section>

    <!-- ─── MENU BUTTON ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Menu Button</h2>
      <p class="dsSectionDesc">
        A round overflow button (kebab icon) that opens a small dropdown of actions. Use it to tuck away
        secondary or rarely-needed actions (e.g. deleting an event) instead of showing a prominent button.
        Pass <code>items</code> as <code>{ label, danger?, onClick }</code>; mark destructive actions with
        <code>danger</code>. While open it dims the rest of the screen with a backdrop that blocks other
        clicks; tapping the backdrop, pressing <code>Esc</code>, or choosing an item closes it.
      </p>
      <div class="dsMenuButtonRow">
        <MenuButton :items="menuItems" aria-label="More options" />
        <span class="dsNote">Click it — the menu opens below, aligned to the right.</span>
      </div>
    </section>

    <!-- ─── ALERTS & FEEDBACK ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Alerts &amp; Feedback</h2>
      <p class="dsSectionDesc">
        See <strong>Automatic Request Handling</strong> above for when to use each component.
        <code>&lt;Alert&gt;</code> is for local errors; <code>&lt;RequestError /&gt;</code> handles all other API errors automatically.
      </p>

      <h3 class="dsSubsectionTitle">Alert</h3>
      <p class="dsSectionDesc">
        Local errors only — 422 form validation and client-side errors (e.g. a missing credential).
        <code>variant</code> is <code>error</code> or <code>warning</code>. <code>title</code> is optional.
        <code>size</code> defaults to <code>medium</code>; use <code>small</code> for inline hints (matches small body text).
      </p>
      <div class="dsAlertStack">
        <Alert title="Something went wrong" message="The email address is already in use. Please try a different one." variant="error" />
        <Alert title="Heads up" message="This event is in the past. New attendees won't be able to sign up." variant="warning" />
        <Alert message="No title, small size — for a quiet inline hint next to other content." variant="warning" size="small" />
      </div>

      <h3 class="dsSubsectionTitle">RequestError</h3>
      <p class="dsSectionDesc">
        Place <code>&lt;RequestError /&gt;</code> once per view. It reads from <code>requestStore.error</code> and shows/hides itself automatically — no props or logic required.
      </p>
    </section>

    <!-- ─── MODAL ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Modal</h2>
      <p class="dsSectionDesc">
        Use <code>&lt;Modal&gt;</code> as an overlay wrapper. It provides the dimmed backdrop, gradient border ring, and centered positioning.
        All modal content goes in the default slot — add your own title, body, and close button inside.
      </p>
      <div class="dsModalPreview">
        <div class="dsMockModal">
          <div class="dsMockModalContent">
            <p style="font-weight: var(--font-weight-bold); margin-bottom: 8px">Delete event?</p>
            <p style="font-size: var(--font-size-small); color: var(--text-gray); margin-bottom: 16px">
              This action cannot be undone.
            </p>
            <div style="display: flex; gap: var(--gap)">
              <Button title="Cancel" theme="secondary" />
              <Button title="Delete" theme="red" />
            </div>
          </div>
        </div>
      </div>
      <div class="dsNote">
        The modal has a 1px gradient border ring (diagonal) and animates in (fade + slight rise). It is centered,
        capped at <code>400px</code> wide, and keeps a <code>--padding</code> margin from the screen edges on small
        devices. The backdrop uses <code>--modal-background-color</code> (<code>rgba(12,29,54,0.30)</code>). Always
        give the user a way to dismiss — add a close button in your slot content.
      </div>
    </section>

    <!-- ─── AUTOMATIC REQUEST HANDLING ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Automatic Request Handling</h2>
      <p class="dsSectionDesc">
        Every Axios request is intercepted globally. Loading and errors are handled automatically — you rarely need to wire these up manually.
      </p>

      <h3 class="dsSubsectionTitle">Loading — Top Progress Bar</h3>
      <p class="dsSectionDesc">
        A 3px animated bar fixed at the very top of the viewport. Appears and disappears automatically with every request.
        Mount <code>&lt;TopProgressBar /&gt;</code> once in <code>App.vue</code> — it's already there.
      </p>
      <div class="dsSpacingList">
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">requestStore</code>
          <span class="dsColorLabel">Counter-based Pinia store — interceptors increment on request start, decrement on finish or error</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">isLoading</code>
          <span class="dsColorLabel">True while any request is in flight — the bar animates as long as this is true</span>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Errors — RequestError</h3>
      <p class="dsSectionDesc">
        The Axios interceptor catches all non-401, non-422 errors and stores them in <code>requestStore.error</code>.
        Place <code>&lt;RequestError /&gt;</code> once per view to display it — the component reads from the store and shows/hides itself automatically.
        Errors clear automatically on navigation.
      </p>
      <div class="dsNote">
        <strong>Error decision tree:</strong><br />
        422 form validation → <code>&lt;Alert&gt;</code> locally in the form.<br />
        Client-side error (no request) → <code>&lt;Alert&gt;</code> locally.<br />
        Any other API error → <code>&lt;RequestError /&gt;</code> — interceptor handles it automatically.<br />
        401 Unauthorized → interceptor redirects to sign-in automatically.
      </div>

      <h3 class="dsSubsectionTitle">Local spinners — Buttons &amp; FAB</h3>
      <p class="dsSectionDesc">
        The top bar covers global feedback, but always also set <code>:loading="true"</code> on the triggering button or FAB.
        This disables the control and shows a spinner right at the interaction point — prevents double-submits and gives the user immediate tactile feedback before the bar appears.
      </p>
      <div class="dsButtonGrid" style="max-width: 500px">
        <div class="dsComponentItem">
          <Button title="Creating event…" theme="primary" :loading="true" />
          <code class="dsToken">&lt;Button :loading="isLoading"&gt;</code>
          <p class="dsColorLabel">Use on every async form submit</p>
        </div>
        <div class="dsComponentItem">
          <div class="dsFabBox" style="width: 80px; height: 80px">
            <div class="dsMockFab">
              <div class="dsMockFabSpinner" />
            </div>
          </div>
          <code class="dsToken">&lt;FloatingActionButton :loading="isLoading"&gt;</code>
          <p class="dsColorLabel">Use on async FAB actions — e.g. creating an event from the list screen</p>
        </div>
      </div>
      <div class="dsNote">
        Both the top bar and a local spinner run simultaneously — that's correct and expected.
        The top bar shows system-level progress; the local spinner anchors the feedback to where the user interacted.
      </div>
    </section>

    <!-- ─── OTHER COMPONENTS ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Other Components</h2>

      <h3 class="dsSubsectionTitle">Divider</h3>
      <p class="dsSectionDesc">
        <code>&lt;Divider /&gt;</code> — a 1px full-width horizontal rule at 30% blue opacity.
        Used between sections within a screen.
      </p>
      <div class="dsDividerPreview">
        <Divider />
      </div>

      <h3 class="dsSubsectionTitle">ProfileInfo</h3>
      <p class="dsSectionDesc">
        <code>&lt;ProfileInfo :user="user" /&gt;</code> — avatar (xl, gradient ring) + nickname, full name, and email.
        Used at the top of the profile screen and profile menu.
      </p>
    </section>

    <!-- ─── LAYOUT PATTERNS ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Layout Patterns</h2>
      <p class="dsSectionDesc">
        The app is mobile-first. All views are rendered inside a flex column that fills the viewport height.
        On desktop, the same structure widens naturally — add breakpoints per-view as needed.
      </p>

      <h3 class="dsSubsectionTitle">Mobile shell (375px)</h3>
      <div class="dsLayoutGrid">
        <div class="dsMockPhone">
          <div class="dsMockPhoneBar">TopProgressBar (fixed, z-index 9999)</div>
          <div class="dsMockPhoneContent">
            <div class="dsMockPhoneView">RouterView · overflow-y: auto · padding: 12px</div>
          </div>
          <div class="dsMockPhoneNav">Navigation (when meta.showNavigation)</div>
        </div>
        <div class="dsLayoutNotes">
          <p><code>.app</code> — flex column, <code>height: 100dvh</code>, <code>overflow: hidden</code></p>
          <p><code>.content</code> — <code>flex-grow: 1</code>, <code>overflow-y: auto</code>, <code>padding: 12px</code></p>
          <p><code>.navigation</code> — <code>flex-shrink: 0</code>, margins on all sides</p>
          <p>Horizontal page padding: <code>--padding (12px)</code></p>
          <p>Navigation bar height: ~57px + 12px margin = 69px from bottom</p>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Desktop (768px+)</h3>
      <div class="dsLayoutGrid">
        <div class="dsMockDesktop">
          <div class="dsMockDesktopNav">Navigation (top bar, full width)</div>
          <div class="dsMockDesktopBody">
            <div class="dsMockDesktopContent">RouterView · max-width: 1100px · margin: 0 auto · padding: 12px</div>
          </div>
        </div>
        <div class="dsLayoutNotes">
          <p><code>.app</code> — flex column, <code>height: 100dvh</code></p>
          <p><code>.navigation</code> — <code>order: -1</code> (floats to top), full width, top bar layout</p>
          <p><code>.content</code> — <code>flex-grow: 1</code>, <code>overflow-y: auto</code></p>
          <p>Content max-width: <code>--content-max-width (1100px)</code>, centered with <code>margin: 0 auto</code></p>
          <p>Horizontal page padding: <code>--padding (12px)</code></p>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">Desktop card pattern</h3>
      <p class="dsSectionDesc">
        On desktop, form and detail views wrap their content in a centered card.
        Apply to the inner container — not the outer view wrapper (which handles padding and centering).
      </p>
      <div class="dsSpacingList" style="margin-bottom: 12px">
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">max-width</code>
          <span class="dsColorLabel"><code>var(--content-max-width)</code> — consistent width across all detail/form screens</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">border-radius</code>
          <span class="dsColorLabel"><code>var(--border-radius)</code> on the card wrapper (partial on sub-elements if needed, e.g. hero image top corners)</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">box-shadow</code>
          <span class="dsColorLabel"><code>var(--box-shadow)</code> — lifts card off the background</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">background</code>
          <span class="dsColorLabel"><code>var(--background-color)</code> — matches page bg on mobile, looks like a card on desktop</span>
        </div>
      </div>
      <div class="dsNote">
        Views currently using this pattern: EventDetailView, NewEventView, ProfileView, ProfileSetupView, PendingApprovalView.
        Apply the same pattern to any new detail or form screen.
      </div>
    </section>

    <!-- ─── VOICE & TONE ─── -->
    <section class="dsSection">
      <h2 class="dsSectionTitle">Voice &amp; Tone</h2>
      <p class="dsSectionDesc">
        Clubrizer is a sports app. The people using it are teammates, coaches, and club members — not enterprise users filling out forms.
        Every string should sound like it was written by a person on the team, not generated by a system.
      </p>

      <h3 class="dsSubsectionTitle">Principles</h3>
      <div class="dsSpacingList">
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">Friendly</code>
          <span class="dsColorLabel">Warm and approachable — like a message from a teammate, not a terms-of-service page</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">Personal</code>
          <span class="dsColorLabel">Always "du", never "Sie" or impersonal constructions. The user is a person, address them directly.</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">Authentic</code>
          <span class="dsColorLabel">Say what you mean, plainly. No corporate speak, no filler words, no over-explanation.</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">Relaxed</code>
          <span class="dsColorLabel">Sports context — keep it easy and light. Short sentences. No drama, even for errors.</span>
        </div>
        <div class="dsSpacingRow">
          <code class="dsToken dsSpacingName">No emojis</code>
          <span class="dsColorLabel">Ever. The design carries the personality — the text doesn't need decoration.</span>
        </div>
      </div>

      <h3 class="dsSubsectionTitle">By context</h3>
      <div class="dsToneGrid">
        <div class="dsToneCard">
          <p class="dsToneContext">Buttons &amp; actions</p>
          <p class="dsToneRule">Short imperative verbs. Say what happens, nothing more.</p>
          <div class="dsToneExamples">
            <span class="dsToneBad">Event erstellen</span>
            <span class="dsToneGood">Erstellen</span>
            <span class="dsToneBad">Jetzt anmelden</span>
            <span class="dsToneGood">Anmelden</span>
          </div>
        </div>
        <div class="dsToneCard">
          <p class="dsToneContext">Errors</p>
          <p class="dsToneRule">Own it briefly, tell them what to do. Never blame the user, never be dramatic.</p>
          <div class="dsToneExamples">
            <span class="dsToneBad">Ein unerwarteter Fehler ist aufgetreten.</span>
            <span class="dsToneGood">Etwas hat nicht geklappt. Bitte versuch es nochmal.</span>
          </div>
        </div>
        <div class="dsToneCard">
          <p class="dsToneContext">Validation</p>
          <p class="dsToneRule">Helpful, not accusatory. Tell them what's needed, not what they did wrong.</p>
          <div class="dsToneExamples">
            <span class="dsToneBad">Das Pflichtfeld darf nicht leer sein.</span>
            <span class="dsToneGood">Bitte gib deinen Vornamen ein.</span>
          </div>
        </div>
        <div class="dsToneCard">
          <p class="dsToneContext">Status &amp; waiting</p>
          <p class="dsToneRule">Keep people calm and in the loop. Avoid bureaucratic nouns.</p>
          <div class="dsToneExamples">
            <span class="dsToneBad">Ausstehende Genehmigung</span>
            <span class="dsToneGood">Dein Account wird geprüft</span>
            <span class="dsToneBad">Sitzung abgelaufen</span>
            <span class="dsToneGood">Meld dich kurz neu an</span>
          </div>
        </div>
      </div>

      <div class="dsNote" style="margin-top: 20px">
        When in doubt: say it out loud. If it sounds like something you'd actually say to a teammate before training, it's probably right.
        If it sounds like a bank letter, rewrite it.
      </div>
    </section>

    <!-- ─── USAGE GUIDELINES ─── -->
    <section class="dsSection dsLastSection">
      <h2 class="dsSectionTitle">Usage Guidelines</h2>
      <div class="dsGuidelinesGrid">
        <div class="dsGuidelineCard dsGuidelineDo">
          <h4 class="dsGuidelineTitle">✓ Do</h4>
          <ul class="dsGuidelineList">
            <li>Use CSS variables from <code>base.css</code> for all colors, sizes, and spacing</li>
            <li>Use <code>&lt;Button&gt;</code> instead of <code>&lt;button&gt;</code></li>
            <li>Use <code>:loading="true"</code> on any button or FAB that triggers an async action</li>
            <li>Place <code>&lt;RequestError /&gt;</code> once per view — the interceptor populates it automatically</li>
            <li>Use <code>&lt;Header&gt;</code> for all page titles and back navigation</li>
            <li>Use scoped CSS with camelCase class names prefixed by component name</li>
            <li>Use <code>&lt;Input&gt;</code> floating labels — no separate <code>&lt;label&gt;</code> needed</li>
            <li>Design mobile-first, then extend with <code>@media (min-width: 768px)</code></li>
            <li>Use <code>--content-max-width</code> for centering content on desktop — define it once, use everywhere</li>
            <li>Set <code>meta.activeNav</code> in the router for each route to drive nav active state</li>
            <li>Enforce access control in the backend — frontend-only restrictions are not secure</li>
            <li>Internationalize all user-facing strings with <code>i18n.global.t()</code></li>
          </ul>
        </div>
        <div class="dsGuidelineCard dsGuidelineDont">
          <h4 class="dsGuidelineTitle">✗ Don't</h4>
          <ul class="dsGuidelineList">
            <li>Hardcode color hex values, pixel sizes, or spacing — use variables</li>
            <li>Add per-component loading indicators for API requests other than button/FAB spinners</li>
            <li>Use <code>&lt;Alert&gt;</code> for general API errors — use <code>&lt;RequestError /&gt;</code> instead</li>
            <li>Use native <code>&lt;select&gt;</code> or <code>&lt;button&gt;</code> directly</li>
            <li>Skip the <code>:loading</code> prop — users can double-submit without it</li>
            <li>Hide UI elements on the frontend as a security measure — backend must also check</li>
            <li>Hardcode user-facing strings in templates or scripts</li>
          </ul>
        </div>
      </div>
    </section>

  </div>
</template>

<style scoped>
.dsPage {
  max-width: 1100px;
  margin: 0 auto;
  padding-bottom: 48px;
}

/* ── Page header ── */
.dsPageHeader {
  text-align: center;
  padding: 32px 0 24px;
  border-bottom: 1px solid var(--gray);
  margin-bottom: 40px;
}

.dsPageTitle {
  font-size: 28px;
  font-weight: var(--font-weight-bold);
  background: var(--horizontal-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px;
}

.dsPageSubtitle {
  color: var(--text-gray);
  font-size: var(--font-size-small);
  max-width: 560px;
  margin: 0 auto;
}

/* ── Sections ── */
.dsSection {
  margin-bottom: 48px;
  padding-bottom: 48px;
  border-bottom: 1px solid var(--gray);
}

.dsLastSection {
  border-bottom: none;
}

.dsSectionTitle {
  font-size: 22px;
  font-weight: var(--font-weight-bold);
  color: var(--blue);
  margin: 0 0 6px;
}

.dsSectionDesc {
  color: var(--text-gray);
  font-size: var(--font-size-small);
  margin: 0 0 20px;
  line-height: 1.5;
}

.dsSubsectionTitle {
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  color: var(--text-light);
  margin: 28px 0 12px;
}

.dsSubsectionTitleNoMargin {
  margin-top: 0;
}

/* ── Colors ── */
.dsColorGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: var(--padding);
}

.dsColorSwatch {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dsColorBlock {
  height: 56px;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.dsColorHex {
  font-size: 11px;
  font-weight: var(--font-weight-medium);
  color: var(--text-light);
  font-family: monospace;
}

.dsColorLabel {
  font-size: 11px;
  color: var(--text-gray);
  line-height: 1.4;
  margin: 0;
}

/* ── Typography ── */
.dsTypographyTable {
  display: flex;
  flex-direction: column;
}

.dsTypographyRow {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: var(--gap);
  padding: 10px 0;
  border-bottom: 1px solid var(--light-gray);
  color: var(--text-light);
}

/* ── Spacing ── */
.dsSpacingGroupLabel {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  color: var(--text-gray);
  margin: 16px 0 6px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dsSpacingList {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin-bottom: 4px;
}

.dsSpacingRow {
  display: flex;
  align-items: center;
  gap: var(--padding);
  flex-wrap: wrap;
}

.dsSpacingName {
  width: 150px;
  flex-shrink: 0;
}

.dsSpacingValue {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  color: var(--text-light);
  width: 40px;
  flex-shrink: 0;
}

/* ── Radius + Shadow ── */
.dsTokenRow {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--padding);
  margin-top: 24px;
}

@media (max-width: 600px) {
  .dsTokenRow {
    grid-template-columns: 1fr;
  }
}

.dsTokenCard {
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
}

.dsRadiusRow {
  display: flex;
  align-items: center;
  gap: var(--padding);
  margin-bottom: var(--padding);
}

.dsRadiusBox {
  width: 48px;
  height: 48px;
  background: var(--horizontal-gradient);
  border-radius: var(--border-radius);
  flex-shrink: 0;
}

.dsRadiusPill {
  border-radius: 50%;
}

.dsShadowDemo {
  display: flex;
  align-items: flex-start;
  gap: var(--padding);
}

.dsShadowBox {
  padding: 10px 14px;
  background: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  font-size: var(--font-size-small);
  white-space: nowrap;
  flex-shrink: 0;
}

/* ── Ghost button preview ── */
.dsGhostPreview {
  background: var(--gradient);
  border-radius: var(--border-radius);
  padding: var(--padding);
}

/* ── Buttons ── */
.dsButtonGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--padding);
  margin-bottom: 8px;
}

.dsComponentItem {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dsComponentItemCenter {
  align-items: center;
}

/* ── Forms ── */
.dsFormGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: var(--padding);
  margin-bottom: 8px;
}

/* ── Event ── */
.dsEventPreview {
  max-width: 360px;
}

/* ── Avatar ── */
.dsAvatarRow {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  align-items: flex-end;
  margin-bottom: 8px;
}

/* ── Navigation ── */
.dsNavPreview {
  max-width: 375px;
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
  margin-bottom: 12px;
}

.dsMockNav {
  display: flex;
  justify-content: space-evenly;
  background: var(--white);
  border-radius: var(--border-radius);
  padding: 12px 0;
  box-shadow: var(--box-shadow);
}

.dsMockNavItem {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  color: var(--gray);
  cursor: default;
}

.dsMockNavActive {
  color: var(--blue);
}

.dsMockNavLabel {
  font-size: 10px;
}

.dsDesktopNavPreview {
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
  margin-bottom: 12px;
}

.dsMockTopNav {
  display: flex;
  align-items: center;
  gap: 24px;
  background: var(--white);
  border-radius: var(--border-radius);
  padding: 10px 20px;
  box-shadow: var(--box-shadow);
}

.dsMockBrand {
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  background: var(--horizontal-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.dsMockTopNavLinks {
  display: flex;
  gap: 24px;
  margin-left: auto;
}

.dsMockTopNavLink {
  font-size: var(--font-size-small);
  color: var(--text-gray);
  cursor: default;
}

.dsMockTopNavLink.dsMockNavActive {
  color: var(--blue);
  font-weight: var(--font-weight-medium);
}

/* ── FAB ── */
.dsFabPreview {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  padding: var(--padding);
  background: var(--light-gray);
  border-radius: var(--border-radius);
  margin-bottom: 12px;
}

.dsFabItem {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.dsFabBox {
  position: relative;
  width: 120px;
  height: 100px;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
}

.dsMockFab {
  width: 56px;
  height: 56px;
  background: var(--gradient);
  border-radius: 50%;
  box-shadow: var(--box-shadow);
  display: flex;
  align-items: center;
  justify-content: center;
}

.dsMockFabMenu {
  position: absolute;
  bottom: calc(100% - 36px);
  right: 0;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  align-items: flex-end;
}

.dsMockFabLabel {
  background: var(--white);
  padding: 8px 14px;
  border-radius: 24px;
  border: 2px solid var(--blue);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  box-shadow: var(--box-shadow);
  white-space: nowrap;
}

.dsMockFabSpinner {
  width: 22px;
  height: 22px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: dsSpin 1s linear infinite;
}

@keyframes dsSpin {
  to { transform: rotate(360deg); }
}

.dsMockFabPill {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--gradient);
  border-radius: 28px;
  padding: 0 20px 0 12px;
  height: 56px;
  box-shadow: var(--box-shadow);
  white-space: nowrap;
}

.dsMockFabPillLabel {
  color: white;
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
}

/* ── Alerts ── */
.dsAlertStack {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  max-width: 500px;
  margin-bottom: 12px;
}

.dsMenuButtonRow {
  display: flex;
  align-items: center;
  gap: var(--padding);
}

/* ── Modal ── */
.dsModalPreview {
  background: rgba(12, 29, 54, 0.15);
  border-radius: var(--border-radius);
  padding: 48px 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.dsMockModal {
  padding: 1px;
  background: var(--gradient);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  max-width: 320px;
  width: 100%;
}

.dsMockModalContent {
  padding: calc(var(--padding) * 2);
  background: var(--white);
  border-radius: calc(var(--border-radius) - 1px);
}

/* ── Header ── */
.dsHeaderPreview {
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
}

/* ── Divider ── */
.dsDividerPreview {
  padding: var(--padding) 0;
}

/* ── Layout ── */
.dsLayoutGrid {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 24px;
  align-items: start;
  margin-bottom: 12px;
}

@media (max-width: 600px) {
  .dsLayoutGrid {
    grid-template-columns: 1fr;
  }
}

.dsMockPhone {
  width: 160px;
  height: 280px;
  border: 2px solid var(--gray);
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  font-size: 10px;
  flex-shrink: 0;
}

.dsMockPhoneBar {
  background: var(--horizontal-gradient);
  color: white;
  padding: 4px 8px;
  text-align: center;
  font-size: 9px;
}

.dsMockPhoneContent {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dsMockPhoneView {
  flex: 1;
  padding: 8px;
  background: var(--light-gray);
  color: var(--text-gray);
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.dsMockPhoneNav {
  background: var(--white);
  border-top: 1px solid var(--gray);
  padding: 6px 8px;
  text-align: center;
  color: var(--text-gray);
}

.dsMockDesktop {
  width: 260px;
  height: 160px;
  border: 2px solid var(--gray);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  font-size: 9px;
  flex-shrink: 0;
}

.dsMockDesktopNav {
  background: var(--white);
  border-bottom: 1px solid var(--gray);
  padding: 5px 8px;
  text-align: center;
  color: var(--text-gray);
  flex-shrink: 0;
}

.dsMockDesktopBody {
  flex: 1;
  background: var(--light-gray);
  display: flex;
  padding: 8px;
}

.dsMockDesktopContent {
  flex: 1;
  background: var(--white);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: var(--text-gray);
  padding: 4px;
}

.dsLayoutNotes {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: var(--font-size-small);
  color: var(--text-light);
}

.dsLayoutNotes p {
  margin: 0;
}

/* ── Guidelines ── */
.dsGuidelinesGrid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--padding);
}

@media (max-width: 600px) {
  .dsGuidelinesGrid {
    grid-template-columns: 1fr;
  }
}

.dsGuidelineCard {
  padding: var(--padding);
  border-radius: var(--border-radius);
  border-left: 4px solid;
}

.dsGuidelineDo {
  background: rgba(40, 232, 162, 0.08);
  border-color: var(--green);
}

.dsGuidelineDont {
  background: rgba(243, 103, 94, 0.08);
  border-color: var(--red);
}

.dsGuidelineTitle {
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--gap);
  color: var(--text-light);
}

.dsGuidelineList {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: var(--font-size-small);
  color: var(--text-light);
  line-height: 1.4;
}

/* ── Voice & Tone ── */
.dsToneGrid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--padding);
  margin-top: 12px;
}

.dsToneCard {
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.dsToneContext {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  color: var(--text-light);
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dsToneRule {
  font-size: var(--font-size-small);
  color: var(--text-gray);
  margin: 0;
  line-height: 1.4;
}

.dsToneExamples {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 4px;
}

.dsToneBad,
.dsToneGood {
  font-size: var(--font-size-small);
  padding: 6px var(--padding);
  border-radius: 8px;
  line-height: 1.4;
}

.dsToneBad {
  background: rgba(243, 103, 94, 0.08);
  color: var(--text-light);
  text-decoration: line-through;
  text-decoration-color: var(--red);
}

.dsToneGood {
  background: rgba(40, 232, 162, 0.1);
  color: var(--text-light);
  font-weight: var(--font-weight-medium);
}

/* ── Shared helpers ── */
.dsToken {
  font-family: monospace;
  font-size: 11px;
  color: var(--text-gray);
  background: var(--light-gray);
  padding: 2px 5px;
  border-radius: 4px;
  white-space: nowrap;
}

.dsNote {
  background: var(--light-gray);
  border-left: 3px solid var(--blue);
  padding: var(--gap) var(--padding);
  border-radius: 0 var(--border-radius) var(--border-radius) 0;
  font-size: var(--font-size-small);
  color: var(--text-light);
  margin-top: var(--gap);
  line-height: 1.5;
}
</style>
