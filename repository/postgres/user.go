package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/repository"
)

const (
	userTableName             = `"user"`
	userTableIDColName        = "id"
	userTableEmailColName     = "email"
	userTablePasswordColName  = "password"
	userTableIsActiveColName  = "is_active"
	userTableCreatedAtColName = "created_at"
	userTableUpdatedAtColName = "updated_at"

	userBioTableName             = "user_bio"
	userBioTableIDColName        = "id"
	userBioTableUserIDColName    = "user_id"
	userBioTableFullNameColName  = "fullname"
	userBioTableLocationColName  = "location"
	userBioTableBioColName       = "bio"
	userBioTableWebColName       = "web"
	userBioTablePictureColName   = "picture"
	userBioTableCreatedAtColName = "created_at"
	userBioTableUpdatedAtColName = "updated_at"
)

type UserRepository interface {
	PrepareStatements(context.Context) error
	TearDownStatements()
	CreateUser(ctx context.Context, user *repository.User) (*repository.User, error)
	GetUserByID(ctx context.Context, userID string) (*repository.User, error)
	GetUserBioByID(ctx context.Context, userID string) (*repository.UserBio, error)
	GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
	UpdateUserActivationStatusByID(
		ctx context.Context,
		userID string,
		isActive bool,
	) (*repository.User, error)
	UpdatePasswordByID(
		ctx context.Context,
		userID string,
		password string,
	) (*repository.User, error)
	UpdateEmailByID(
		ctx context.Context,
		userID string,
		email string,
	) (*repository.User, error)
	UpdateUserBioByID(
		ctx context.Context,
		userID string,
		ub *repository.UserBio,
	) (*repository.UserBio, error)
	UpdatePictureByID(
		ctx context.Context,
		userID string,
		picturePath string,
	) (*repository.UserBio, error)
}

type BaseUserRepository struct {
	db         *sql.DB
	statements *userStatements
}

type userStatements struct {
	createUserStmt                     *sql.Stmt
	createUserBioStmt                  *sql.Stmt
	getUserByIDStmt                    *sql.Stmt
	getUserByEmailStmt                 *sql.Stmt
	getUserBioByIDStmt                 *sql.Stmt
	updateUserActivationStatusByIDStmt *sql.Stmt
	updatePasswordByIDStmt             *sql.Stmt
	updateEmailByIDStmt                *sql.Stmt
	updateUserBioByIDStmt              *sql.Stmt
	updatePictureByIDStmt              *sql.Stmt
}

func NewBaseUserRepository(db *sql.DB) *BaseUserRepository {
	return &BaseUserRepository{
		db: db,
	}
}

func (r *BaseUserRepository) scanUser(u *repository.User, row *sql.Row) error {
	return row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}

func (r *BaseUserRepository) scanUserBio(ub *repository.UserBio, row *sql.Row) error {
	return row.Scan(
		&ub.ID,
		&ub.Fullname,
		&ub.Location,
		&ub.Bio,
		&ub.Web,
		&ub.Picture,
		&ub.CreatedAt,
		&ub.UpdatedAt,
	)
}

