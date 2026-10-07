package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/fazil-syed/bifrost/internal/authentication"
)

type Handler struct {
	authenticationService authentication.AuthenticationService
}

func NewHandler(authenticationService authentication.AuthenticationService) *Handler {
	return &Handler{
		authenticationService: authenticationService,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, loginSession, err := h.authenticationService.LoginWithPassword(r.Context(), request.Email, request.Password)
	if err != nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "bifrost_session",
		Value:    loginSession.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  loginSession.ExpiresAt,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(LoginResponse{Success: true})
}
