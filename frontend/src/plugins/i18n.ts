import { createI18n } from 'vue-i18n'

export default createI18n({
  locale: navigator.language,
  fallbackLocale: 'en',
  messages: {
    en: {
      team: 'LISC-2010', // TODO move to a config or env variable
      signIn: {
        welcomeTo: 'Welcome to',
        getStarted: 'Sign in to get started',
        failedToSignIn: 'Failed to sign in',
        failedToSignInWithGoogle: 'Failed to sign in with Google',
        failedToParseGoogleToken: 'Failed to parse Google login token'
      },
      profile: {
        header: 'My Profile',
        logout: 'Logout'
      },
      events: {
        header: 'Events',
        development: {
          title: 'Development Preview',
          message: 'This app is still in development. Please reach out to Kathi if you notice any bugs.'
        },
        new: {
          header: 'Create Event',
          title: 'Title',
          date: 'Date',
          time: 'Time',
          location: 'Location',
          category: 'Category',
          description: 'Description',
          create: 'Create',
          failedToCreate: 'Failed to create event',
          required: 'Required',
          pastDate: 'Date cannot be in the past'
        }
      },
      select: {
        loading: 'Loading...'
      }
    },
    de: {
      team: 'LISC-2010',
      signIn: {
        welcomeTo: 'Willkommen bei',
        getStarted: 'Melde dich an, um loszulegen',
        failedToSignIn: 'Anmeldung fehlgeschlagen',
        failedToSignInWithGoogle: 'Google Anmeldung fehlgeschlagen',
        failedToParseGoogleToken: 'Google-Anmelde-Token konnte nicht verarbeitet werden'
      },
      profile: {
        header: 'Mein Profil',
        logout: 'Abmelden'
      },
      events: {
        header: 'Events',
        development: {
          title: 'Entwicklungs-Vorschau',
          message: 'Diese App ist in der Entwicklung. Bitte melde dich bei Kathi, wenn dir Fehler auffallen.'
        },
        new: {
          header: 'Neues Event',
          title: 'Name',
          date: 'Datum',
          time: 'Uhrzeit',
          location: 'Ort',
          category: 'Kategorie',
          description: 'Beschreibung',
          create: 'Speichern',
          failedToCreate: 'Event konnte nicht erstellt werden',
          required: 'Pflichtfeld',
          pastDate: 'Datum darf nicht in der Vergangenheit liegen'
        }
      }
    }
  }
})
