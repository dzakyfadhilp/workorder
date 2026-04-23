package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"workorder-api/model"
	"workorder-api/repository"
)

type WorkorderHandler struct {
	repo *repository.WorkorderRepository
}

func NewWorkorderHandler(repo *repository.WorkorderRepository) *WorkorderHandler {
	return &WorkorderHandler{repo: repo}
}

func (h *WorkorderHandler) UpdateWorkorder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.WorkorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := validateRequest(&req); err != nil {
		log.Printf("WARN: Validation failed: %v", err)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.UpsertWorkorder(&req); err != nil {
		log.Printf("ERROR: Failed to save workorder: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to save workorder")
		return
	}

	log.Printf("INFO: Workorder updated successfully - wonum: %s, status: %s", req.Req.Wonum, req.Req.Status)
	
	response := model.APIResponse{
		Success: true,
		Message: "Workorder updated successfully",
		Data: map[string]string{
			"wonum":  req.Req.Wonum,
			"status": req.Req.Status,
			"siteid": req.Req.Siteid,
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func validateRequest(req *model.WorkorderRequest) error {
	if strings.TrimSpace(req.Req.Wonum) == "" {
		return &ValidationError{Field: "wonum", Message: "wonum is required"}
	}
	if strings.TrimSpace(req.Req.Status) == "" {
		return &ValidationError{Field: "status", Message: "status is required"}
	}
	if strings.TrimSpace(req.Req.Siteid) == "" {
		return &ValidationError{Field: "siteid", Message: "siteid is required"}
	}
	return nil
}

func respondError(w http.ResponseWriter, code int, message string) {
	response := model.APIResponse{
		Success: false,
		Message: "Request failed",
		Error:   message,
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
