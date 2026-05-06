package clients

import (
	"net/http"
	"time"
)

type BillingClient struct {
	*ProxyClient
}

func NewBillingClient(baseURL string) *BillingClient {
	return &BillingClient{
		ProxyClient: newProxyClient(baseURL, "billing service unavailable", &http.Client{
			Timeout: 5 * time.Second,
		}),
	}
}
