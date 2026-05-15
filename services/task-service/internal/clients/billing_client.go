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
	baseURL              string
	internalServiceToken string
	client               *http.Client
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

func NewBillingClient(baseURL, internalServiceToken string) *BillingClient {
	return &BillingClient{
		baseURL:              strings.TrimRight(baseURL, "/"),
		internalServiceToken: internalServiceToken,
		client:               &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *BillingClient) Charge(userID, taskID string, amount float64) error {
	return c.postBillingEvent("/internal/billing/charge", userID, taskID, amount)
}

func (c *BillingClient) Refund(userID, taskID string, amount float64) error {
	return c.postBillingEvent("/internal/billing/refund", userID, taskID, amount)
}

func (c *BillingClient) GetUser(userID string) (*UserDTO, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/internal/users/"+url.PathEscape(userID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Internal-Service-Token", c.internalServiceToken)

	resp, err := c.client.Do(req)
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

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Service-Token", c.internalServiceToken)

	resp, err := c.client.Do(req)
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
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("decode billing error response: %w", err)
	}
	if body.Message == "" {
		return fmt.Errorf("decode billing error response: missing message")
	}

	return &DownstreamError{
		StatusCode: resp.StatusCode,
		Message:    body.Message,
	}
}
