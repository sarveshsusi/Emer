package repositories

import (
	"database/sql"

	"github.com/google/uuid"
)
type AuditRepo struct {
	db *sql.DB
}
func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}
func (a *AuditRepo) Log(
	userID uuid.UUID,
	action string,
	ip string,
	userAgent string,
) error {
	_, err := a.db.Exec(
		`INSERT INTO audit_logs (user_id, action, ip, user_agent)
		 VALUES ($1, $2, $3, $4)`,
		userID,
		action,
		ip,
		userAgent,
	)
	return err
}

