package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RegisterHandler struct {
	service *Service
}

func NewRegisterHandler(service *Service) *RegisterHandler {
	return &RegisterHandler{
		service: service,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Register user
	if err := h.service.CreateWithLoginAndPassword(r.Context(), req.Login, req.Password); err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
