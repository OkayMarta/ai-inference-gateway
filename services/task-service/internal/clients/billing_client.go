package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BillingClient struct {
	baseURL string
	client  *http.Client
}

type UserDTO struct {
	ID           string  `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	TokenBalance float64 `json:"tokenBalance"`
	Role         string  `json:"role"`
}

type DownstreamError struct {
	StatusCode int
	Message    string
}

func (e *DownstreamError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("billing service returned status %d", e.StatusCode)
}

func NewBillingClient(baseURL string) *BillingClient {
	return &BillingClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *BillingClient) Charge(userID, taskID string, amount float64) error {
	return c.postBillingEvent("/internal/billing/charge", userID, taskID, amount)
}

func (c *BillingClient) Refund(userID, taskID string, amount float64) error {
	return c.postBillingEvent("/internal/billing/refund", userID, taskID, amount)
}

func (c *BillingClient) GetUser(userID string) (*UserDTO, error) {
	resp, err := c.client.Get(c.baseURL + "/internal/users/" + url.PathEscape(userID))
	if err != nil {
		return nil, fmt.Errorf("billing service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeDownstreamError(resp)
	}

	var user UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode billing user response: %w", err)
	}

	return &user, nil
}

func (c *BillingClient) postBillingEvent(path, userID, taskID string, amount float64) error {
	body, err := json.Marshal(map[string]any{
		"userId": userID,
		"taskId": taskID,
		"amount": amount,
	})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.baseURL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("billing service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return decodeDownstreamError(resp)
	}

	return nil
}

func decodeDownstreamError(resp *http.Response) error {
	var body struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	return &DownstreamError{
		StatusCode: resp.StatusCode,
		Message:    body.Message,
	}
}
