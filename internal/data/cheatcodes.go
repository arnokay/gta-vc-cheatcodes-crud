package data

import (
	"time"
)

type Movie struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Code        *string   `json:"code"`
	Description *string   `json:"description"`
	Tags        []string  `json:"tags"`
	Version     int32     `json:"version"`
}
