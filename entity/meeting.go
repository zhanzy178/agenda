package entity

import (
	"fmt"
	"time"
)

var timeLayout string = "2006-01-02 15:04:05 Mon"

type Meeting struct {
	Id int

	// Base meeting information
	Title             string
	InitiatorName     string
	ParticipatorsName []string
	StartTime         time.Time
	EndTime           time.Time
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

func (m Meeting) Conflict(me Meeting) bool {
	return m.StartTime.Before(me.EndTime) && me.StartTime.Before(m.EndTime)
}

func (m *Meeting) String() string {
	indent := "  "
	str := fmt.Sprintf("******[Meeting: %d]******\n", m.Id)
	str += fmt.Sprintln("Title: ", m.Title)
	str += fmt.Sprintln("Initiator: ", m.InitiatorName)
	str += fmt.Sprintln("Time: ")
	str += fmt.Sprintf("%sSince (%s)\n%sTo (%s)\n",
		indent, m.StartTime.Format(timeLayout), indent, m.EndTime.Format(timeLayout))

	str += fmt.Sprintln("Participators: ")
	for i, s := range m.ParticipatorsName {
		if i%3 == 0 {
			if i != 0 {
				str += fmt.Sprint("\n")
			}
			str += fmt.Sprint(indent)
		}
		str += fmt.Sprint(s, ", ")
	}
	return str
}
