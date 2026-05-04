package handlers

import (
	"encoding/json"
	"net/http"

	"billing-service/internal/models"
	"billing-service/internal/services"
)

type BillingHandler struct {
	service *services.BillingService
}

func NewBillingHandler(s *services.BillingService) *BillingHandler {
	return &BillingHandler{service: s}
}

type billingRequest struct {
	UserID string  `json:"userId"`
	TaskID string  `json:"taskId"`
	Amount float64 `json:"amount"`
}

type billingResponse struct {
	TransactionID string  `json:"transactionId"`
	UserID        string  `json:"userId"`
	TaskID        string  `json:"taskId"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
}

func (h *BillingHandler) Charge(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBillingRequest(w, r)
	if !ok {
		return
	}

	tx, err := h.service.Charge(req.UserID, req.TaskID, req.Amount)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusCreated, toBillingResponse(tx))
}

func (h *BillingHandler) Refund(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBillingRequest(w, r)
	if !ok {
		return
	}

	tx, err := h.service.Refund(req.UserID, req.TaskID, req.Amount)
	if err != nil {
		respondServiceError(w, r, err)
		return
	}

	respondJSON(w, http.StatusCreated, toBillingResponse(tx))
}

func decodeBillingRequest(w http.ResponseWriter, r *http.Request) (billingRequest, bool) {
	var req billingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid request body")
		return req, false
	}

	return req, true
}

func toBillingResponse(tx *models.Transaction) billingResponse {
	return billingResponse{
		TransactionID: tx.ID,
		UserID:        tx.UserID,
		TaskID:        tx.TaskID,
		Amount:        tx.Amount,
		Type:          tx.Type,
	}
}
