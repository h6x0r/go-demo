package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/h6x0r/go-demo/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (id, firstname, lastname, email, age, created)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Firstname, user.Lastname, user.Email, user.Age, user.Created)
	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, firstname, lastname, email, age, created FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Age, &user.Created)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user models.User) error {
	query := `UPDATE users
              SET firstname = $2, lastname = $3, email = $4, age = $5
              WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Firstname, user.Lastname, user.Email, user.Age)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
