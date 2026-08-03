package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const resendAPIURL = "https://api.resend.com/emails"

type Client struct {
	apiKey      string
	fromAddress string
	httpClient  *http.Client
}

func NewClient(apiKey, fromAddress string) *Client {
	return &Client{
		apiKey:      apiKey,
		fromAddress: fromAddress,
		httpClient:  &http.Client{},
	}
}

type sendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Bcc     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

// send delivers one email to one or more visible recipients via Resend.
func (c *Client) send(ctx context.Context, to []string, subject, html string) error {
	return c.doSend(ctx, sendEmailRequest{
		From:    c.fromAddress,
		To:      to,
		Subject: subject,
		Html:    html,
	})
}

// sendBCC delivers one email whose recipients are all BCC'd, so none of them sees the
// others' addresses. Resend still requires a `to`, so the sender address is used there.
func (c *Client) sendBCC(ctx context.Context, bcc []string, subject, html string) error {
	return c.doSend(ctx, sendEmailRequest{
		From:    c.fromAddress,
		To:      []string{c.fromAddress},
		Bcc:     bcc,
		Subject: subject,
		Html:    html,
	})
}

// doSend performs the Resend HTTP call. All the public helpers build their German copy and
// delegate here so the HTTP plumbing lives in one place.
func (c *Client) doSend(ctx context.Context, body sendEmailRequest) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API returned status %d", resp.StatusCode)
	}

	return nil
}

// All copy below is hardcoded German until proper i18n is implemented:
// https://github.com/KatharinaSick/clubrizer/issues/15

func (c *Client) SendOTP(ctx context.Context, to, code string) error {
	return c.send(ctx, []string{to}, "Dein LISC-2010 Anmelde-Code",
		fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Dein Anmelde-Code</h2>
  <p style="color:#555;margin-bottom:24px">Gib diesen Code in LISC-2010 ein, um dich anzumelden:</p>
  <p style="font-size:36px;font-weight:bold;letter-spacing:12px;text-align:center;padding:20px;background:#f5f5f5;border-radius:8px;margin:0 0 24px">%s</p>
  <p style="color:#888;font-size:13px">Dieser Code ist 1 Stunde lang gültig. Wenn du dich nicht angemeldet hast, kannst du diese E-Mail ignorieren.</p>
</div>`, code))
}

// SendAccountApproved tells an account holder their request was approved and they can now
// sign in.
func (c *Client) SendAccountApproved(ctx context.Context, to string) error {
	return c.send(ctx, []string{to}, "Dein LISC-2010 Konto wurde freigeschaltet",
		`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Willkommen bei LISC-2010!</h2>
  <p style="color:#555;margin-bottom:24px">Deine Anfrage wurde freigegeben. Du kannst dich jetzt in der LISC-2010 App anmelden und loslegen.</p>
</div>`)
}

// SendAccountRejected tells an account holder their request was declined. Mirrors the
// wording of the rejected-account screen: a friendly dead end that points them at the club.
func (c *Client) SendAccountRejected(ctx context.Context, to string) error {
	return c.send(ctx, []string{to}, "Dein LISC-2010 Konto-Antrag",
		`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Deine Anfrage wurde abgelehnt</h2>
  <p style="color:#555;margin-bottom:24px">Leider wurde deine Anfrage abgelehnt. Wenn du denkst, dass das ein Fehler war, wende dich bitte direkt an deinen Verein.</p>
</div>`)
}

// SendKidApproved tells a parent that a kid they added after their account was already
// approved has now been approved and can take part.
func (c *Client) SendKidApproved(ctx context.Context, to, kidName string) error {
	if kidName == "" {
		kidName = "Dein Kind"
	}
	return c.send(ctx, []string{to}, "Freigabe für dein Kind bei LISC-2010",
		fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">%s wurde freigeschaltet</h2>
  <p style="color:#555;margin-bottom:24px">%s wurde freigegeben und kann jetzt bei Terminen mitmachen.</p>
</div>`, kidName, kidName))
}

// SendKidRejected tells a parent that a kid they added after their account was already
// approved has been declined.
func (c *Client) SendKidRejected(ctx context.Context, to, kidName string) error {
	if kidName == "" {
		kidName = "dein Kind"
	}
	return c.send(ctx, []string{to}, "Anfrage für dein Kind bei LISC-2010",
		fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Anfrage abgelehnt</h2>
  <p style="color:#555;margin-bottom:24px">Die Anfrage für %s wurde leider abgelehnt. Wenn du denkst, dass das ein Fehler war, wende dich bitte direkt an deinen Verein.</p>
</div>`, kidName))
}

// SendNewApprovalRequest notifies admins that a new account is waiting in the approval
// queue. All admins are BCC'd in one email, so no admin sees another's address.
//
// TODO(#50): link straight to the Manage Members screen once the backend has a configurable
// frontend base URL (part of the transactional-email overhaul).
func (c *Client) SendNewApprovalRequest(ctx context.Context, to []string, applicantName string) error {
	return c.sendBCC(ctx, to, "Neue Beitrittsanfrage bei LISC-2010",
		fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Neue Beitrittsanfrage</h2>
  <p style="color:#555;margin-bottom:24px">%s möchte LISC-2010 beitreten und wartet auf deine Freigabe. Öffne die LISC-2010 App, um die Anfrage zu prüfen.</p>
</div>`, applicantName))
}
