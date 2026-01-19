package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"ticketapp/internal/services"
)

func (h *AuthHandler) issueTokens(
	w http.ResponseWriter,
	userID string,
	role string,
) {
	// generate access token
	accessToken, err := h.jwt.GenerateAccessToken(userID, role)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	// generate refresh token
	refreshToken := uuid.NewString()
	refreshHash := services.HashToken(refreshToken)

	// store refresh token (hashed)
	err = h.tokenRepo.Store(
		uuid.MustParse(userID),
		refreshHash,
		time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		http.Error(w, "token storage failed", http.StatusInternalServerError)
		return
	}

	// set secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
	})

	// response
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

