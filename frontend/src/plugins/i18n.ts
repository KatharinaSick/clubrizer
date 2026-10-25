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
        back: 'Back',
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
          step3Title: 'Tell us who\'s joining',
          step3: 'Choose if the account is just for you, for you and your kids, or only your kids. If you have kids, add their first names.',
          step4Title: 'Wait for approval',
          step4: 'An admin checks that you\'re a member and unlocks your account. Once that\'s done, someone will let you know.',
          step5Title: 'Set up your profile',
          step5: 'Add your name, and a photo if you like. If you added kids, add their details too. Now you\'re all set.',
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
          androidNoteTitle: 'Android',
          androidNote: 'Make sure to use Chrome for this. If your phone shows a security warning while installing, it\'s safe to tap "Install anyway".',
          whyTitle: 'Why does my phone show a warning?',
          whyText: 'Installing just adds this same website to your home screen. To do that, your phone wraps it in a small app, and that wrapper lags a bit behind the newest Android version, so your phone flags it. It\'s a version check, not a sign anything is unsafe. Tap "Install anyway" to continue.',
          whySource1: 'How these warnings work (Google)',
        },
        footer: 'Still have questions? Just get in touch with your club and someone will help you.',
      },
      profile: {
        header: 'My Profile',
        logout: 'Logout',
        edit: 'Edit Profile',
        myKids: 'My Kids',
        manageMembers: 'Manage Members',
        menuLabel: 'Open menu',
        cancel: 'Cancel',
        howItWorks: 'How it works',
        changelog: "What's New",
        notifications: 'Notifications',
      },
      notifications: {
        header: 'Notifications',
        off: {
          title: 'Never miss a thing',
          intro: 'Get a heads-up on this device when something new happens in your club.',
          turnOn: 'Turn on',
        },
        on: {
          status: 'On for this device',
          sendTest: 'Send test',
          testSent: 'Sent! It should pop up in a few seconds.',
          testMissing: 'Nothing showing up? Check that notifications are allowed for your browser in your device settings and that "Do Not Disturb" is off.',
          testUnreachable: 'We couldn\'t reach this device. Please turn notifications off and on again, then try once more.',
          turnOffThisDevice: 'Turn off on this device',
        },
        homeScreen: {
          title: 'One quick step first',
          intro: 'On iPhone or iPad, notifications only work once the app is on your home screen:',
          step1: 'Tap the Share button at the bottom (the square with an arrow pointing up).',
          step2: 'Scroll down a little and choose "Add to Home Screen".',
          step3: 'Tap "Add" in the top right.',
          step4: 'Open the app from your home screen and come back here.',
        },
        unsupported: 'This browser can\'t show notifications. Try Chrome, Firefox, Edge or Safari.',
        blocked: {
          title: 'Notifications are blocked',
          message: 'To get them, allow notifications for this app in your phone or browser settings, then come back here.',
        },
        deviceError: 'That didn\'t work on this device. Please try again.',
        settings: {
          title: 'What you\'ll hear about',
          newEvents: {
            label: 'New events',
            description: 'When someone adds an event.',
          },
          responses: {
            label: 'Who\'s coming',
            description: 'When someone replies to an event you created.',
          },
          comments: {
            label: 'Comments',
            description: 'On events you created or are going to.',
          },
          eventChanges: {
            label: 'Cancelled events',
            description: 'When an event you\'re going to is cancelled, back on or removed.',
          },
          memberRequests: {
            label: 'Join requests',
            description: 'When someone wants to join the club. You\'ll still get an email.',
          },
        },
      },
      pushNotifications: {
        test: {
          title: 'It works!',
          body: 'Notifications are all set on this device.',
        },
        memberRequest: {
          title: 'New join request',
          body: 'Someone wants to join the club. Tap to take a look.',
        },
        accountApproved: {
          title: 'Welcome to the club!',
          body: 'You\'re in. Tap to get started.',
        },
        accountRejected: {
          title: 'Your request was declined',
          body: 'Your request to join the club wasn\'t approved.',
        },
        kidApproved: {
          title: '{name} is in!',
          body: 'You can now sign {name} up for events.',
        },
        kidRejected: {
          title: '{name} wasn\'t approved',
          body: '{name}\'s request to join the club wasn\'t approved.',
        },
        eventCancelled: {
          title: 'Event cancelled',
          body: '{title} ({date}) is cancelled.',
        },
        eventUncancelled: {
          title: 'Event is back on',
          body: '{title} ({date}) is happening after all.',
        },
        eventDeleted: {
          title: 'Event removed',
          body: '{title} ({date}) was removed.',
        },
        eventComment: {
          title: '{name} commented on {title}',
          body: '{comment}',
        },
        newEvent: {
          title: 'New event: {title}',
          body: '{date}',
          groupTitle: '{count} new events',
          groupBody: '{titles}',
        },
        eventResponse: {
          title: '{title}',
          body: '{answer}',
          groupTitle: '{count} new replies to {title}',
          groupBody: '{answers}',
          yes: '{name} is coming',
          no: '{name} can\'t make it',
        },
        fallback: {
          title: 'LISC-2010',
          body: 'Something new happened in your club.',
        },
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
      accountSetup: {
        header: 'Complete your details',
        intro: 'Before you get started, please add a few details so your club knows who\'s who. Adding a photo is optional, but it really helps everyone recognise each other.',
        yourData: 'Your data',
        firstName: 'First Name',
        lastName: 'Last Name',
        nickName: 'Nickname (optional)',
        changePhoto: 'Add a photo',
        kidLastName: 'Last Name',
        kidPhoto: 'Add a photo',
        save: 'Save',
        firstNameRequired: 'Please enter your first name',
        lastNameRequired: 'Please enter your last name',
        kidLastNameRequired: 'Please enter a last name',
      },
      events: {
        header: 'Events',
        back: 'Back',
        noEvents: 'No events yet.',
        fab: 'New Event',
        filter: {
          allEvents: 'All events',
        },
        past: {
          show: 'Show past events',
          title: 'Past events',
          empty: 'No past events.',
        },
        cancelled: 'Cancelled',
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
            cancel: 'Cancel event',
            restore: 'Restore event',
          },
          deleteConfirm: {
            title: 'Delete this event?',
            message: 'This permanently deletes the event for everyone. This cannot be undone.',
            confirm: 'Delete',
            cancel: 'Cancel',
          },
          cancelConfirm: {
            title: 'Cancel this event?',
            message: 'The event stays visible for everyone but will be marked as cancelled.',
            confirm: 'Yes, cancel it',
            back: 'Keep event',
          },
          restoreConfirm: {
            title: 'Restore this event?',
            message: 'The event will be active again and no longer shown as cancelled.',
            confirm: 'Yes, restore it',
            back: 'Keep as cancelled',
          },
          cancelledBanner: 'This event has been cancelled.',
          attendees: {
            going: '{count} going',
            notGoing: '{count} not going',
            goingLabel: 'Going',
            notGoingLabel: 'Not going',
            noResponses: 'No responses yet',
            kid: 'Kid',
            showAll: 'Show all',
            less: 'Show less',
          },
          comments: {
            title: 'Comments',
            placeholder: 'Write a comment…',
            post: 'Post',
            empty: 'No comments yet. Write the first one!',
          }
        }
      },
      select: {
        loading: 'Loading…'
      },
      userProfileModal: {
        kidOf: 'Kid of {parent}',
      },
      onboarding: {
        title: 'Who are you signing up?',
        hint: 'Choose who this account is for.',
        justMe: 'Just me',
        meAndKids: 'Me and my kids',
        onlyKids: 'Only my kids',
        kidsTitle: 'Add your kids',
        kidsHint: 'Enter the first name of each kid you\'re signing up.',
        kidNamePlaceholder: 'First name',
        addAnotherKid: 'Add another child (optional)',
        removeKid: 'Remove',
        continue: 'Continue',
        back: 'Back',
        needKid: 'Please add at least one kid.',
        cancel: 'Cancel',
      },
      kids: {
        firstName: 'First Name',
        lastName: 'Last Name',
        add: 'Add Kid',
        editTitle: 'Edit Kid',
        empty: 'No kids added yet.',
        firstNameRequired: 'Please enter a first name',
        lastNameRequired: 'Please enter a last name',
        pendingBadge: 'Pending',
        edit: 'Edit',
        remove: 'Remove',
        save: 'Save',
        cancel: 'Cancel',
        changePhoto: 'Change photo',
        removeConfirm: {
          title: 'Remove this kid?',
          message: 'This removes {name} and their event responses. This cannot be undone.',
          confirm: 'Remove',
          cancel: 'Cancel',
        },
      },
      manageKids: {
        header: 'My Kids',
      },
      manageMembers: {
        header: 'Manage Members',
        requests: {
          title: 'Approval requests',
          empty: 'No pending requests right now.',
          newMember: 'New member',
          guardian: 'Manages kids only',
          newKids: 'New kids to review',
          kids: 'Kids',
          approve: 'Approve',
          reject: 'Decline',
        },
        members: {
          title: 'Members',
          empty: 'No members yet.',
          guardian: 'Guardian',
          kids: 'Kids',
        },
      },
      notificationsAnnouncement: {
        title: 'New: notifications',
        message: 'Get a heads-up on your phone when there\'s a new event, someone replies or an event is cancelled.',
        settingsHint: 'You can change this anytime in your profile under Notifications.',
        turnOn: 'Turn on',
        later: 'Later',
        noThanks: 'No thanks',
        deviceError: 'That didn\'t work on this device. Please try again.',
        notAllowed: 'Notifications weren\'t allowed. Tap "Turn on" to try again.',
        homeScreen: {
          intro: 'On iPhone or iPad, add the app to your home screen first:',
          step1: 'Tap the Share button at the bottom (the square with an arrow pointing up).',
          step2: 'Scroll down a little and choose "Add to Home Screen".',
          step3: 'Tap "Add" in the top right.',
          step4: 'Open the app from your home screen. We\'ll ask you again there.',
        },
        done: {
          title: 'You\'re all set!',
          message: 'You can choose what you hear about in your profile under Notifications.',
          close: 'Got it',
        },
        blocked: {
          title: 'Notifications are blocked',
          message: 'No problem! If you change your mind, you can allow them in your phone or browser settings.',
          close: 'Got it',
        },
      },
      pendingApproval: {
        pending: {
          title: 'You\'re almost in!',
          message: 'An admin will check your request soon. Hang tight, it won\'t be long.',
          checkStatus: 'Check again',
          stillPending: 'Not yet. Check back a little later.',
          yourKids: 'Your kids: {names}'
        },
        notifications: {
          question: 'Want a heads-up as soon as you\'re in?',
          turnOn: 'Notify me',
          on: 'Great! We\'ll let you know on this device.',
          homeScreen: 'On iPhone or iPad, add the app to your home screen and open it from there first. Then you can turn this on.',
          deviceError: 'That didn\'t work on this device. Please try again.',
        },
        rejected: {
          title: 'Your request wasn\'t approved',
          message: 'If you think this is a mistake, just reach out to the club directly.'
        },
        signOut: 'Sign out'
      },
      changelog: {
        header: "What's New",
        entries: [
          {
            date: 'October 25, 2026',
            items: [
              'You can now get notifications on your phone, e.g. about new events, replies and comments. Turn them on in your profile under Notifications.',
              'You choose what you hear about.',
              'The screen you see while waiting for approval got a fresh look, and it can let you know as soon as you\'re in.',
            ],
          },
          {
            date: 'August 11, 2026',
            items: [
              'You can now filter events by category.',
              'Past events are shown separately at the bottom of the list.',
              'There is now a "What\'s New" page. You\'re looking at it.',
            ],
          },
          {
            date: 'August 6, 2026',
            items: [
              'Tap on any profile photo to see who it is.',
              'The attendee list on event details now takes up less space.',
              'Events stay visible until 4 hours after they start.',
              'The member count in the admin view now correctly counts kids and excludes guardians.',
            ],
          },
          {
            date: 'August 5, 2026',
            items: [
              'The year is now shown as a divider in the event list.',
              'Menus now respond faster on mobile.',
            ],
          },
          {
            date: 'August 4, 2026',
            items: [
              'You can now leave comments on events.',
            ],
          },
          {
            date: 'August 3, 2026',
            items: [
              'You can now manage your kids and respond to events on their behalf.',
              'Admins can approve or decline new members and kids directly in the app.',
            ],
          },
          {
            date: 'July 10, 2026',
            items: [
              'Admins can now cancel events and restore them later.',
              'Tap on a profile photo to see that person\'s name.',
              'People who are going are now shown first in the attendee list.',
            ],
          },
          {
            date: 'July 9, 2026',
            items: [
              'Admins can now delete events.',
            ],
          },
          {
            date: 'July 7, 2026',
            items: [
              'Fixed a bug where a valid sign-in code was sometimes shown as invalid.',
            ],
          },
          {
            date: 'June 21, 2026',
            items: [
              'The app now has a proper layout on larger screens.',
            ],
          },
          {
            date: 'June 15, 2026',
            items: [
              'You can now install the app on your phone\'s home screen.',
            ],
          },
          {
            date: 'June 6, 2026',
            items: [
              'Sign in now works with a 6-digit code sent to your email.',
              'The event organizer is now shown on the event detail page.',
            ],
          },
          {
            date: 'March 31, 2026',
            items: [
              'New accounts now need to be approved by an admin before they can access the app.',
            ],
          },
          {
            date: 'March 28, 2026',
            items: [
              'You can now say whether you are going to an event.',
            ],
          },
          {
            date: 'March 26, 2026',
            items: [
              'Tap on an event to see its full details.',
            ],
          },
          {
            date: 'March 6, 2026',
            items: [
              'Only upcoming events are shown in the list.',
              'Events in the past can no longer be created.',
            ],
          },
          {
            date: 'February 15, 2026',
            items: [
              'Events can now be created and browsed in a list.',
            ],
          },
          {
            date: 'February 26, 2025',
            items: [
              'The app is now available in German.',
            ],
          },
          {
            date: 'February 13, 2025',
            items: [
              'The app launched. Sign in with your Google account to get started.',
            ],
          },
        ],
      },
      navigation: {
        events: 'Events',
        profile: 'Profile',
        signIn: 'Sign in',
        changelog: "What's New",
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
        back: 'Zurück',
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
          step3Title: 'Sag uns, wer dabei ist',
          step3: 'Wähle, ob der Account nur für dich, für dich und deine Kinder oder nur für deine Kinder ist. Wenn du Kinder hast, trag ihre Vornamen ein.',
          step4Title: 'Warte auf die Freigabe',
          step4: 'Ein Admin prüft, ob du Mitglied bist, und schaltet dein Konto frei. Sobald das erledigt ist, meldet sich jemand bei dir.',
          step5Title: 'Richte dein Profil ein',
          step5: 'Trag deinen Namen ein und, wenn du magst, ein Foto. Wenn du Kinder hinzugefügt hast, ergänze auch ihre Angaben. Fertig, jetzt kann es losgehen.',
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
          androidNoteTitle: 'Android',
          androidNote: 'Verwende am besten Chrome. Wenn dein Handy beim Installieren eine Sicherheitswarnung zeigt, kannst du bedenkenlos auf „Trotzdem installieren“ tippen.',
          whyTitle: 'Warum zeigt mein Handy eine Warnung?',
          whyText: 'Beim Installieren wird nur dieselbe Website auf deinen Startbildschirm gelegt. Dafür verpackt dein Handy sie in eine kleine App, und diese Verpackung hinkt der neuesten Android-Version etwas hinterher, deshalb die Warnung. Es ist eine Versionsprüfung, kein Zeichen, dass etwas unsicher ist. Tippe auf „Trotzdem installieren“.',
          whySource1: 'Wie diese Warnungen funktionieren (Google)',
        },
        footer: 'Noch Fragen? Melde dich einfach bei deinem Verein, dann hilft dir jemand weiter.',
      },
      profile: {
        header: 'Mein Profil',
        logout: 'Abmelden',
        edit: 'Profil bearbeiten',
        myKids: 'Meine Kinder',
        manageMembers: 'Mitglieder verwalten',
        menuLabel: 'Menü öffnen',
        cancel: 'Abbrechen',
        howItWorks: 'So funktioniert die App',
        changelog: 'Neu in der App',
        notifications: 'Benachrichtigungen',
      },
      notifications: {
        header: 'Benachrichtigungen',
        off: {
          title: 'Nichts mehr verpassen',
          intro: 'Erfahr auf diesem Gerät sofort, wenn sich in deinem Verein etwas tut.',
          turnOn: 'Einschalten',
        },
        on: {
          status: 'Auf diesem Gerät an',
          sendTest: 'Testen',
          testSent: 'Gesendet! Sie sollte gleich auftauchen.',
          testMissing: 'Nichts angekommen? Schau, ob Benachrichtigungen für deinen Browser in den Einstellungen deines Geräts erlaubt sind und „Nicht stören“ aus ist.',
          testUnreachable: 'Wir konnten dieses Gerät nicht erreichen. Schalte die Benachrichtigungen aus und wieder ein und versuch es dann nochmal.',
          turnOffThisDevice: 'Auf diesem Gerät ausschalten',
        },
        homeScreen: {
          title: 'Noch ein kurzer Schritt',
          intro: 'Auf dem iPhone oder iPad klappen Benachrichtigungen erst, wenn die App auf deinem Home-Bildschirm ist:',
          step1: 'Tipp unten auf das Teilen-Symbol (das Quadrat mit dem Pfeil nach oben).',
          step2: 'Scroll ein Stück nach unten und wähle „Zum Home-Bildschirm“.',
          step3: 'Tipp oben rechts auf „Hinzufügen“.',
          step4: 'Öffne die App über deinen Home-Bildschirm und komm wieder hierher.',
        },
        unsupported: 'Dieser Browser kann leider keine Benachrichtigungen anzeigen. Probier es mit Chrome, Firefox, Edge oder Safari.',
        blocked: {
          title: 'Benachrichtigungen sind blockiert',
          message: 'Erlaub Benachrichtigungen für diese App in den Einstellungen deines Handys oder Browsers und komm dann wieder hierher.',
        },
        deviceError: 'Das hat auf diesem Gerät nicht geklappt. Bitte versuch es nochmal.',
        settings: {
          title: 'Worüber du Bescheid bekommst',
          newEvents: {
            label: 'Neue Events',
            description: 'Wenn jemand ein Event einträgt.',
          },
          responses: {
            label: 'Zu- und Absagen',
            description: 'Wenn jemand auf eines deiner Events antwortet.',
          },
          comments: {
            label: 'Kommentare',
            description: 'Bei Events, die du erstellt oder zugesagt hast.',
          },
          eventChanges: {
            label: 'Abgesagte Events',
            description: 'Wenn ein Event, zu dem du kommst, abgesagt, wieder angesetzt oder entfernt wird.',
          },
          memberRequests: {
            label: 'Beitrittsanfragen',
            description: 'Wenn jemand dem Verein beitreten möchte. Eine E-Mail bekommst du trotzdem.',
          },
        },
      },
      pushNotifications: {
        test: {
          title: 'Es funktioniert!',
          body: 'Benachrichtigungen sind auf diesem Gerät startklar.',
        },
        memberRequest: {
          title: 'Neue Beitrittsanfrage',
          body: 'Jemand möchte dem Verein beitreten. Tipp hier, um sie anzusehen.',
        },
        accountApproved: {
          title: 'Willkommen im Verein!',
          body: 'Du bist dabei. Tipp hier, um loszulegen.',
        },
        accountRejected: {
          title: 'Deine Anfrage wurde abgelehnt',
          body: 'Deine Anfrage, dem Verein beizutreten, wurde leider nicht angenommen.',
        },
        kidApproved: {
          title: '{name} ist dabei!',
          body: 'Du kannst {name} jetzt für Events anmelden.',
        },
        kidRejected: {
          title: '{name} wurde nicht angenommen',
          body: '{name} wurde leider nicht in den Verein aufgenommen.',
        },
        eventCancelled: {
          title: 'Event abgesagt',
          body: '{title} ({date}) ist abgesagt.',
        },
        eventUncancelled: {
          title: 'Event findet doch statt',
          body: '{title} ({date}) ist wieder angesetzt.',
        },
        eventDeleted: {
          title: 'Event entfernt',
          body: '{title} ({date}) wurde entfernt.',
        },
        eventComment: {
          title: '{name} hat {title} kommentiert',
          body: '{comment}',
        },
        newEvent: {
          title: 'Neues Event: {title}',
          body: '{date}',
          groupTitle: '{count} neue Events',
          groupBody: '{titles}',
        },
        eventResponse: {
          title: '{title}',
          body: '{answer}',
          groupTitle: '{count} neue Antworten zu {title}',
          groupBody: '{answers}',
          yes: '{name} kommt',
          no: '{name} kann nicht',
        },
        fallback: {
          title: 'LISC-2010',
          body: 'In deinem Verein gibt es etwas Neues.',
        },
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
      accountSetup: {
        header: 'Daten vervollständigen',
        intro: 'Bevor du loslegst, ergänze bitte ein paar Angaben, damit dein Verein weiß, wer wer ist. Ein Foto ist optional, hilft aber allen, sich gegenseitig zu erkennen.',
        yourData: 'Deine Daten',
        firstName: 'Vorname',
        lastName: 'Nachname',
        nickName: 'Spitzname (optional)',
        changePhoto: 'Foto hinzufügen',
        kidLastName: 'Nachname',
        kidPhoto: 'Foto hinzufügen',
        save: 'Speichern',
        firstNameRequired: 'Bitte gib deinen Vornamen ein',
        lastNameRequired: 'Bitte gib deinen Nachnamen ein',
        kidLastNameRequired: 'Bitte gib einen Nachnamen ein',
      },
      events: {
        header: 'Events',
        back: 'Zurück',
        noEvents: 'Noch keine Events.',
        fab: 'Neues Event',
        filter: {
          allEvents: 'Alle Events',
        },
        past: {
          show: 'Vergangene Events anzeigen',
          title: 'Vergangene Events',
          empty: 'Keine vergangenen Events.',
        },
        cancelled: 'Abgesagt',
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
            cancel: 'Event absagen',
            restore: 'Event wiederherstellen',
          },
          deleteConfirm: {
            title: 'Event wirklich löschen?',
            message: 'Das Event wird für alle dauerhaft gelöscht. Das kann nicht rückgängig gemacht werden.',
            confirm: 'Löschen',
            cancel: 'Abbrechen',
          },
          cancelConfirm: {
            title: 'Event absagen?',
            message: 'Das Event bleibt für alle sichtbar, wird aber als abgesagt markiert.',
            confirm: 'Ja, absagen',
            back: 'Event behalten',
          },
          restoreConfirm: {
            title: 'Event wiederherstellen?',
            message: 'Das Event ist dann wieder aktiv und wird nicht mehr als abgesagt angezeigt.',
            confirm: 'Ja, wiederherstellen',
            back: 'Abgesagt lassen',
          },
          cancelledBanner: 'Dieses Event wurde abgesagt.',
          attendees: {
            going: '{count} Zusagen',
            notGoing: '{count} Absagen',
            goingLabel: 'Zusagen',
            notGoingLabel: 'Absagen',
            noResponses: 'Noch keine Rückmeldungen',
            kid: 'Kind',
            showAll: 'Alle anzeigen',
            less: 'Weniger',
          },
          comments: {
            title: 'Kommentare',
            placeholder: 'Schreib einen Kommentar…',
            post: 'Senden',
            empty: 'Noch keine Kommentare. Schreib den ersten!',
          }
        }
      },
      notificationsAnnouncement: {
        title: 'Neu: Benachrichtigungen',
        message: 'Erfahr am Handy sofort, wenn es ein neues Event gibt, jemand antwortet oder ein Event abgesagt wird.',
        settingsHint: 'Du kannst das jederzeit in deinem Profil unter Benachrichtigungen ändern.',
        turnOn: 'Einschalten',
        later: 'Später',
        noThanks: 'Nein danke',
        deviceError: 'Das hat auf diesem Gerät nicht geklappt. Bitte versuch es nochmal.',
        notAllowed: 'Benachrichtigungen wurden nicht erlaubt. Tipp auf „Einschalten“, um es nochmal zu versuchen.',
        homeScreen: {
          intro: 'Auf dem iPhone oder iPad füg die App zuerst zu deinem Home-Bildschirm hinzu:',
          step1: 'Tipp unten auf das Teilen-Symbol (das Quadrat mit dem Pfeil nach oben).',
          step2: 'Scroll ein Stück nach unten und wähle „Zum Home-Bildschirm“.',
          step3: 'Tipp oben rechts auf „Hinzufügen“.',
          step4: 'Öffne die App über deinen Home-Bildschirm. Dort fragen wir dich nochmal.',
        },
        done: {
          title: 'Alles startklar!',
          message: 'In deinem Profil unter Benachrichtigungen kannst du aussuchen, worüber du Bescheid bekommst.',
          close: 'Alles klar',
        },
        blocked: {
          title: 'Benachrichtigungen sind blockiert',
          message: 'Kein Problem! Falls du es dir anders überlegst, kannst du sie in den Einstellungen deines Handys oder Browsers erlauben.',
          close: 'Alles klar',
        },
      },
      pendingApproval: {
        pending: {
          title: 'Fast geschafft!',
          message: 'Ein Admin schaut sich deine Anfrage bald an. Es dauert nicht mehr lange.',
          checkStatus: 'Nochmal schauen',
          stillPending: 'Noch nicht. Schau später nochmal vorbei.',
          yourKids: 'Deine Kinder: {names}'
        },
        notifications: {
          question: 'Willst du Bescheid bekommen, sobald du dabei bist?',
          turnOn: 'Bescheid geben',
          on: 'Super! Wir sagen dir auf diesem Gerät Bescheid.',
          homeScreen: 'Auf dem iPhone oder iPad füg die App zuerst zu deinem Home-Bildschirm hinzu und öffne sie von dort. Dann kannst du das hier einschalten.',
          deviceError: 'Das hat auf diesem Gerät nicht geklappt. Bitte versuch es nochmal.',
        },
        rejected: {
          title: 'Deine Anfrage wurde nicht angenommen',
          message: 'Wenn du denkst, dass das ein Fehler ist, melde dich einfach direkt beim Verein.'
        },
        signOut: 'Abmelden'
      },
      select: {
        loading: 'Wird geladen…'
      },
      userProfileModal: {
        kidOf: 'Kind von {parent}',
      },
      onboarding: {
        title: 'Wen meldest du an?',
        hint: 'Wähle, für wen dieser Account ist.',
        justMe: 'Nur mich',
        meAndKids: 'Mich und meine Kinder',
        onlyKids: 'Nur meine Kinder',
        kidsTitle: 'Kinder hinzufügen',
        kidsHint: 'Gib den Vornamen von jedem Kind ein, das du anmeldest.',
        kidNamePlaceholder: 'Vorname',
        addAnotherKid: 'Weiteres Kind (optional)',
        removeKid: 'Entfernen',
        continue: 'Weiter',
        back: 'Zurück',
        needKid: 'Bitte füge mindestens ein Kind hinzu.',
        cancel: 'Abbrechen',
      },
      kids: {
        firstName: 'Vorname',
        lastName: 'Nachname',
        add: 'Kind hinzufügen',
        editTitle: 'Kind bearbeiten',
        empty: 'Noch keine Kinder hinzugefügt.',
        firstNameRequired: 'Bitte gib einen Vornamen ein',
        lastNameRequired: 'Bitte gib einen Nachnamen ein',
        pendingBadge: 'In Prüfung',
        edit: 'Bearbeiten',
        remove: 'Entfernen',
        save: 'Speichern',
        cancel: 'Abbrechen',
        changePhoto: 'Foto ändern',
        removeConfirm: {
          title: 'Kind entfernen?',
          message: 'Damit werden {name} und die Event-Rückmeldungen entfernt. Das kann nicht rückgängig gemacht werden.',
          confirm: 'Entfernen',
          cancel: 'Abbrechen',
        },
      },
      manageKids: {
        header: 'Meine Kinder',
      },
      manageMembers: {
        header: 'Mitglieder verwalten',
        requests: {
          title: 'Anfragen',
          empty: 'Aktuell keine offenen Anfragen.',
          newMember: 'Neues Mitglied',
          guardian: 'Verwaltet nur Kinder',
          newKids: 'Neue Kinder zur Freigabe',
          kids: 'Kinder',
          approve: 'Freigeben',
          reject: 'Ablehnen',
        },
        members: {
          title: 'Mitglieder',
          empty: 'Noch keine Mitglieder.',
          guardian: 'Betreuer:in',
          kids: 'Kinder',
        },
      },
      changelog: {
        header: 'Neu in der App',
        entries: [
          {
            date: '25. Oktober 2026',
            items: [
              'Du kannst jetzt Benachrichtigungen aufs Handy bekommen, z. B. zu neuen Events, Antworten und Kommentaren. Schalte sie in deinem Profil unter Benachrichtigungen ein.',
              'Du suchst aus, worüber du Bescheid bekommst.',
              'Der Bildschirm, während du auf deine Freischaltung wartest, sieht jetzt frischer aus und kann dir Bescheid geben, sobald du dabei bist.',
            ],
          },
          {
            date: '11. August 2026',
            items: [
              'Du kannst Events jetzt nach Kategorie filtern.',
              'Vergangene Events findest du jetzt am Ende der Liste.',
              'Es gibt jetzt eine Seite „Neu in der App". Du schaust gerade drauf.',
            ],
          },
          {
            date: '6. August 2026',
            items: [
              'Tippe auf ein Profilbild, um zu sehen, wer dahinter steckt.',
              'Die Teilnehmerliste bei einem Event nimmt jetzt weniger Platz ein.',
              'Events bleiben noch 4 Stunden nach dem Start in der Liste sichtbar.',
              'Die Mitgliederzahl in der Verwaltung zählt Kinder jetzt korrekt und schließt Betreuer aus.',
            ],
          },
          {
            date: '5. August 2026',
            items: [
              'Das Jahr wird jetzt als Trennlinie in der Eventliste angezeigt.',
              'Menüs reagieren auf dem Handy jetzt schneller.',
            ],
          },
          {
            date: '4. August 2026',
            items: [
              'Du kannst jetzt Kommentare bei Events hinterlassen.',
            ],
          },
          {
            date: '3. August 2026',
            items: [
              'Du kannst deine Kinder jetzt verwalten und Events für sie beantworten.',
              'Admins können neue Mitglieder und Kinder direkt in der App freigeben oder ablehnen.',
            ],
          },
          {
            date: '10. Juli 2026',
            items: [
              'Admins können Events jetzt absagen und später wiederherstellen.',
              'Tippe auf ein Profilbild, um den Namen der Person zu sehen.',
              'Zusagen werden in der Teilnehmerliste zuerst angezeigt.',
            ],
          },
          {
            date: '9. Juli 2026',
            items: [
              'Admins können Events jetzt löschen.',
            ],
          },
          {
            date: '7. Juli 2026',
            items: [
              'Fehler behoben, bei dem ein gültiger Anmeldecode manchmal als ungültig angezeigt wurde.',
            ],
          },
          {
            date: '21. Juni 2026',
            items: [
              'Die App hat jetzt ein optimiertes Layout auf größeren Bildschirmen.',
            ],
          },
          {
            date: '15. Juni 2026',
            items: [
              'Du kannst die App jetzt auf deinem Startbildschirm installieren.',
            ],
          },
          {
            date: '6. Juni 2026',
            items: [
              'Die Anmeldung funktioniert jetzt mit einem 6-stelligen Code per E-Mail.',
              'Der Ersteller eines Events wird jetzt in den Event-Details angezeigt.',
            ],
          },
          {
            date: '31. März 2026',
            items: [
              'Neue Accounts müssen jetzt erst von einem Admin freigegeben werden, bevor sie die App nutzen können.',
            ],
          },
          {
            date: '28. März 2026',
            items: [
              'Du kannst jetzt angeben, ob du bei einem Event dabei bist.',
            ],
          },
          {
            date: '26. März 2026',
            items: [
              'Tippe auf ein Event, um alle Details zu sehen.',
            ],
          },
          {
            date: '6. März 2026',
            items: [
              'In der Eventliste werden jetzt nur bevorstehende Events angezeigt.',
              'Events in der Vergangenheit können nicht mehr erstellt werden.',
            ],
          },
          {
            date: '15. Februar 2026',
            items: [
              'Events können jetzt erstellt und in einer Liste angesehen werden.',
            ],
          },
          {
            date: '26. Februar 2025',
            items: [
              'Die App ist jetzt auch auf Deutsch verfügbar.',
            ],
          },
          {
            date: '13. Februar 2025',
            items: [
              'Die App ist gestartet. Melde dich mit deinem Google-Konto an, um loszulegen.',
            ],
          },
        ],
      },
      navigation: {
        events: 'Events',
        profile: 'Profil',
        signIn: 'Anmelden',
        changelog: 'Neu in der App',
      },
      request: {
        errorTitle: 'Etwas ist schiefgelaufen',
        unexpectedError: 'Etwas hat nicht geklappt. Bitte versuch es nochmal.'
      }
    }
  }
})
