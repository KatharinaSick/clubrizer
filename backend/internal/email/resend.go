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
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

func (c *Client) SendOTP(ctx context.Context, to, code string) error {
	// Hardcoded German until proper i18n is implemented: https://github.com/KatharinaSick/clubrizer/issues/15
	body := sendEmailRequest{
		From:    c.fromAddress,
		To:      []string{to},
		Subject: "Dein Clubrizer Anmelde-Code",
		Html: fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px;margin:0 auto">
  <h2 style="margin-bottom:8px">Dein Anmelde-Code</h2>
  <p style="color:#555;margin-bottom:24px">Gib diesen Code in Clubrizer ein, um dich anzumelden:</p>
  <p style="font-size:36px;font-weight:bold;letter-spacing:12px;text-align:center;padding:20px;background:#f5f5f5;border-radius:8px;margin:0 0 24px">%s</p>
  <p style="color:#888;font-size:13px">Dieser Code ist 1 Stunde lang gültig. Wenn du dich nicht angemeldet hast, kannst du diese E-Mail ignorieren.</p>
</div>`, code),
	}

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
