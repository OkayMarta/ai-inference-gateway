package clients

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AuthHeaders struct {
	UserID string
	Email  string
	Role   string
}

type ProxyClient struct {
	baseURL              string
	unavailableMessage   string
	internalServiceToken string
	client               *http.Client
}

func newProxyClient(baseURL, unavailableMessage, internalServiceToken string, client *http.Client) *ProxyClient {
	return &ProxyClient{
		baseURL:              strings.TrimRight(baseURL, "/"),
		unavailableMessage:   unavailableMessage,
		internalServiceToken: internalServiceToken,
		client:               client,
	}
}

func (c *ProxyClient) UnavailableMessage() string {
	return c.unavailableMessage
}

func (c *ProxyClient) Do(r *http.Request, downstreamPath string, auth *AuthHeaders) (*http.Response, error) {
	target, err := url.Parse(c.baseURL + downstreamPath)
	if err != nil {
		return nil, err
	}
	target.RawQuery = r.URL.RawQuery

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), r.Body)
	if err != nil {
		return nil, err
	}

	copyForwardHeaders(req.Header, r.Header)
	if auth != nil {
		req.Header.Set("X-User-ID", auth.UserID)
		req.Header.Set("X-User-Email", auth.Email)
		req.Header.Set("X-User-Role", auth.Role)
		req.Header.Set("X-Internal-Service-Token", c.internalServiceToken)
	}

	return c.client.Do(req)
}

func CopyResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func copyForwardHeaders(dst, src http.Header) {
	for _, key := range []string{"Accept", "Content-Type", "Authorization"} {
		if value := src.Get(key); value != "" {
			dst.Set(key, value)
		}
	}
}
