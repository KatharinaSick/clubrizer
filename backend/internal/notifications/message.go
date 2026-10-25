package notifications

// message is the payload delivered to the service worker. It only describes what happened —
// the service worker turns Type into a translated text (pushNotifications.<type> in the
// frontend's i18n), filling in Params, so all user-facing wording stays in the frontend.
type message struct {
	Type string `json:"type"`
	// Params are inserted into the translated text. A "date" param is an RFC 3339 timestamp
	// that the service worker formats in the device's language and timezone.
	Params map[string]string `json:"params,omitempty"`
	// URL is the in-app path opened when the notification is tapped.
	URL string `json:"url,omitempty"`
}

const (
	typeTest             = "test"
	typeMemberRequest    = "memberRequest"
	typeAccountApproved  = "accountApproved"
	typeAccountRejected  = "accountRejected"
	typeKidApproved      = "kidApproved"
	typeKidRejected      = "kidRejected"
	typeEventCancelled   = "eventCancelled"
	typeEventUncancelled = "eventUncancelled"
	typeEventDeleted     = "eventDeleted"
	typeEventComment     = "eventComment"
	typeNewEvent         = "newEvent"
	typeEventResponse    = "eventResponse"
)

// Categories group message types into the switches on the notification settings screen. A
// user who switched a category off gets no push notifications of its types; emails are not
// affected. Types without a category (test, approval decisions) are always sent.
const (
	categoryNewEvents      = "newEvents"
	categoryResponses      = "responses"
	categoryComments       = "comments"
	categoryEventChanges   = "eventChanges"
	categoryMemberRequests = "memberRequests"
)

var typeCategories = map[string]string{
	typeNewEvent:         categoryNewEvents,
	typeEventResponse:    categoryResponses,
	typeEventComment:     categoryComments,
	typeEventCancelled:   categoryEventChanges,
	typeEventUncancelled: categoryEventChanges,
	typeEventDeleted:     categoryEventChanges,
	typeMemberRequest:    categoryMemberRequests,
}
