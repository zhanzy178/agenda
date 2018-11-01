package entity

import "time"

// Log to record user login information
type Log struct {
	UserId      int
	Token       int
	LastLogDate time.Time
}
