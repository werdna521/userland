package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, u *repository.User) error {
	query := fmt.Sprintf(
		`INSERT INTO %s(%s, %s, %s) 
		 VALUES($1, $2, $3) 
		 RETURNING id`,
		userTableName,
		userTableFullNameColName,
		userTableEmailColName,
		userTablePasswordColName,
	)

	stmt, err := ur.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to prepare statement")
		return err
	}
	defer stmt.Close()

	var userID int64
	err = stmt.QueryRowContext(ctx, u.FullName, u.Email, u.Password).Scan(&userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to execute statement")
		return err
	}

	u.ID = userID

	return nil
}
