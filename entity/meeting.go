package entity

import "time"

type Meeting struct {
	Id int

	// Base meeting information
	Title           string
	InitiatorId     *int
	ParticipatorsId []int
	StartTime       time.Time
	EndTime         time.Time
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
