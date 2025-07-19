package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)


// TriggerN8NWebhook sends the user data to an n8n webhook endpoint
func TriggerN8NWebhook(users []User) error {
	webhookURL := os.Getenv("N8N_WEBHOOK_URL")
	if webhookURL == "" {
		log.Println("Skipping webhook – N8N_WEBHOOK_URL is not set")
		return nil
	}

	body, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Webhook sent successfully. Status: %s. Response: %s", resp.Status, string(respBody))
	return nil
}

