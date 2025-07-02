package repository

import (
	"context"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxURLRepository struct {
	db *pgxpool.Pool
}

func NewPgxURLRepository(db *pgxpool.Pool) *pgxURLRepository {
	return &pgxURLRepository{db}
}

func (r *pgxURLRepository) SaveURL(u *model.MonitoredURL) error {
	query := `INSERT INTO monitored_urls (id, url, user_id, is_active)
			  VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query, u.ID, u.URL, u.UserID, u.IsActive)

	return err
}

func (r *pgxURLRepository) ExistsByUserAndURL(ctx context.Context, url string, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM monitored_urls WHERE user_id = $1 AND url = $2)`
	err := r.db.QueryRow(ctx, query, userID, url).Scan(&exists)
	return exists, err
}

func (r *pgxURLRepository) GetURLByUserID(userID int) ([]*model.MonitoredURL, error) {
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

func (r *pgxURLRepository) GetURLByID(id string) (*model.MonitoredURL, error) {
	query := `SELECT id, url, user_id, is_active, created_at FROM monitored_urls WHERE id = $1`

	var url model.MonitoredURL
	err := r.db.QueryRow(context.Background(), query, id).Scan(&url.ID, &url.URL, &url.UserID, &url.IsActive, &url.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *pgxURLRepository) DeleteURL(urlID string) error {
	query := `DELETE FROM monitored_urls WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, urlID)

	return err
}

func (r *pgxURLRepository) GetAllActiveURLs() ([]*model.MonitoredURL, error) {
	query := `SELECT id, url, user_id, is_active, created_at FROM monitored_urls WHERE is_active = true`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*model.MonitoredURL

	for rows.Next() {
		var u model.MonitoredURL
		if err := rows.Scan(&u.ID, &u.URL, &u.UserID, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, &u)
	}

	return urls, rows.Err()
}

func (r *pgxURLRepository) SaveURLLog(log *model.URLLog) error {
	query := `
		INSERT INTO url_logs (id, url_id, status_code, response_time_ms, checked_at, is_up)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(context.Background(), query, log.ID, log.URLID, log.StatusCode, log.ResponseTimeMs, log.CheckedAt, log.IsUp)

	return err
}

func (r *pgxURLRepository) GetLogsByURLID(urlID string) ([]*model.URLLog, error) {
	query := `
		SELECT id, url_id, status_code, response_time_ms, checked_at, is_up
		FROM url_logs
		WHERE url_id = $1
		ORDER BY checked_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.URLLog
	for rows.Next() {
		var log model.URLLog
		err := rows.Scan(&log.ID, &log.URLID, &log.StatusCode, &log.ResponseTimeMs, &log.CheckedAt, &log.IsUp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, rows.Err()
}
