package handlers

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"gateway-service/internal/clients"
	"gateway-service/internal/middleware"
)

type ProxyHandler struct {
	billing *clients.BillingClient
	task    *clients.TaskClient
}

func NewProxyHandler(billing *clients.BillingClient, task *clients.TaskClient) *ProxyHandler {
	return &ProxyHandler{billing: billing, task: task}
}

func (h *ProxyHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, h.billing.ProxyClient, "/api/auth/register", nil)
}

func (h *ProxyHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, h.billing.ProxyClient, "/api/auth/login", nil)
}

func (h *ProxyHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, h.billing.ProxyClient, "/api/auth/forgot-password", nil)
}

func (h *ProxyHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, h.billing.ProxyClient, "/api/auth/reset-password", nil)
}

func (h *ProxyHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		RespondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.forward(w, r, h.billing.ProxyClient, "/internal/users/"+url.PathEscape(user.ID), authHeaders(user))
}

func (h *ProxyHandler) Models(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, h.task.ProxyClient, "/api/models", nil)
}

func (h *ProxyHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		RespondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.forward(w, r, h.task.ProxyClient, "/api/tasks", authHeaders(user))
}

func (h *ProxyHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUser(r)
	if !ok {
		RespondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	h.forward(w, r, h.task.ProxyClient, "/api/tasks/"+chi.URLParam(r, "id"), authHeaders(user))
}

func (h *ProxyHandler) forward(w http.ResponseWriter, r *http.Request, client *clients.ProxyClient, path string, auth *clients.AuthHeaders) {
	resp, err := client.Do(r, path, auth)
	if err != nil {
		RespondError(w, r, http.StatusServiceUnavailable, client.UnavailableMessage())
		return
	}

	clients.CopyResponse(w, resp)
}

func authHeaders(user *middleware.AuthUser) *clients.AuthHeaders {
	return &clients.AuthHeaders{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
}
