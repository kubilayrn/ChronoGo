package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
	_ = godotenv.Load()

	url := os.Getenv("WEBHOOK_URL")
	authKey := os.Getenv("WEBHOOK_AUTH_KEY")

	if url == "" {
		log.Fatal("WEBHOOK_URL environment variable is required")
	}
	if authKey == "" {
		log.Fatal("WEBHOOK_AUTH_KEY environment variable is required")
	}

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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var webhookResp WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		responseStr := string(body)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			// Check if it's HTML or plain text (not JSON)
			if len(responseStr) > 0 && (responseStr[0] != '{' && responseStr[0] != '[') {
				mockID := uuid.New()
				log.Printf("Warning: Webhook returned non-JSON response (HTML/text), using mock messageId: %s", mockID.String())
				return &mockID, nil
			}
		}
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, responseStr)
	}

	// check response messageId
	if webhookResp.MessageID == uuid.Nil {
		mockID := uuid.New()
		log.Printf("Warning: Webhook returned empty messageId, using mock: %s", mockID.String())
		return &mockID, nil
	}

	return &webhookResp.MessageID, nil
}
