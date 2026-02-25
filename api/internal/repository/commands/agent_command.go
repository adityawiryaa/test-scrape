package commands

import (
	"context"
	"database/sql"

	"github.com/adityawiryaa/api/domain/entity"
)

type AgentCommand struct {
	db *sql.DB
}

func NewAgentCommand(db *sql.DB) *AgentCommand {
	return &AgentCommand{db: db}
}

func (r *AgentCommand) Save(ctx context.Context, agent *entity.Agent) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO agents (id, hostname, ip_address, port, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET hostname=?, ip_address=?, port=?, status=?, updated_at=?`,
		agent.ID, agent.Hostname, agent.IPAddress, agent.Port, agent.Status, agent.CreatedAt, agent.UpdatedAt,
		agent.Hostname, agent.IPAddress, agent.Port, agent.Status, agent.UpdatedAt,
	)
	return err
}
