package dto

type AgentDTO struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
}

type RegistrationResponseDTO struct {
	AgentID             string `json:"agent_id"`
	Status              string `json:"status"`
	PollURL             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
}
