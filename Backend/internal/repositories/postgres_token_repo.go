package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"ticketapp/internal/models"
)

type PostgresRefreshTokenRepo struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepo(db *sql.DB) *PostgresRefreshTokenRepo {
	return &PostgresRefreshTokenRepo{db: db}
}
func (r *PostgresRefreshTokenRepo) Store(
	userID uuid.UUID,
	hash string,
	exp time.Time,
) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3)`,
		userID, hash, exp,
	)
	return err
}

func (r *PostgresRefreshTokenRepo) GetValid(hash string) (*models.RefreshToken, error) {
	t := &models.RefreshToken{}
	err := r.db.QueryRow(
		`SELECT id, user_id FROM refresh_tokens
		 WHERE token_hash=$1 AND revoked=false AND expires_at > NOW()`,
		hash,
	).Scan(&t.ID, &t.UserID)
	return t, err
}

func (r *PostgresRefreshTokenRepo) Revoke(tokenID uuid.UUID) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens SET revoked=true WHERE id=$1`,
		tokenID,
	)
	return err
}

func (r *PostgresRefreshTokenRepo) RevokeAll(userID uuid.UUID) error {
	_, err := r.db.Exec(
		`UPDATE refresh_tokens SET revoked=true WHERE user_id=$1`,
		userID,
	)
	return err
}
