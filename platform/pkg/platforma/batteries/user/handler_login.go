package user

import (
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	service *Service
}

func NewLoginHandler(service *Service) *LoginHandler {
	return &LoginHandler{
		service: service,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.service.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
