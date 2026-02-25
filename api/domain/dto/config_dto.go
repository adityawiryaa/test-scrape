package dto

type ConfigDTO struct {
	ID                  string            `json:"id"`
	Version             int64             `json:"version"`
	Data                map[string]string `json:"data"`
	PollIntervalSeconds int               `json:"poll_interval_seconds"`
}
