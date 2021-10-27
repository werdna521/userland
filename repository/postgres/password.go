package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

const (
	passwordHistoryTableName             = "password_history"
	passwordHistoryTableUserIDColName    = "user_id"
	passwordHistoryTablePasswordColName  = "password"
	passwordHistoryTableCreatedAtColName = "created_at"
)

type PasswordHistoryRepository interface {
	PrepareStatements(context.Context) error
	TearDownStatements()
	CreatePasswordHistoryRecord(
		ctx context.Context,
		fp *repository.PasswordHistory,
	) (*repository.PasswordHistory, error)
	GetLastNPasswordHashes(ctx context.Context, userID string, n int) ([]string, error)
}

type BasePasswordHistoryRepository struct {
	db         *sql.DB
	statements *PasswordHistoryStatements
}

type PasswordHistoryStatements struct {
	createPasswordHistoryRecordStmt *sql.Stmt
	getLastNPasswordHashesStmt      *sql.Stmt
}

func NewBasePasswordHistoryRepository(db *sql.DB) *BasePasswordHistoryRepository {
	return &BasePasswordHistoryRepository{
		db: db,
	}
}

func (r *BasePasswordHistoryRepository) PrepareStatements(ctx context.Context) error {
	log.Info().Msg("preparing create forgot password record statement")
	query := fmt.Sprintf(
		`INSERT INTO %s
		 VALUES(DEFAULT, $1, $2, $3, $4)
		 RETURNING id`,
		passwordHistoryTableName,
	)
	createPasswordHistoryRecordStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare create forgot password record statement")
		return err
	}

	log.Info().Msg("preparing get last n password hashes statement")
	query = fmt.Sprintf(
		`SELECT %s
		 FROM %s
		 WHERE %s = $1
		 ORDER BY %s
		 LIMIT $2`,
		passwordHistoryTablePasswordColName,
		passwordHistoryTableName,
		passwordHistoryTableUserIDColName,
		passwordHistoryTableCreatedAtColName,
	)
	getLastNPasswordHashesStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare get last n password hashes statement")
		return err
	}

	r.statements = &PasswordHistoryStatements{
		createPasswordHistoryRecordStmt: createPasswordHistoryRecordStmt,
		getLastNPasswordHashesStmt:      getLastNPasswordHashesStmt,
	}

	return nil
}

func (r *BasePasswordHistoryRepository) TearDownStatements() {
	defer r.statements.createPasswordHistoryRecordStmt.Close()
}

func (r *BasePasswordHistoryRepository) CreatePasswordHistoryRecord(
	ctx context.Context,
	fp *repository.PasswordHistory,
) (*repository.PasswordHistory, error) {
	now := time.Now()

	log.Info().Msg("running statement to create forgot password record")
	err := r.statements.createPasswordHistoryRecordStmt.
		QueryRowContext(ctx, fp.UserID, fp.Password, now, now).
		Scan(&fp.ID)

	return fp, err
}

func (r *BasePasswordHistoryRepository) GetLastNPasswordHashes(
	ctx context.Context,
	userID string,
	n int,
) ([]string, error) {
	log.Info().Msg("running statement to get last n password hashes")
	rows, err := r.statements.getLastNPasswordHashesStmt.QueryContext(ctx, userID, fmt.Sprint(n))
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msg("no password history")
		return nil, repository.NewNotFoundError()
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get last n password hashes")
		return nil, err
	}
	defer rows.Close()

	var hashes []string
	for rows.Next() {
		var hash string
		err := rows.Scan(&hash)
		if err != nil {
			log.Error().Err(err).Msg("fail to scan password hash")
			return nil, err
		}
		hashes = append(hashes, hash)
	}

	return hashes, nil
}
