package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/1sh-repalto/url-monitoring-api/internal/model"
)

type pgxURLRepository struct {
	db *pgxpool.Pool
}

func NewPgxURLRepository(db *pgxpool.Pool) *pgxURLRepository {
	return &pgxURLRepository{db}
}

func(r *pgxURLRepository) SaveURL(u *model.MonitoredURL) error {
	query := `INSERT INTO monitored_urls (id, url, user_id, is_active)
			  VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query, u.ID, u.URL, u.UserID, u.IsActive)

	return err
}

func(r *pgxURLRepository) GetURLByUserID(userID int) ([]*model.MonitoredURL, error) {
	query := `SELECT id, url, user_id, is_active, created_at
			  FROM monitored_urls
			  WHERE user_id = $1
	`

	rows, err := r.db.Query(context.Background(), query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*model.MonitoredURL

	for rows.Next() {
		var u model.MonitoredURL

		err := rows.Scan(&u.ID, &u.URL, &u.UserID, &u.IsActive, &u.CreatedAt)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return urls, nil
}

func(r *pgxURLRepository) GetURLByID(id string) (*model.MonitoredURL, error) {
	query := `SELECT id, url, user_id, is_active, created_at FROM monitored_urls WHERE id = $1`

	var url model.MonitoredURL
	err := r.db.QueryRow(context.Background(), query, id).Scan(&url.ID, &url.URL, &url.UserID, &url.IsActive, &url.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func(r *pgxURLRepository) DeleteURL(urlID string) error {
	query := `DELETE FROM monitored_urls WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, urlID)

	return err
}