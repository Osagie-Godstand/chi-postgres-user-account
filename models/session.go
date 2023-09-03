package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
