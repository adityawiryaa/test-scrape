package request

type UpdateConfigRequest struct {
	Data                map[string]string `json:"data" binding:"required"`
	PollIntervalSeconds int               `json:"poll_interval_seconds"`
}
