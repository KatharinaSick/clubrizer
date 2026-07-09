<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/Button.vue'
import Divider from '@/components/Divider.vue'
import Alert from '@/components/Alert.vue'

const router = useRouter()
const auth = useAuthStore()

interface JoinStep {
  title: string
  text: string
  linkTo?: string
  linkLabel?: string
  note?: string
}

const joinSteps: JoinStep[] = [
  { title: 'gettingStarted.join.step1Title', text: 'gettingStarted.join.step1', linkTo: '/signin', linkLabel: 'gettingStarted.join.step1Link' },
  { title: 'gettingStarted.join.step2Title', text: 'gettingStarted.join.step2', note: 'gettingStarted.join.step2Note' },
  { title: 'gettingStarted.join.step3Title', text: 'gettingStarted.join.step3' },
  { title: 'gettingStarted.join.step4Title', text: 'gettingStarted.join.step4' },
]

interface InstallPlatform {
  title: string
  steps: string[]
}

const installPlatforms: InstallPlatform[] = [
  {
    title: 'gettingStarted.install.iphoneTitle',
    steps: [
      'gettingStarted.install.iphoneStep1',
      'gettingStarted.install.iphoneStep2',
      'gettingStarted.install.iphoneStep3',
      'gettingStarted.install.iphoneStep4',
    ],
  },
  {
    title: 'gettingStarted.install.androidTitle',
    steps: [
      'gettingStarted.install.androidStep1',
      'gettingStarted.install.androidStep2',
      'gettingStarted.install.androidStep3',
      'gettingStarted.install.androidStep4',
    ],
  },
]

interface WhySource {
  label: string
  url: string
}

// Background reading for curious users about the Android install warning.
const whySources: WhySource[] = [
  { label: 'gettingStarted.install.whySource1', url: 'https://developers.google.com/android/play-protect/warning-dev-guidance' },
]
</script>

<template>
  <div class="gettingStartedPage">

    <!-- Hero -->
    <header class="gettingStartedHero">
      <p class="gettingStartedWelcomeTo">{{ $t('gettingStarted.hero.welcomeTo') }}</p>
      <h1 class="gettingStartedTitle">{{ $t('team') }}</h1>
      <p class="gettingStartedIntro">{{ $t('gettingStarted.hero.intro') }}</p>
    </header>

    <!-- How to join -->
    <section class="gettingStartedSection">
      <h2 class="gettingStartedSectionTitle">{{ $t('gettingStarted.join.title') }}</h2>
      <ol class="gettingStartedSteps">
        <li v-for="(step, index) in joinSteps" :key="index" class="gettingStartedStep">
          <span class="gettingStartedStepNumber">{{ index + 1 }}</span>
          <div class="gettingStartedStepBody">
            <h3 class="gettingStartedStepTitle">{{ $t(step.title) }}</h3>
            <i18n-t
              v-if="step.linkTo && step.linkLabel"
              :keypath="step.text"
              tag="p"
              class="gettingStartedStepText"
            >
              <template #link>
                <RouterLink :to="step.linkTo" class="gettingStartedStepLink">{{ $t(step.linkLabel) }}</RouterLink>
              </template>
            </i18n-t>
            <p v-else class="gettingStartedStepText">{{ $t(step.text) }}</p>
            <Alert
              v-if="step.note"
              variant="warning"
              size="small"
              :message="$t(step.note)"
              class="gettingStartedStepNote"
            />
          </div>
        </li>
      </ol>
      <div class="gettingStartedCta">
        <Button
          v-if="auth.isLoggedIn"
          :title="$t('gettingStarted.join.ctaHome')"
          @click="router.push('/events')"
        />
        <Button
          v-else
          :title="$t('gettingStarted.join.cta')"
          @click="router.push('/signin')"
        />
      </div>
    </section>

    <Divider />

    <!-- Add it to your phone -->
    <section class="gettingStartedSection">
      <h2 class="gettingStartedSectionTitle">{{ $t('gettingStarted.install.title') }}</h2>
      <p class="gettingStartedSectionIntro">{{ $t('gettingStarted.install.intro') }}</p>
      <div class="gettingStartedPlatforms">
        <div
          v-for="platform in installPlatforms"
          :key="platform.title"
          class="gettingStartedPlatform"
        >
          <h3 class="gettingStartedPlatformTitle">{{ $t(platform.title) }}</h3>
          <ol class="gettingStartedPlatformSteps">
            <li
              v-for="stepKey in platform.steps"
              :key="stepKey"
              class="gettingStartedPlatformText"
            >
              {{ $t(stepKey) }}
            </li>
          </ol>
        </div>
      </div>

      <Alert
        variant="info"
        size="small"
        :title="$t('gettingStarted.install.androidNoteTitle')"
        :message="$t('gettingStarted.install.androidNote')"
        class="gettingStartedInstallNote"
      >
        <details class="gettingStartedWhy">
          <summary class="gettingStartedWhySummary">{{ $t('gettingStarted.install.whyTitle') }}</summary>
          <p class="gettingStartedWhyText">{{ $t('gettingStarted.install.whyText') }}</p>
          <ul class="gettingStartedWhySources">
            <li v-for="source in whySources" :key="source.url">
              <a
                :href="source.url"
                target="_blank"
                rel="noopener noreferrer"
                class="gettingStartedWhyLink"
              >{{ $t(source.label) }}</a>
            </li>
          </ul>
        </details>
      </Alert>
    </section>

    <p class="gettingStartedFooter">{{ $t('gettingStarted.footer') }}</p>
  </div>
