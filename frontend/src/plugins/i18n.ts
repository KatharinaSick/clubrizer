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
        back: 'Back',
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
          required: 'Required',
          pastDate: 'Date cannot be in the past'
        },
        detail: {
          wontGo: 'Won\'t Go',
          going: 'Going',
          at: 'at',
          attendees: {
            going: '{count} going',
            notGoing: '{count} not going',
            noResponses: 'No responses yet',
          }
        }
      },
      select: {
        loading: 'Loading...'
      },
      pendingApproval: {
        pending: {
          title: 'Approval Pending',
          message: 'Your account is pending approval. Hang tight — an admin will review it shortly.'
        },
        rejected: {
          title: 'Account Declined',
          message: 'Your account request was declined. If you think this was a mistake, please reach out directly to the club.'
        }
      },
      request: {
        errorTitle: 'Something went wrong',
        unexpectedError: 'An unexpected error occurred. Please try again.'
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
        back: 'Zurück',
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
          required: 'Pflichtfeld',
          pastDate: 'Datum darf nicht in der Vergangenheit liegen'
        },
        detail: {
          wontGo: 'Absagen',
          going: 'Zusagen',
          at: 'um',
          attendees: {
            going: '{count} Zusagen',
            notGoing: '{count} Absagen',
            noResponses: 'Noch keine Rückmeldungen',
          }
        }
      },
      pendingApproval: {
        pending: {
          title: 'Ausstehende Genehmigung',
          message: 'Dein Account wartet auf Genehmigung. Bleib dran – ein Administrator wird es in Kürze prüfen.'
        },
        rejected: {
          title: 'Account abgelehnt',
          message: 'Deine Account Erstellung wurde abgelehnt. Wenn du glaubst, dass das ein Fehler ist, wende dich direkt an den Verein.'
        }
      },
      request: {
        errorTitle: 'Etwas ist schiefgelaufen',
        unexpectedError: 'Ein unerwarteter Fehler ist aufgetreten. Bitte versuche es erneut.'
      }
    }
  }
})
