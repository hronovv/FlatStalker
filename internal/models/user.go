package models

import "time"

type User struct {
	ID            int64
	ChatID        int64
	Plan          string
	PlanExpiresAt *time.Time
}
