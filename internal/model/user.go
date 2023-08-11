package model

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID
	Quota
}

type Quota struct {
	MonthSizeLimit  float64 `json:"month_size_limit"`
	MinuteRateLimit int     `json:"minute_rate_limit"`
}
