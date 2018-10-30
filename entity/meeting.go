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

}
func (m *Meeting) Quit(u *User) error {

}

func (m *Meeting) authorize(ini *User) error {

}
