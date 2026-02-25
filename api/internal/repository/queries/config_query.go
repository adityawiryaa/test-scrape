package queries

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/adityawiryaa/api/domain/entity"
)

type ConfigQuery struct {
	db *sql.DB
}

func NewConfigQuery(db *sql.DB) *ConfigQuery {
	return &ConfigQuery{db: db}
}

func (r *ConfigQuery) GetLatestConfig(ctx context.Context) (*entity.Config, error) {
	cfg := &entity.Config{}
	var data string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, version, data, poll_interval_seconds, created_at FROM configs ORDER BY version DESC LIMIT 1`,
	).Scan(&cfg.ID, &cfg.Version, &data, &cfg.PollIntervalSeconds, &cfg.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &cfg.Data); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (r *ConfigQuery) GetConfigByVersion(ctx context.Context, version int64) (*entity.Config, error) {
	cfg := &entity.Config{}
	var data string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, version, data, poll_interval_seconds, created_at FROM configs WHERE version = ?`, version,
	).Scan(&cfg.ID, &cfg.Version, &data, &cfg.PollIntervalSeconds, &cfg.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &cfg.Data); err != nil {
		return nil, err
	}
	return cfg, nil
}
