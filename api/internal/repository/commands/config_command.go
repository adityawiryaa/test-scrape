package commands

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/adityawiryaa/api/domain/entity"
)

type ConfigCommand struct {
	db *sql.DB
}

func NewConfigCommand(db *sql.DB) *ConfigCommand {
	return &ConfigCommand{db: db}
}

func (r *ConfigCommand) SaveConfig(ctx context.Context, cfg *entity.Config) error {
	data, err := json.Marshal(cfg.Data)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO configs (id, version, data, poll_interval_seconds, created_at) VALUES (?, ?, ?, ?, ?)`,
		cfg.ID, cfg.Version, string(data), cfg.PollIntervalSeconds, cfg.CreatedAt,
	)
	return err
}
