package user

import (
	"encoding/json"
	"net/http"
	"time"
)

type LogoutHandler struct {
	service *Service
}

func NewLogoutHandler(service *Service) *LogoutHandler {
	return &LogoutHandler{
		service: service,
	}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		// No session cookie found, still return success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	// Delete session from database
	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear session cookie by setting it to expire immediately
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
