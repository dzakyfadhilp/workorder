package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"workorder-api/model"
	"workorder-api/repository"
	"workorder-api/utils"
)

type DispatcherHandler struct {
	workorderRepo *repository.WorkorderRepository
}

func NewDispatcherHandler(workorderRepo *repository.WorkorderRepository) *DispatcherHandler {
	return &DispatcherHandler{
		workorderRepo: workorderRepo,
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

	// Route ke handler berdasarkan function name
	switch genericReq.Function {
	case "ff_updateWorkorder":
		h.handleUpdateWorkorder(w, &genericReq)
	default:
		log.Printf("[%s] WARN: Unknown function: %s", requestID, genericReq.Function)
		respondError(w, requestID, http.StatusBadRequest, "Unknown function: "+genericReq.Function)
	}
}

// Handler untuk ff_updateWorkorder
func (h *DispatcherHandler) handleUpdateWorkorder(w http.ResponseWriter, genericReq *model.GenericRequest) {
	var payload model.WorkorderRequest
	if err := json.Unmarshal(genericReq.Payload, &payload); err != nil {
		log.Printf("[%s] ERROR: Failed to parse ff_updateWorkorder payload: %v", genericReq.RequestID, err)
		respondError(w, genericReq.RequestID, http.StatusBadRequest, "Invalid payload for ff_updateWorkorder")
		return
	}

	// Validasi
	if err := validateWorkorderRequest(&payload); err != nil {
		log.Printf("[%s] WARN: Validation failed: %v", genericReq.RequestID, err)
		respondError(w, genericReq.RequestID, http.StatusBadRequest, err.Error())
		return
	}

	// Save ke table workorder_updates
	if err := h.workorderRepo.UpsertWorkorder(&payload, genericReq.RequestID); err != nil {
		log.Printf("[%s] ERROR: Failed to save workorder: %v", genericReq.RequestID, err)
		respondError(w, genericReq.RequestID, http.StatusInternalServerError, "Failed to save workorder")
		return
	}

	log.Printf("[%s] INFO: Workorder updated successfully - wonum: %s, status: %s", 
		genericReq.RequestID, payload.Req.Wonum, payload.Req.Status)

	response := model.APIResponse{
		RequestID: genericReq.RequestID,
		Success:   true,
		Message:   "Workorder updated successfully",
		Data: map[string]string{
			"function": "ff_updateWorkorder",
			"wonum":    payload.Req.Wonum,
			"status":   payload.Req.Status,
			"siteid":   payload.Req.Siteid,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func validateWorkorderRequest(req *model.WorkorderRequest) error {
	if req.Req.Wonum == "" {
		return &ValidationError{Field: "wonum", Message: "wonum is required"}
	}
	if req.Req.Status == "" {
		return &ValidationError{Field: "status", Message: "status is required"}
	}
	if req.Req.Siteid == "" {
		return &ValidationError{Field: "siteid", Message: "siteid is required"}
	}
	return nil
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

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
