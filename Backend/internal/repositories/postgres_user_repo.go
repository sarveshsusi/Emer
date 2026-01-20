package repositories

import (
	"context"
	"errors"
	"time"

	"ticketapp/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepo(db *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetByUsername(username string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRow(
		context.Background(),
   `SELECT id, username, email, password_hash, role, is_active
		 FROM users WHERE email=$1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.IsActive)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return u, nil
}


func (r *PostgresUserRepo) GetByID(id uuid.UUID) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, username, password_hash, role, is_active
		 FROM users WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.IsActive)

	if err != nil {
		return nil, err
	}

	return u, nil
}


func (r *PostgresUserRepo) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, email, password_hash, role, is_active
		 FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresUserRepo) GetOTPSecret(userID uuid.UUID) (string, error) {
	var secret string

	err := r.db.QueryRow(
		context.Background(),
		`SELECT secret FROM otp_secrets WHERE user_id=$1`,
		userID,
	).Scan(&secret)

	return secret, err
}

func (r *PostgresUserRepo) StoreResetToken(
	userID uuid.UUID,
	hash string,
	exp time.Time,
) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO password_resets (user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3)`,
		userID, hash, exp,
	)
	return err
}


func (r *PostgresUserRepo) ValidateResetToken(hash string) (uuid.UUID, error) {
	var userID uuid.UUID

	err := r.db.QueryRow(
		context.Background(),
		`SELECT user_id FROM password_resets
		 WHERE token_hash=$1 AND expires_at > NOW()`,
		hash,
	).Scan(&userID)

	return userID, err
}


func (r *PostgresUserRepo) UpdatePassword(
	userID uuid.UUID,
	passwordHash string,
) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET password_hash=$1 WHERE id=$2`,
		passwordHash, userID,
	)
	return err
}

func (r *PostgresUserRepo) Create(user models.User) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users (id, email, username, password_hash, role, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	)
	return err
}


func (r *PostgresUserRepo) Disable(userID uuid.UUID) error {
	cmd, err := r.db.Exec(
		context.Background(),
		`UPDATE users SET is_active=false WHERE id=$1`,
		userID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

