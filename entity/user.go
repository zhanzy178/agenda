package entity

import (
	"fmt"
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
	password      string
	signin        bool
	lastSigninLog time.Time
	pidCookie     int // Here we take bash pid as cookie to maintain every user state
}

// Printing function
func (u *User) PrintLong() {
	fmt.Printf("[User: %d]\n", u.Id)
	fmt.Println("Name:   ", u.Name)
	fmt.Println("E-mail: ", u.Email)
	fmt.Println("Number: ", u.Number)
}
func (u *User) PrintShort() {
	fmt.Printf("id=%d, name=%s, e-mail=%s, number=%s\n",
		u.Id, u.Name, u.Email, u.Number)
}
