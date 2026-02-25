package queries

import (
	"context"
	"database/sql"

	"github.com/adityawiryaa/api/domain/entity"
)

type AgentQuery struct {
	db *sql.DB
}

func NewAgentQuery(db *sql.DB) *AgentQuery {
	return &AgentQuery{db: db}
}

func (r *AgentQuery) FindByID(ctx context.Context, id string) (*entity.Agent, error) {
	agent := &entity.Agent{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, hostname, ip_address, port, status, created_at, updated_at FROM agents WHERE id = ?`, id,
	).Scan(&agent.ID, &agent.Hostname, &agent.IPAddress, &agent.Port, &agent.Status, &agent.CreatedAt, &agent.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return agent, nil
}
