package clients

import (
	"net/http"
	"time"
)

type TaskClient struct {
	*ProxyClient
}

func NewTaskClient(baseURL, internalServiceToken string) *TaskClient {
	return &TaskClient{
		ProxyClient: newProxyClient(baseURL, "task service unavailable", internalServiceToken, &http.Client{
			Timeout: 5 * time.Second,
		}),
	}
}
