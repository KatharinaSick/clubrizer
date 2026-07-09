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
        howItWorks: 'New here? See how it works',
      },
      gettingStarted: {
        hero: {
          welcomeTo: 'Welcome to',
          intro: 'The app for your club. Stay connected, wherever you are.',
        },
        join: {
          title: 'How to get started',
          step1Title: 'Sign in with your email address',
          step1: 'Head to the {link} and use the same email address you gave your club when you became a member. That way we know it\'s really you.',
          step1Link: 'sign-in page',
          step2Title: 'Enter your code',
          step2: 'We send a 6-digit code to that email. Type it in and you\'re signed in.',
          step2Note: '💡 If you don\'t see the email within a few minutes, check your spam or junk folder.',
          step3Title: 'Wait for approval',
          step3: 'An admin checks that you\'re a member and unlocks your account. Once that\'s done, someone will let you know.',
          step4Title: 'Set up your profile',
          step4: 'Add your name and, if you like, a photo. Now you\'re all set.',
          cta: 'Go to sign in',
          ctaHome: 'Go to home',
        },
        install: {
          title: 'Install the app on your phone',
          intro: 'Add the app to your home screen and it opens just like any other app. You don\'t need an app store.',
          iphoneTitle: 'iPhone',
          iphoneStep1: 'Open this page on your phone in Safari.',
          iphoneStep2: 'Tap the Share button at the bottom (the square with an arrow pointing up).',
          iphoneStep3: 'Scroll down a little and choose "Add to Home Screen".',
          iphoneStep4: 'Tap "Add" in the top right.',
          androidTitle: 'Android',
          androidStep1: 'Open this page on your phone in Chrome.',
          androidStep2: 'Tap the menu in the top right (three dots).',
          androidStep3: 'Choose "Install app" or "Add to Home screen".',
          androidStep4: 'Confirm in the popup that appears by tapping "Install" or "Add".',
        },
        footer: 'Still have questions? Just get in touch with your club and someone will help you.',
      },
      profile: {
        header: 'My Profile',
        logout: 'Logout',
        edit: 'Edit Profile',
        cancel: 'Cancel',
        howItWorks: 'How it works',
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
          menu: {
            label: 'More options',
            delete: 'Delete event',
          },
          deleteConfirm: {
            title: 'Delete this event?',
            message: 'This permanently deletes the event for everyone. This cannot be undone.',
            confirm: 'Delete',
            cancel: 'Cancel',
          },
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
        signIn: 'Sign in',
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
        howItWorks: 'Neu hier? So geht\'s',
      },
      gettingStarted: {
        hero: {
          welcomeTo: 'Willkommen bei',
          intro: 'Die App für deinen Verein. Bleib in Verbindung, egal wo du bist.',
        },
        join: {
          title: 'So meldest du dich an',
          step1Title: 'Melde dich mit deiner E-Mail-Adresse an',
          step1: 'Geh zur {link} und nutze die E-Mail-Adresse, mit der du dich im Verein angemeldet hast. So wissen wir, dass du es wirklich bist.',
          step1Link: 'Anmeldung',
          step2Title: 'Gib deinen Code ein',
          step2: 'Wir schicken dir einen 6-stelligen Code an diese E-Mail. Tipp ihn ein und du bist angemeldet.',
          step2Note: '💡 Wenn die E-Mail nicht innerhalb weniger Minuten da ist, schau in deinem Spam-Ordner nach.',
          step3Title: 'Warte auf die Freigabe',
          step3: 'Ein Admin prüft, ob du Mitglied bist, und schaltet dein Konto frei. Sobald das erledigt ist, meldet sich jemand bei dir.',
          step4Title: 'Richte dein Profil ein',
          step4: 'Trag deinen Namen ein und, wenn du magst, ein Foto. Fertig, jetzt kann es losgehen.',
          cta: 'Zur Anmeldung',
          ctaHome: 'Zur Startseite',
        },
        install: {
          title: 'App am Handy installieren',
          intro: 'Leg die App auf deinen Startbildschirm. Dann öffnest du sie wie jede andere App, ganz ohne App Store.',
          iphoneTitle: 'iPhone',
          iphoneStep1: 'Öffne diese Seite am Handy in Safari.',
          iphoneStep2: 'Tipp unten auf das Teilen-Symbol (das Quadrat mit dem Pfeil nach oben).',
          iphoneStep3: 'Scroll ein Stück nach unten und wähle „Zum Home-Bildschirm“.',
          iphoneStep4: 'Tipp oben rechts auf „Hinzufügen“.',
          androidTitle: 'Android',
          androidStep1: 'Öffne diese Seite am Handy in Chrome.',
          androidStep2: 'Tipp oben rechts auf das Menü (drei Punkte).',
          androidStep3: 'Wähle „App installieren“ oder „Zum Startbildschirm hinzufügen“.',
          androidStep4: 'Bestätige im Fenster, das aufgeht, mit „Installieren“ oder „Hinzufügen“.',
        },
        footer: 'Noch Fragen? Melde dich einfach bei deinem Verein, dann hilft dir jemand weiter.',
      },
      profile: {
        header: 'Mein Profil',
        logout: 'Abmelden',
        edit: 'Profil bearbeiten',
        cancel: 'Abbrechen',
        howItWorks: 'So funktioniert die App',
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
          menu: {
            label: 'Mehr Optionen',
            delete: 'Event löschen',
          },
          deleteConfirm: {
            title: 'Event wirklich löschen?',
            message: 'Das Event wird für alle dauerhaft gelöscht. Das kann nicht rückgängig gemacht werden.',
            confirm: 'Löschen',
            cancel: 'Abbrechen',
          },
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
        signIn: 'Anmelden',
      },
      request: {
        errorTitle: 'Etwas ist schiefgelaufen',
        unexpectedError: 'Etwas hat nicht geklappt. Bitte versuch es nochmal.'
      }
    }
  }
})
