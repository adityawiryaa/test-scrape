package entity

type RegistrationRequest struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	Port      int    `json:"port" binding:"required"`
}

type RegistrationResponse struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}
