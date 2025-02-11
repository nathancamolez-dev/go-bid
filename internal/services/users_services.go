package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/nathancamolez-dev/go-bid/internal/store/pgstore"
)

var ErrDuplicatedEmailOrUsername = errors.New("duplicated email or username")

type UserServices struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService(pool *pgxpool.Pool) UserServices {
	return UserServices{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserServices) CreateUser(
	ctx context.Context,
	userName,
	password,
	email,
	bio string,
) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}
	args := pgstore.CreateUserParams{
		UserName:     userName,
		PasswordHash: hash,
		Email:        email,
		Bio:          bio,
	}
	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation{
			return uuid.UUID{}, ErrDuplicatedEmailOrUsername
		}
		return uuid.UUID{}, err
	}
	return id, nil
}
