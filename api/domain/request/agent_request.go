package request

type RegisterAgentRequest struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	Port      int    `json:"port" binding:"required"`
}
