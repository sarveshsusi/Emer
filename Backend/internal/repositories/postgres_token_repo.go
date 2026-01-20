package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ticketapp/internal/models"
)

type PostgresRefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRefreshTokenRepo(db *pgxpool.Pool) *PostgresRefreshTokenRepo {
	return &PostgresRefreshTokenRepo{db: db}
}

func (r *PostgresRefreshTokenRepo) Store(
	userID uuid.UUID,
	hash string,
	exp time.Time,
) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1,$2,$3)`,
		userID, hash, exp,
	)
	return err
}


func (r *PostgresRefreshTokenRepo) GetValid(hash string) (*models.RefreshToken, error) {
	t := &models.RefreshToken{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, user_id FROM refresh_tokens
		 WHERE token_hash=$1 AND revoked=false AND expires_at > NOW()`,
		hash,
	).Scan(&t.ID, &t.UserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return t, nil
}


func (r *PostgresRefreshTokenRepo) Revoke(tokenID uuid.UUID) error {
	cmd, err := r.db.Exec(
		context.Background(),
		`UPDATE refresh_tokens SET revoked=true WHERE id=$1`,
		tokenID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("token not found")
	}
	return nil
}


func (r *PostgresRefreshTokenRepo) RevokeAll(userID uuid.UUID) error {
	_, err := r.db.Exec(
		context.Background(),
		`UPDATE refresh_tokens SET revoked=true WHERE user_id=$1`,
		userID,
	)
	return err
}
