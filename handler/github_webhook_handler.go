package handler

import (
	"github.com/hayashiki/mentions/usecase"
	"net/http"
)

type WebhookHandler struct {
	process usecase.WebhookProcess
}

func NewWebhookHandler(
	process usecase.WebhookProcess) WebhookHandler {
	return WebhookHandler{
		process: process,
	}
}

func (h *WebhookHandler) PostWebhook(w http.ResponseWriter, r *http.Request) {

	err := h.process.Do(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
