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
	Text    string   `json:"text"`
}

func (c *Client) SendOTP(ctx context.Context, to, code string) error {
	body := sendEmailRequest{
		From:    c.fromAddress,
		To:      []string{to},
		Subject: "Your Clubrizer sign-in code",
		Text:    fmt.Sprintf("Your Clubrizer sign-in code is: %s\n\nThis code expires in 1 hour.", code),
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
