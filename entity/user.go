package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhanzongyuan/agenda/validate"
)

// User is a user with name, email information in agenda system
type User struct {
	Id int

	// Base user information for register.
	Name    string
	Email   string
	Number  string
	LastLog time.Time

	// User account satety, authorization and login state
	password string
	login    bool
	token    int // Here we take bash pid as token to maintain every user state
}

// New User
func NewUser(id int, name string, password string, email string, number string) (*User, error) {
	if err := validate.IsNameValid(name); err != nil {
		return nil, err
	}
	if err := validate.IsPasswordValid(password); err != nil {
		return nil, err
	}

	if err := validate.IsEmailValid(email); err != nil {
		return nil, err
	}
	if err := validate.IsNumberValid(number); len(number) != 0 && err != nil {
		return nil, err
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
	str += fmt.Sprintln("Name:    ", u.Name)
	str += fmt.Sprintln("E-mail:  ", u.Email)
	str += fmt.Sprintln("Number:  ", u.Number)
	str += fmt.Sprintln("State:   ", state)
	str += fmt.Sprintln("LastLog: ", u.LastLog)
	return str
}

// Auth
func (u *User) Auth(name string, password string) bool {
	return name == u.Name && password == u.password
}

// Check Pid
func (u *User) CheckToken(token int) bool {
	return token == u.token
}

// Set Pid
func (u *User) UpdateToken(token int) {
	u.token = token
}

// login flag
func (u *User) IsLogin() bool {
	return u.login
}
func (u *User) Login() {
	u.login = true
	u.LastLog = time.Now()
}
func (u *User) Logout() {
	u.login = false
}

// JSON
func (u *User) UnmarshalJSON(b []byte) error {
	user := struct {
		Id       int
		Name     string
		Email    string
		Number   string
		LastLog  time.Time
		Password string
		Login    bool
		Token    int //
	}{}
	err := json.Unmarshal(b, &user)
	if err != nil {
		return err
	}
	u.Id = user.Id
	u.Name = user.Name
	u.Email = user.Email
	u.Number = user.Number
	u.LastLog = user.LastLog
	u.password = user.Password
	u.login = user.Login
	u.token = user.Token
	return nil
}

func (u User) MarshalJSON() ([]byte, error) {
	user := struct {
		Id       int
		Name     string
		Email    string
		Number   string
		LastLog  time.Time
		Password string
		Login    bool
		Token    int
	}{}
	user.Id = u.Id
	user.Name = u.Name
	user.Email = u.Email
	user.Number = u.Number
	user.LastLog = u.LastLog
	user.Password = u.password
	user.Login = u.login
	user.Token = u.token

	return json.Marshal(user)
}
