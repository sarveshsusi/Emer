package repositories

import (
	"time"

	"github.com/google/uuid"
	"ticketapp/internal/models"
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)

	// 2FA
	GetOTPSecret(userID uuid.UUID) (string, error)

	// Password reset
	StoreResetToken(userID uuid.UUID, hash string, exp time.Time) error
	ValidateResetToken(hash string) (uuid.UUID, error)
	UpdatePassword(userID uuid.UUID, passwordHash string) error
}

type RefreshTokenRepository interface {
	Store(userID uuid.UUID, hash string, exp time.Time) error
	GetValid(hash string) (*models.RefreshToken, error)
	Revoke(tokenID uuid.UUID) error
	RevokeAll(userID uuid.UUID) error
}
