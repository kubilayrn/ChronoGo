package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type WebhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string    `json:"message"`
	MessageID uuid.UUID `json:"messageId"`
}

type WebhookSender struct {
	client  *http.Client
	url     string
	authKey string
}

func NewWebhookSender() *WebhookSender {
	url := getEnv("WEBHOOK_URL", "https://webhook.site/c3f13233-1ed4-429e-9649-8133b3b9c9cd")
	authKey := getEnv("WEBHOOK_AUTH_KEY", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")

	return &WebhookSender{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		url:     url,
		authKey: authKey,
	}
}

func (s *WebhookSender) SendMessage(to, content string) (*uuid.UUID, error) {
	payload := WebhookRequest{
		To:      to,
		Content: content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", s.authKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &webhookResp.MessageID, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
