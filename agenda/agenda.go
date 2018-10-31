package agenda

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zhanzongyuan/agenda/entity"
)

var agenda Agenda

type User entity.User
type Meeting entity.Meeting
type Agenda struct {
	UserList    []User
	MeetingList []Meeting

	userDiskFile    os.File
	meetingDiskFile os.File
	loginDiskFile   os.File
}

// Get System Agenda
func SystemAgenda() *Agenda {
	return &agenda
}

// Config Agenda disk data directory
func (agd *Agenda) InitConfig(dataDir string) error {
	// Check directory
	fi, err := os.Lstat(dataDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	if !fi.Mode().IsDir() {
		errMsg := fmt.Sprintf("'%s' is not a directory.", dataDir)
		fmt.Fprintln(os.Stderr, errMsg)
		return errors.New(errMsg)
	}

	// Config user, meeting, login disk file
	agenda.userDiskFile = filepath.Join(dataDir, "user.json")
	agenda.meetingDistFile = filepath.Join(dataDir, "meeting,json")
	agenda.loginDiskFile = filepath.Join(dataDir, "curUser.txt")

	// Load data
	if err := agenda.load(); err != nil {
		return err
	}

	return nil
}

// Disk Storage
func (agd *Agenda) load() error {
	userPath, meetingPath, loginPath := agenda.userDiskFile, agenda.meetingDistFile, agenda.loginDiskFile

	return nil

}
func (agd *Agenda) sync() error {
	return nil

}

// User Management
func (agd *Agenda) Register(name string, password string, email string) (*User, error) {
	return nil, nil

}
func (agd *Agenda) Login(name string, password string) error {
	return nil

}
func (agd *Agenda) Logout(name string) error {
	return nil

}
func (agd *Agenda) CheckUsers(name_list []string) {

}
func (agd *Agenda) FindUser(name string) *User {
	return nil

}
func (agd *Agenda) RemoveUser(name string) error {
	return nil

}

// Meeting Management
func (agd *Agenda) NewMeeting(title string, st time.Time, et time.Time, initiator *User) (*Meeting, error) {
	return nil, nil

}
func (agd *Agenda) FindMeeting(title string) (*Meeting, error) {
	return nil, nil

}

// Package Function
func InitConfig(dataDir string) error {
	return agenda.InitConfig(dataDir)
}
func Register(name string, password string, email string) (*User, error) {
	return nil, nil

}
func Login(name string, password string) error {
	return nil

}
func Logout(name string) error {
	return nil

}
func CheckUsers(name_list []string) {

}
func FindUser(name string) *User {
	return nil

}
func RemoveUser(name string) error {
	return nil

}
func NewMeeting(title string, st time.Time, et time.Time, initiator *User) (*Meeting, error) {
	return nil, nil

}
func FindMeeting(title string) (*Meeting, error) {
	return nil, nil

}