func (r *BaseUserRepository) PrepareStatements(ctx context.Context) error {
	log.Info().Msg("preparing create user statement")
	query := fmt.Sprintf(
		`INSERT INTO %s
		 VALUES(DEFAULT, $1, $2, $3, $4, $5)
		 RETURNING id`,
		userTableName,
	)
	createUserStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare create user statement")
		return err
	}

	log.Info().Msg("preparing create user bio statement")
	query = fmt.Sprintf(
		`INSERT INTO %s
		 VALUES(DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		userBioTableName,
	)
	createUserBioStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare create user bio statement")
		return err
	}

	log.Info().Msg("preparing get user by ID statement")
	query = fmt.Sprintf(
		`SELECT *
		 FROM %s
		 WHERE %s = $1`,
		userTableName,
		userTableIDColName,
	)
	getUserByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare get user by ID statement")
		return err
	}

	log.Info().Msg("preparing get user by email statement")
	query = fmt.Sprintf(
		`SELECT * 
		 FROM %s
		 WHERE %s = $1`,
		userTableName,
		userTableEmailColName,
	)
	getUserByEmailStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare get user by email statement")
		return err
	}

	log.Info().Msg("preparing get user bio by id statement")
	query = fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s
		 FROM %s
		 WHERE %s = $1`,
		userBioTableIDColName,
		userBioTableFullNameColName,
		userBioTableLocationColName,
		userBioTableBioColName,
		userBioTableWebColName,
		userBioTablePictureColName,
		userBioTableCreatedAtColName,
		userBioTableUpdatedAtColName,
		userBioTableName,
		userBioTableUserIDColName,
	)
	getUserBioByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare get user bio by id statement")
		return err
	}

	log.Info().Msg("preparing update activation status by email statement")
	query = fmt.Sprintf(
		`UPDATE %s
		 SET 
		   %s = $1,
			 %s = $2
		 WHERE %s = $3
		 RETURNING *`,
		userTableName,
		userTableIsActiveColName,
		userTableUpdatedAtColName,
		userTableIDColName,
	)
	updateUserActivationStatusByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update activation status by email statement")
		return err
	}

	log.Info().Msg("preparing update password by email statement")
	query = fmt.Sprintf(
		`UPDATE %s
		 SET 
		   %s = $1,
			 %s = $2
		 WHERE %s = $3
		 RETURNING *`,
		userTableName,
		userTablePasswordColName,
		userTableUpdatedAtColName,
		userTableIDColName,
	)
	UpdatePasswordByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update password by email statement")
		return err
	}

	log.Info().Msg("preparing update email by id statement")
	query = fmt.Sprintf(
		`UPDATE %s
		 SET
		   %s = $1,
			 %s = $2
		 WHERE %s = $3
		 RETURNING *`,
		userTableName,
		userTableEmailColName,
		userTableUpdatedAtColName,
		userTableIDColName,
	)
	updateEmailByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update email by id statement")
		return err
	}

	log.Info().Msg("preparing update user bio by id statement")
	query = fmt.Sprintf(
		`UPDATE %s
		 SET 
		   %s = $1,
			 %s = CASE
							WHEN $2 = '' THEN %s
							ELSE $2
						END,
			 %s = CASE
							WHEN $3 = '' THEN %s
							ELSE $3
						END,
			 %s = CASE
							WHEN $4 = '' THEN %s
							ELSE $4
						END,
			 %s = $5
		 WHERE %s = $6
		 RETURNING %s, %s, %s, %s, %s, %s, %s, %s`,
		userBioTableName,
		userBioTableFullNameColName,
		userBioTableLocationColName,
		userBioTableLocationColName,
		userBioTableBioColName,
		userBioTableBioColName,
		userBioTableWebColName,
		userBioTableWebColName,
		userBioTableUpdatedAtColName,
		userBioTableUserIDColName,
		userBioTableIDColName,
		userBioTableFullNameColName,
		userBioTableLocationColName,
		userBioTableBioColName,
		userBioTableWebColName,
		userBioTablePictureColName,
		userBioTableCreatedAtColName,
		userBioTableUpdatedAtColName,
	)
	updateUserBioByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update user bio by id statement")
		return err
	}

	log.Info().Msg("preparing update picture by id statement")
	query = fmt.Sprintf(
		`UPDATE %s
		 SET 
		   %s = $1,
			 %s = $2
		 WHERE %s = $3
		 RETURNING %s, %s, %s, %s, %s, %s, %s, %s`,
		userBioTableName,
		userBioTablePictureColName,
		userBioTableUpdatedAtColName,
		userBioTableUserIDColName,
		userBioTableIDColName,
		userBioTableFullNameColName,
		userBioTableLocationColName,
		userBioTableBioColName,
		userBioTableWebColName,
		userBioTablePictureColName,
		userBioTableCreatedAtColName,
		userBioTableUpdatedAtColName,
	)
	updatePictureByIDStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare update picture by id statement")
		return err
	}

	r.statements = &userStatements{
		createUserStmt:                     createUserStmt,
		createUserBioStmt:                  createUserBioStmt,
		getUserByIDStmt:                    getUserByIDStmt,
		getUserByEmailStmt:                 getUserByEmailStmt,
		getUserBioByIDStmt:                 getUserBioByIDStmt,
		updateUserActivationStatusByIDStmt: updateUserActivationStatusByIDStmt,
		updatePasswordByIDStmt:             UpdatePasswordByIDStmt,
		updateEmailByIDStmt:                updateEmailByIDStmt,
		updateUserBioByIDStmt:              updateUserBioByIDStmt,
		updatePictureByIDStmt:              updatePictureByIDStmt,
	}

	return nil
}

