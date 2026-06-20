import { createI18n } from 'vue-i18n'

export default createI18n({
  locale: navigator.language.split('-')[0],
  fallbackLocale: 'en',
  messages: {
    en: {
      team: 'LISC-2010', // TODO move to a config or env variable
      signIn: {
        welcomeTo: 'Welcome to',
        getStarted: 'Sign in to get started',
        emailLabel: 'Email address',
        emailNote: 'Please use the email address you registered with at your club.',
        sendCode: 'Send Code',
        codeTitle: 'Check your email!',
        codeNote: 'We sent a 6-digit code to {email}. Enter it below. If you don\'t see it, check your spam folder.',
        codeLabel: '6-digit code',
        verify: 'Verify',
        backToEmail: 'Back',
        resendCode: 'Resend Code',
        codeSent: 'A new code was sent to your email.',
        emailRequired: 'Please enter your email address',
        codeInvalid: 'Please enter the 6-digit code from your email',
      },
      profile: {
        header: 'My Profile',
        logout: 'Logout',
        edit: 'Edit Profile',
        cancel: 'Cancel',
      },
      profileSetup: {
        header: 'Set Up Your Profile',
        firstName: 'First Name',
        lastName: 'Last Name',
        nickName: 'Nickname (optional)',
        picture: 'Profile Picture (optional)',
        save: 'Save',
        firstNameRequired: 'Please enter your first name',
        lastNameRequired: 'Please enter your last name',
      },
      events: {
        header: 'Events',
        back: 'Back',
        noEvents: 'No events yet.',
        fab: 'New Event',
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
          pastDate: 'Date cannot be in the past',
          timeInvalid: 'Please enter a time in HH:MM format'
        },
        createdBy: 'by',
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
        loading: 'Loading…'
      },
      pendingApproval: {
        pending: {
          title: 'Approval Pending',
          message: 'Your account is pending approval. Hang tight — an admin will review it shortly.',
          checkStatus: 'Check Status',
          stillPending: 'Still pending — check back a little later.'
        },
        rejected: {
          title: 'Account Declined',
          message: 'Your account request was declined. If you think this was a mistake, please reach out directly to the club.'
        },
        cancel: 'Cancel'
      },
      navigation: {
        events: 'Events',
        profile: 'Profile',
      },
      request: {
        errorTitle: 'Something went wrong',
        unexpectedError: 'Something went wrong. Please try again.'
      }
    },
    de: {
      team: 'LISC-2010',
      signIn: {
        welcomeTo: 'Willkommen bei',
        getStarted: 'Gib deine E-Mail-Adresse ein, um einen Anmeldecode zu erhalten',
        emailLabel: 'E-Mail-Adresse',
        emailNote: 'Bitte verwende die E-Mail-Adresse, mit der du dich im Verein angemeldet hast.',
        sendCode: 'Code senden',
        codeTitle: 'Schau in deine E-Mails!',
        codeNote: 'Wir haben einen 6-stelligen Code an {email} gesendet. Gib ihn hier ein. Falls du ihn nicht siehst, schau auch im Spam-Ordner nach.',
        codeLabel: '6-stelliger Code',
        verify: 'Bestätigen',
        backToEmail: 'Zurück',
        resendCode: 'Code erneut senden',
        codeSent: 'Ein neuer Code wurde an deine E-Mail geschickt.',
        emailRequired: 'Bitte gib deine E-Mail-Adresse ein',
        codeInvalid: 'Bitte gib den 6-stelligen Code aus deiner E-Mail ein',
      },
      profile: {
        header: 'Mein Profil',
        logout: 'Abmelden',
        edit: 'Profil bearbeiten',
        cancel: 'Abbrechen',
      },
      profileSetup: {
        header: 'Profil einrichten',
        firstName: 'Vorname',
        lastName: 'Nachname',
        nickName: 'Spitzname (optional)',
        picture: 'Profilbild (optional)',
        save: 'Speichern',
        firstNameRequired: 'Bitte gib deinen Vornamen ein',
        lastNameRequired: 'Bitte gib deinen Nachnamen ein',
      },
      events: {
        header: 'Events',
        back: 'Zurück',
        noEvents: 'Noch keine Events.',
        fab: 'Neues Event',
        development: {
          title: 'Entwicklungs-Vorschau',
          message: 'Diese App ist noch in der Entwicklung. Wenn dir Fehler auffallen, melde dich gerne bei Kathi.'
        },
        new: {
          header: 'Neues Event',
          title: 'Name',
          date: 'Datum',
          time: 'Uhrzeit',
          location: 'Ort',
          category: 'Kategorie',
          description: 'Beschreibung',
          create: 'Erstellen',
          required: 'Pflichtfeld',
          pastDate: 'Datum darf nicht in der Vergangenheit liegen',
          timeInvalid: 'Bitte gib eine Uhrzeit im Format HH:MM ein'
        },
        createdBy: 'von',
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
          title: 'Dein Account wird geprüft',
          message: 'Dein Account wird gerade geprüft. Ein Admin schaut sich das bald an.',
          checkStatus: 'Status prüfen',
          stillPending: 'Noch ausstehend — schau später nochmal vorbei.'
        },
        rejected: {
          title: 'Antrag abgelehnt',
          message: 'Dein Antrag wurde leider abgelehnt. Wenn du denkst, dass das ein Fehler ist, melde dich direkt beim Verein.'
        },
        cancel: 'Abbrechen'
      },
      select: {
        loading: 'Wird geladen…'
      },
      navigation: {
        events: 'Events',
        profile: 'Profil',
      },
      request: {
        errorTitle: 'Etwas ist schiefgelaufen',
        unexpectedError: 'Etwas hat nicht geklappt. Bitte versuch es nochmal.'
      }
    }
  }
})
