package mapper

import (
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/entity"
)

func ToAgentDTO(agent *entity.Agent) dto.AgentDTO {
	return dto.AgentDTO{
		ID:        agent.ID,
		Hostname:  agent.Hostname,
		IPAddress: agent.IPAddress,
		Port:      agent.Port,
		Status:    agent.Status,
	}
}

func ToRegistrationResponseDTO(agent *entity.Agent, pollIntervalSeconds int) dto.RegistrationResponseDTO {
	return dto.RegistrationResponseDTO{
		AgentID:             agent.ID,
		Status:              agent.Status,
		PollURL:             "/config",
		PollIntervalSeconds: pollIntervalSeconds,
	}
}
