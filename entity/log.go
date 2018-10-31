package entity

import "time"

// Log to record user login information
type Log struct {
	UserId      int
	Pid         int
	LastLogDate time.Time
}
