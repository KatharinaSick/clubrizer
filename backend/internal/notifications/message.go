package notifications

import (
	"bytes"
	"encoding/json"
)

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

// maxParamLength caps each param (in characters). A push payload may be at most ~4KB, and
// a notification only shows a few lines anyway, so a long comment or event title is cut off
// rather than making the push fail for everyone.
const maxParamLength = 200

// encode turns the message into the push payload, shortening long params.
func encode(msg message) ([]byte, error) {
	params := make(map[string]string, len(msg.Params))
	for k, v := range msg.Params {
		params[k] = shorten(v, maxParamLength)
	}
	msg.Params = params

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	// Keep <, > and & as they are instead of \u003c etc., which take six bytes each.
	enc.SetEscapeHTML(false)
	if err := enc.Encode(msg); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

func shorten(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}

const (
	typeTest                     = "test"
	typeMemberRequest            = "memberRequest"
	typeAccountApproved          = "accountApproved"
	typeAccountRejected          = "accountRejected"
	typeKidApproved              = "kidApproved"
	typeKidRejected              = "kidRejected"
	typeEventCancelled           = "eventCancelled"
	typeEventUncancelled         = "eventUncancelled"
	typeEventTimeChanged         = "eventTimeChanged"
	typeEventPlaceChanged        = "eventPlaceChanged"
	typeEventTimeAndPlaceChanged = "eventTimeAndPlaceChanged"
	typeEventComment             = "eventComment"
	typeNewEvent                 = "newEvent"
	typeEventResponse            = "eventResponse"
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
	typeNewEvent:                 categoryNewEvents,
	typeEventResponse:            categoryResponses,
	typeEventComment:             categoryComments,
	typeEventCancelled:           categoryEventChanges,
	typeEventUncancelled:         categoryEventChanges,
	typeEventTimeChanged:         categoryEventChanges,
	typeEventPlaceChanged:        categoryEventChanges,
	typeEventTimeAndPlaceChanged: categoryEventChanges,
	typeMemberRequest:            categoryMemberRequests,
}
