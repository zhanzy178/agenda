package entity

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// User is a user with name, email information in agenda system
// E.g. TODO
type User struct {
	Id int

	// Base user information for register.
	Name   string
	Email  string
	Number string

	// User account satety, authorization and login state
	password  string
	login     bool
	lastLog   time.Time
	pidCookie int // Here we take bash pid as cookie to maintain every user state
}

// New User
func NewUser(id int, name string, password string, email string, number string) (*User, error) {
	if len(name) < 3 {
		return nil, errors.New("Username number must be longer than 2!")
	}
	if name[0] <= 'z' && name[0] >= 'a' {
		return nil, errors.New("First letter of username must be capitalized!")
	}
	if len(password) < 4 {
		return nil, errors.New("password number must longer than 3!")
	}
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	validNumber := regexp.MustCompile(`^1[34578]\d{9}$`)

	if !validEmail.MatchString(email) {
		return nil, errors.New("Invalid email!")
	}
	if len(number) != 0 && !validNumber.MatchString(number) {
		return nil, errors.New("Invalid number!")
	}

	return &User{Id: id, Name: name, password: password, Email: email, Number: number}, nil
}

// Printing function
func (u User) String() string {
	state := "Online"
	if !u.login {
		state = "Offline"
	}
	str := fmt.Sprintf("[User: %d]\n", u.Id)
	str += fmt.Sprintln("Name:   ", u.Name)
	str += fmt.Sprintln("E-mail: ", u.Email)
	str += fmt.Sprintln("Number: ", u.Number)
	str += fmt.Sprintln("State:  ", state)
	return str
}
