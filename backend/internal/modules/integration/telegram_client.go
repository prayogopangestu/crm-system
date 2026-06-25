package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TelegramClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewTelegramClient() *TelegramClient {
	return &TelegramClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.telegram.org",
	}
}

func NewTelegramClientWithBaseURL(baseURL string, client *http.Client) *TelegramClient {
	return &TelegramClient{httpClient: client, baseURL: strings.TrimRight(baseURL, "/")}
}

func (c *TelegramClient) Send(ctx context.Context, token, chatID, message string) error {
	if token == "" || chatID == "" {
		return errors.New("telegram token and chat id are required")
	}
	body, err := json.Marshal(map[string]string{"chat_id": chatID, "text": message})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/bot%s/sendMessage", c.baseURL, token)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	result, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer result.Body.Close()
	if result.StatusCode < 200 || result.StatusCode >= 300 {
		return fmt.Errorf("telegram returned status %d", result.StatusCode)
	}
	return nil
}
