package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, username, passwordHash string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1,$2) RETURNING id, username, password_hash, created_at",
		username, passwordHash,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		"SELECT id, username, password_hash, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		"SELECT id, username, password_hash, created_at FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
