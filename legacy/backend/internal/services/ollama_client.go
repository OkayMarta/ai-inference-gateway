package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaClient — HTTP-клієнт для спілкування з локальною нейромережею Ollama
type OllamaClient struct {
	baseURL string
	client  *http.Client
}

// NewOllamaClient створює клієнта. Cтавимо великий таймаут (5 хвилин), бо генерація тексту важкими моделями може займати багато часу
func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

// Структури для парсингу JSON-відповідей від Ollama
type OllamaModelInfo struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}

type ollamaTagsResponse struct {
	Models []OllamaModelInfo `json:"models"`
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Think  bool   `json:"think"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
}

// ListModels робить запит GET /api/tags до Ollama, щоб дізнатись, які моделі завантажені на ПК
func (c *OllamaClient) ListModels() ([]OllamaModelInfo, error) {
	quick := &http.Client{Timeout: 5 * time.Second}
	resp, err := quick.Get(c.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Ollama at %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	var result ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}
	return result.Models, nil
}

// Generate відправляє запит (POST /api/generate) в Ollama і чекає на згенерований текст
func (c *OllamaClient) Generate(model, prompt string) (string, error) {
	body, _ := json.Marshal(ollamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // Ми не використовуємо стрімінг, чекаємо всю відповідь цілком
		Think:  false,
	})

	resp, err := c.client.Post(c.baseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("Ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Ollama response: %w", err)
	}

	var result ollamaGenerateResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	response := strings.TrimSpace(result.Response)
	if response == "" {
		return "[Model returned empty response]", nil
	}
	return response, nil
}