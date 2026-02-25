package entity

import "time"

type Config struct {
	ID                  string            `json:"id"`
	Version             int64             `json:"version"`
	Data                map[string]string `json:"data"`
	PollIntervalSeconds int               `json:"poll_interval_seconds"`
	CreatedAt           time.Time         `json:"created_at"`
}
