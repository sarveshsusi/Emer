package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticketapp/internal/repositories"
	"ticketapp/internal/services"
	"ticketapp/internal/utils"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.RefreshTokenRepository
	jwt       *services.JWTService
	otp       *services.OTPService
}


func NewAuthHandler(
	userRepo repositories.UserRepository,
	tokenRepo repositories.RefreshTokenRepository,
	jwt *services.JWTService,
	otp *services.OTPService,
) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwt:       jwt,
		otp:       otp,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil || utils.ComparePassword(user.PasswordHash, req.Password) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		fmt.Print(err)
		return
	}

	// 🔐 OPTIONAL 2FA FLOW (keep commented until ready)
	/*
	if user.Is2FAEnabled {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"requires_2fa": true,
			"user_id":      user.ID.String(),
		})
		return
	}
	*/

	// ✅ SINGLE SOURCE OF TRUTH FOR TOKENS
	h.issueTokens(w, user.ID.String(), user.Role)
}


func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	hash := services.HashToken(c.Value)

	token, err := h.tokenRepo.GetValid(hash)
	if err != nil {
		// 🚨 TOKEN REUSE / INVALID TOKEN
		// revoke ALL sessions for this user
		// (hash is untrusted at this point)
		http.Error(w, "token reuse detected", http.StatusUnauthorized)
		return
	}

	// rotate refresh token
	_ = h.tokenRepo.Revoke(token.ID)

	newRefresh := uuid.NewString()
	newHash := services.HashToken(newRefresh)

	_ = h.tokenRepo.Store(
		token.UserID,
		newHash,
		time.Now().Add(7*24*time.Hour),
	)

	user, err := h.userRepo.GetByID(token.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	h.issueTokens(w, user.ID.String(), user.Role)
}


func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	secret, err := h.userRepo.GetOTPSecret(uid)
	if err != nil || !h.otp.Verify(secret, req.Code) {
		http.Error(w, "invalid otp", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(uid)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// ✅ CORRECT CALL — 3 ARGUMENTS
	h.issueTokens(w, user.ID.String(), user.Role)
}


func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email string }
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return
	}

	token := uuid.NewString()
	hash := services.HashToken(token)

	h.userRepo.StoreResetToken(user.ID, hash, time.Now().Add(15*time.Minute))

	// send email with token (link)
	w.WriteHeader(200)
}
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string
		Password string
	}
	json.NewDecoder(r.Body).Decode(&req)

	hash := services.HashToken(req.Token)
	userID, err := h.userRepo.ValidateResetToken(hash)
	if err != nil {
		http.Error(w, "invalid or expired", 400)
		return
	}

	pw, _ := utils.HashPassword(req.Password)
	h.userRepo.UpdatePassword(userID, pw)
	h.tokenRepo.RevokeAll(userID)

	w.WriteHeader(200)
}
