package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/1sh-repalto/url-monitoring-api/internal/model"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *pgxUserRepository {
	return &pgxUserRepository{db}
}

func (r *pgxUserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	user := &model.User{}

	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *pgxUserRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`

	_, err := r.db.Exec(context.Background(), query, user.Email, user.Password)

	return err
}