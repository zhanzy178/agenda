package entity

import "time"

type Meeting struct {
	Id int

	// Base meeting information
	Title         string
	Initiator     User
	Participators []User
	StartTime     time.Time
	EndTime       time.Time
}

// Add or Remove participator
func (m *Meeting) Join(u *User) error {
	return nil
}
func (m *Meeting) Quit(u *User) error {
	return nil
}

func (m *Meeting) authorize(ini *User) error {
	return nil
}
