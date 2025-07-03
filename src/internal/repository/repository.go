package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/entities"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

var _ service.Repository = (*Repository)(nil)

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS(SELECT FROM app_user.users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query row: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *entities.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetRefreshToken(ctx context.Context, id uuid.UUID) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) UpdateRefreshToken(ctx context.Context, token string) error {
	//TODO implement me
	panic("implement me")
}