func (r *BaseUserRepository) TearDownStatements() {
	defer r.statements.createUserStmt.Close()
	defer r.statements.getUserByEmailStmt.Close()
	defer r.statements.updateUserActivationStatusByIDStmt.Close()
	defer r.statements.updatePasswordByIDStmt.Close()
}

func (r *BaseUserRepository) CreateUser(
	ctx context.Context,
	u *repository.User,
) (*repository.User, error) {
	now := time.Now()

	log.Info().Msg("running statement to create user")
	err := r.statements.createUserStmt.
		QueryRowContext(ctx, u.Email, u.Password, u.IsActive, now, now).
		Scan(&u.ID)

	if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
		log.Error().Err(err).Msg("violated unique email constraint")
		return nil, repository.NewUniqueViolationError()
	}

	log.Info().Msg("running statement to create user bio")
	err = r.statements.createUserBioStmt.
		QueryRowContext(ctx, u.ID, u.UserBio.Fullname, "", "", "", "", now, now).
		Scan(&u.UserBio.ID)

	return u, err
}

func (r *BaseUserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to get user by id")
	row := r.statements.getUserByIDStmt.QueryRowContext(ctx, userID)
	err := r.scanUser(u, row)
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msg("failed to find user")
		return nil, repository.NewNotFoundError()
	}

	return u, err
}

func (r *BaseUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*repository.User, error) {
	u := &repository.User{}

	log.Info().Msg("running statement to get user by email")
	row := r.statements.getUserByEmailStmt.QueryRowContext(ctx, email)
	err := r.scanUser(u, row)
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msg("failed to find a user")
		return nil, repository.NewNotFoundError()
	}

	return u, err
}

func (r *BaseUserRepository) GetUserBioByID(
	ctx context.Context,
	userID string,
) (*repository.UserBio, error) {
	ub := &repository.UserBio{}

	log.Info().Msg("running statement to get user bio by id")
	row := r.statements.getUserBioByIDStmt.QueryRowContext(ctx, userID)
	err := r.scanUserBio(ub, row)
	if err == sql.ErrNoRows {
		log.Error().Err(err).Msg("failed to find user bio")
		return nil, repository.NewNotFoundError()
	}

	return ub, err
}

func (r *BaseUserRepository) UpdateUserActivationStatusByID(
	ctx context.Context,
	userID string,
	isActive bool,
) (*repository.User, error) {
	u := &repository.User{}
	now := time.Now()

	log.Info().Msg("running statement to update user activation status by email")
	row := r.statements.updateUserActivationStatusByIDStmt.QueryRowContext(ctx, isActive, now, userID)
	err := r.scanUser(u, row)

	return u, err
}

func (r *BaseUserRepository) UpdatePasswordByID(
	ctx context.Context,
	userID string,
	password string,
) (*repository.User, error) {
	u := &repository.User{}
	now := time.Now()

	log.Info().Msg("running statement to update password by email")
	row := r.statements.updatePasswordByIDStmt.QueryRowContext(ctx, password, now, userID)
	err := r.scanUser(u, row)

	return u, err
}

func (r *BaseUserRepository) UpdateEmailByID(
	ctx context.Context,
	userID string,
	email string,
) (*repository.User, error) {
	u := &repository.User{}
	now := time.Now()

	log.Info().Msg("running statement to update email by id")
	row := r.statements.updateEmailByIDStmt.QueryRowContext(ctx, email, now, userID)
	err := r.scanUser(u, row)

	return u, err
}

func (r *BaseUserRepository) UpdateUserBioByID(
	ctx context.Context,
	userID string,
	ub *repository.UserBio,
) (*repository.UserBio, error) {
	now := time.Now()

	log.Info().Msg("running statement to update user bio by id")
	row := r.statements.updateUserBioByIDStmt.
		QueryRowContext(ctx, ub.Fullname, ub.Location, ub.Bio, ub.Web, now, userID)
	err := r.scanUserBio(ub, row)

	return ub, err
}

func (r *BaseUserRepository) UpdatePictureByID(
	ctx context.Context,
	userID string,
	picturePath string,
) (*repository.UserBio, error) {
	ub := &repository.UserBio{}
	now := time.Now()

	log.Info().Msg("running statement to update picture by id")
	row := r.statements.updatePictureByIDStmt.QueryRowContext(ctx, picturePath, now, userID)
	err := r.scanUserBio(ub, row)

	return ub, err
}
