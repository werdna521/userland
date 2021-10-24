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
	forgotPasswordTableName = "forgot_passwords"
)

type ForgotPasswordRepository interface {
	PrepareStatements(context.Context) error
	TearDownStatements()
	CreateForgotPasswordRecord(
		ctx context.Context,
		fp *repository.ForgotPassword,
	) (*repository.ForgotPassword, error)
}

type BaseForgotPasswordRepository struct {
	db         *sql.DB
	statements *forgotPasswordStatements
}

type forgotPasswordStatements struct {
	createForgotPasswordRecordStmt *sql.Stmt
}

func NewBaseForgotPasswordRepository(db *sql.DB) *BaseForgotPasswordRepository {
	return &BaseForgotPasswordRepository{
		db: db,
	}
}

func (r *BaseForgotPasswordRepository) PrepareStatements(ctx context.Context) error {
	query := fmt.Sprintf(
		`INSERT INTO %s
		 VALUES(DEFAULT, $1, $2, $3, $4)
		 RETURNING id`,
		forgotPasswordTableName,
	)
	log.Info().Msg("preparing create forgot password record statement")
	createForgotPasswordRecordStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare create forgot password record statement")
		return err
	}

	r.statements = &forgotPasswordStatements{
		createForgotPasswordRecordStmt: createForgotPasswordRecordStmt,
	}

	return nil
}

func (r *BaseForgotPasswordRepository) TearDownStatements() {
	defer r.statements.createForgotPasswordRecordStmt.Close()
}

func (r *BaseForgotPasswordRepository) CreateForgotPasswordRecord(
	ctx context.Context,
	fp *repository.ForgotPassword,
) (*repository.ForgotPassword, error) {
	now := time.Now()

	log.Info().Msg("running statement to create forgot password record")
	err := r.statements.createForgotPasswordRecordStmt.
		QueryRowContext(ctx, fp.UserID, fp.OldPassword, now, now).
		Scan(&fp.ID)

	return fp, err
}
