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
	createUserStmt     *sql.Stmt
	getUserByEmailStmt *sql.Stmt
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

	query = fmt.Sprintf(
		`SELECT * FROM %s
		 WHERE %s = $1`,
		userTableName,
		userTableEmailColName,
	)
	log.Info().Msg("preparing get user by email statement")
	getUserByEmailStmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to prepare get user by email statement")
		return err
	}

	ur.statements = &userQueryStatements{
		createUserStmt:     createUserStmt,
		getUserByEmailStmt: getUserByEmailStmt,
	}

	return nil
}

func (ur *UserRepository) TearDownStatements() {
	defer ur.statements.createUserStmt.Close()
}

func (ur *UserRepository) CreateUser(ctx context.Context, u *repository.User) error {
	log.Info().Msg("running statement to create user")
	err := ur.statements.createUserStmt.
		QueryRowContext(ctx, u.Fullname, u.Email, u.Password).
		Scan(&u.ID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to execute statement")
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to get user by email")
	err := ur.statements.getUserByEmailStmt.
		QueryRowContext(ctx, email).
		Scan(&u.ID, &u.Fullname, &u.Email, &u.Password, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Error().Stack().Err(err).Msg("failed to find a user")
		return nil, repository.NewNotFoundError()
	}
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to execute statement")
		return nil, err
	}

	return u, nil
}
