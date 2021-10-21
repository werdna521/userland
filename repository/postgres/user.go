package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

type UserRepository struct {
	db         *sql.DB
	statements *userQueryStatements
}

type userQueryStatements struct {
	createUserStmt *sql.Stmt
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) PrepareStatements(ctx context.Context) error {
	query := fmt.Sprintf(
		`INSERT INTO %s(%s, %s, %s) 
		 VALUES($1, $2, $3) 
		 RETURNING id`,
		userTableName,
		userTableFullNameColName,
		userTableEmailColName,
		userTablePasswordColName,
	)
	log.Info().Msg("preparing create user statement")
	createUserStmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to prepare create user statement")
		return err
	}

	ur.statements = &userQueryStatements{
		createUserStmt: createUserStmt,
	}

	return nil
}

func (ur *UserRepository) TearDownStatements() {
	defer ur.statements.createUserStmt.Close()
}

func (ur *UserRepository) CreateUser(ctx context.Context, u *repository.User) error {
	var userID int64
	err := ur.statements.createUserStmt.
		QueryRowContext(ctx, u.Fullname, u.Email, u.Password).
		Scan(&userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to execute statement")
		return err
	}

	u.ID = userID

	return nil
}
