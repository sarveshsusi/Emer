package handlers

import (
	"encoding/json"
	"net/http"

	"ticketapp/internal/models"
	"ticketapp/internal/repositories"
	"ticketapp/internal/services"
	"ticketapp/internal/utils"

	"github.com/google/uuid"
)
type AdminHandler struct {
	userRepo repositories.UserRepository
	emailSvc *services.EmailService
}
func NewAdminHandler(
	userRepo repositories.UserRepository,
	emailSvc *services.EmailService,
) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
		emailSvc: emailSvc,
	}
}
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"` // support | customer
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Role != "support" && req.Role != "customer" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	tempPassword := uuid.NewString()[:10]
	passwordHash, _ := utils.HashPassword(tempPassword)

	user := models.User{
		ID:                    uuid.New(),
		Email:                 req.Email,
		Username:              req.Email,
		PasswordHash:          passwordHash,
		Role:                  req.Role,
		IsActive:              true,
		PasswordResetRequired: true,
		Is2FAEnabled:          false,
	}

	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	// Send onboarding email
	h.emailSvc.SendUserInvite(
		user.Email,
		user.Username,
		tempPassword,
	)

	w.WriteHeader(http.StatusCreated)
}
func (h *AdminHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	uid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := h.userRepo.Disable(uid); err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
