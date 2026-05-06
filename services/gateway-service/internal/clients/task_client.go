package clients

import (
	"net/http"
	"time"
)

type TaskClient struct {
	*ProxyClient
}

func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		ProxyClient: newProxyClient(baseURL, "task service unavailable", &http.Client{
			Timeout: 5 * time.Second,
		}),
	}
}
