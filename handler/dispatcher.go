package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"workorder-api/model"
	"workorder-api/queue"
	"workorder-api/utils"
)

type DispatcherHandler struct {
	publisher *queue.Publisher
}

func NewDispatcherHandler(publisher *queue.Publisher) *DispatcherHandler {
	return &DispatcherHandler{
		publisher: publisher,
	}
}

// Main dispatcher - route berdasarkan function name
func (h *DispatcherHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Generate request ID untuk tracking
	requestID := utils.GenerateRequestID()

	var genericReq model.GenericRequest
	if err := json.NewDecoder(r.Body).Decode(&genericReq); err != nil {
		log.Printf("[%s] ERROR: Failed to decode request body: %v", requestID, err)
		respondError(w, requestID, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	genericReq.RequestID = requestID
	log.Printf("[%s] INFO: Received function: %s", requestID, genericReq.Function)

	// Validate function name
	if !isValidFunction(genericReq.Function) {
		log.Printf("[%s] WARN: Unknown function: %s", requestID, genericReq.Function)
		respondError(w, requestID, http.StatusBadRequest, "Unknown function: "+genericReq.Function)
		return
	}

	// Publish to queue (async processing)
	msg := &queue.Message{
		RequestID: requestID,
		Function:  genericReq.Function,
		Payload:   genericReq.Payload,
		Timestamp: time.Now(),
	}

	if err := h.publisher.Publish(msg); err != nil {
		log.Printf("[%s] ERROR: Failed to publish message: %v", requestID, err)
		respondError(w, requestID, http.StatusInternalServerError, "Failed to queue request")
		return
	}

	// Return immediate response (202 Accepted)
	response := model.APIResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "Request accepted and queued for processing",
		Data: map[string]interface{}{
			"function": genericReq.Function,
			"status":   "queued",
		},
	}

	w.WriteHeader(http.StatusAccepted) // 202 Accepted
	json.NewEncoder(w).Encode(response)
}

func isValidFunction(function string) bool {
	validFunctions := map[string]bool{
		"ff_updateWorkorder": true,
		// Tambahkan function lain di sini
	}
	return validFunctions[function]
}

func respondError(w http.ResponseWriter, requestID string, code int, message string) {
	response := model.APIResponse{
		RequestID: requestID,
		Success:   false,
		Message:   "Request failed",
		Error:     message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
