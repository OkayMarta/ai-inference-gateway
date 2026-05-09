package clients

import (
	"net/http"
	"time"
)

type BillingClient struct {
	*ProxyClient
}

func NewBillingClient(baseURL, internalServiceToken string) *BillingClient {
	return &BillingClient{
		ProxyClient: newProxyClient(baseURL, "billing service unavailable", internalServiceToken, &http.Client{
			Timeout: 5 * time.Second,
		}),
	}
}