</template>

<style scoped>
.gettingStartedPage {
  width: 100%;
  max-width: 640px;
  margin: 0 auto;
  padding: var(--padding) 0 48px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 48px;
}

/* Hero */
.gettingStartedHero {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.gettingStartedWelcomeTo {
  margin: 0;
  color: var(--text-gray);
  font-size: var(--font-size-medium);
}

.gettingStartedTitle {
  margin: 0;
  font-size: 32px;
  font-weight: var(--font-weight-bold);
  background: var(--horizontal-gradient);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: var(--blue);
  line-height: 1.2;
}

.gettingStartedIntro {
  margin: var(--gap) 0 0;
  font-size: var(--font-size-medium);
  line-height: 1.5;
}

/* Sections */
.gettingStartedSection {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.gettingStartedSectionTitle {
  margin: 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.gettingStartedSectionIntro {
  margin: -12px 0 0;
  color: var(--text-light);
  font-size: var(--font-size-small);
  line-height: 1.5;
}

/* Join steps */
.gettingStartedSteps {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.gettingStartedStep {
  display: flex;
  align-items: flex-start;
  gap: var(--padding);
}

.gettingStartedStepNumber {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--gradient);
  color: var(--white);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

.gettingStartedStepBody {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 4px;
}

.gettingStartedStepTitle {
  margin: 0;
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-medium);
}

.gettingStartedStepText {
  margin: 0;
  color: var(--text-light);
  font-size: var(--font-size-small);
  line-height: 1.5;
}

.gettingStartedStepNote {
  margin-top: var(--gap);
}

.gettingStartedStepLink {
  color: var(--blue);
  font-weight: var(--font-weight-medium);
  text-decoration: underline;
}

.gettingStartedStepLink:active {
  opacity: 0.7;
}

.gettingStartedCta {
  margin-top: var(--gap);
}

/* Install */
.gettingStartedPlatforms {
  display: flex;
  flex-direction: column;
  gap: var(--padding);
}

.gettingStartedPlatform {
  background: var(--light-gray);
  border-radius: var(--border-radius);
  padding: var(--padding);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.gettingStartedPlatformTitle {
  margin: 0;
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
}

.gettingStartedPlatformSteps {
  margin: var(--gap) 0 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.gettingStartedPlatformText {
  margin: 0;
  padding-left: 4px;
  color: var(--text-light);
  font-size: var(--font-size-small);
  line-height: 1.5;
}

/* Why disclosure */
.gettingStartedWhySummary {
  cursor: pointer;
  color: var(--blue);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
}

.gettingStartedWhySummary:active {
  opacity: 0.7;
}

.gettingStartedWhyText {
  margin: var(--padding) 0 0;
  color: var(--text-light);
  font-size: var(--font-size-small);
  line-height: 1.5;
}

.gettingStartedWhySources {
  margin: var(--gap) 0 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.gettingStartedWhyLink {
  color: var(--blue);
  font-size: var(--font-size-small);
  text-decoration: underline;
}

.gettingStartedWhyLink:active {
  opacity: 0.7;
}

/* Footer */
.gettingStartedFooter {
  margin: 0;
  text-align: center;
  color: var(--text-gray);
  font-size: var(--font-size-small);
  line-height: 1.5;
}

@media (min-width: 768px) {
  .gettingStartedPlatforms {
    flex-direction: row;
  }

  .gettingStartedPlatform {
    flex: 1;
  }

  .gettingStartedCta {
    display: flex;
    justify-content: center;
  }

  .gettingStartedCta :deep(.button) {
    width: auto;
    padding-inline: calc(var(--padding) * 3);
  }
}
</style>
