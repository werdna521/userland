package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/log/repository"
)

const (
	auditLogTableName = "audit_log"
)

type LogRepository interface {
	PrepareStatements(ctx context.Context) error
	CreateAuditLog(
		ctx context.Context,
		l *repository.AuditLog,
	) (*repository.AuditLog, error)
}

type BaseLogRepository struct {
	db         *sql.DB
	statements *auditLogStatements
}

type auditLogStatements struct {
	createAuditLogStmt *sql.Stmt
}

func NewBaseLogRepository(db *sql.DB) *BaseLogRepository {
	return &BaseLogRepository{
		db: db,
	}
}

func (r *BaseLogRepository) PrepareStatements(ctx context.Context) error {
	log.Info().Msg("preparing store audit log statement")
	query := fmt.Sprintf(
		`INSERT INTO %s
		 VALUES(DEFAULT, $1, $2, $3, $4, $5)
		 RETURNING id`,
		auditLogTableName,
	)
	createAuditLogStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare store audit log statement")
		return err
	}

	statements := &auditLogStatements{
		createAuditLogStmt: createAuditLogStmt,
	}
	r.statements = statements
	return nil
}

func (r *BaseLogRepository) CreateAuditLog(
	ctx context.Context,
	l *repository.AuditLog,
) (*repository.AuditLog, error) {
	now := time.Now()

	log.Info().Msg("running statement to create audit log")
	err := r.statements.createAuditLogStmt.
		QueryRowContext(ctx, l.UserID, l.RemoteIP, l.AuditType, now, now).
		Scan(&l.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create audit log")
		return nil, err
	}

	return l, nil
}
