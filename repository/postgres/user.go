package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

const (
	userTableName             = "users"
	userTableIDColName        = "id"
	userTableFullNameColName  = "fullname"
	userTableEmailColName     = "email"
	userTablePasswordColName  = "password"
	userTableIsActiveColName  = "is_active"
	userTableCreatedAtColName = "created_at"
	userTableUpdatedAtColName = "updated_at"
)

type UserRepository interface {
	PrepareStatements(context.Context) error
	TearDownStatements()
	CreateUser(ctx context.Context, user *repository.User) (*repository.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	UpdateUserActivationStatusByEmail(
		ctx context.Context,
		email string,
		isActive bool,
	) (*repository.User, error)
	UpdatePasswordByEmail(
		ctx context.Context,
		email string,
		password string,
	) (*repository.User, error)
}

type BaseUserRepository struct {
	db         *sql.DB
	statements *userStatements
}

type userStatements struct {
	createUserStmt                        *sql.Stmt
	getUserByEmailStmt                    *sql.Stmt
	updateUserActivationStatusByEmailStmt *sql.Stmt
	updatePasswordByEmailStmt             *sql.Stmt
}

func NewUserRepository(db *sql.DB) *BaseUserRepository {
	return &BaseUserRepository{
		db: db,
	}
}

func (r *BaseUserRepository) PrepareStatements(ctx context.Context) error {
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
	createUserStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare create user statement")
		return err
	}

	query = fmt.Sprintf(
		`SELECT * FROM %s
		 WHERE %s = $1`,
		userTableName,
		userTableEmailColName,
	)
	log.Info().Msg("preparing get user by email statement")
	getUserByEmailStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare get user by email statement")
		return err
	}

	query = fmt.Sprintf(
		`UPDATE %s
		 SET %s = $1
		 WHERE %s = $2
		 RETURNING *`,
		userTableName,
		userTableIsActiveColName,
		userTableEmailColName,
	)
	log.Info().Msg("preparing update activation status by email statement")
	updateUserActivationStatusByEmailStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update activation status by email statement")
		return err
	}

	query = fmt.Sprintf(
		`UPDATE %s
		 SET %s = $1
		 WHERE %s = $2
		 RETURNING *`,
		userTableName,
		userTablePasswordColName,
		userTableEmailColName,
	)
	log.Info().Msg("preparing update password by email statement")
	updatePasswordByEmailStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update password by email statement")
	}

	r.statements = &userStatements{
		createUserStmt:                        createUserStmt,
		getUserByEmailStmt:                    getUserByEmailStmt,
		updateUserActivationStatusByEmailStmt: updateUserActivationStatusByEmailStmt,
		updatePasswordByEmailStmt:             updatePasswordByEmailStmt,
	}

	return nil
}

func (r *BaseUserRepository) TearDownStatements() {
	defer r.statements.createUserStmt.Close()
	defer r.statements.getUserByEmailStmt.Close()
	defer r.statements.updateUserActivationStatusByEmailStmt.Close()
	defer r.statements.updatePasswordByEmailStmt.Close()
}

func (r *BaseUserRepository) CreateUser(
	ctx context.Context,
	u *repository.User,
) (*repository.User, error) {
	log.Info().Msg("running statement to create user")
	err := r.statements.createUserStmt.
		QueryRowContext(ctx, u.Fullname, u.Email, u.Password).
		Scan(&u.ID)

	return u, err
}

func (r *BaseUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to get user by email")
	err := r.statements.getUserByEmailStmt.
		QueryRowContext(ctx, email).
		Scan(&u.ID, &u.Fullname, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msg("failed to find a user")
		return nil, repository.NewNotFoundError()
	}

	return u, err
}

func (r *BaseUserRepository) UpdateUserActivationStatusByEmail(
	ctx context.Context,
	email string,
	isActive bool,
) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to update user activation status by email")
	err := r.statements.updateUserActivationStatusByEmailStmt.
		QueryRowContext(ctx, isActive, email).
		Scan(&u.ID, &u.Fullname, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}

func (r *BaseUserRepository) UpdatePasswordByEmail(
	ctx context.Context,
	email string,
	password string,
) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to update password by email")
	err := r.statements.updatePasswordByEmailStmt.
		QueryRowContext(ctx, password, email).
		Scan(&u.ID, &u.Fullname, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	return u, err
}
