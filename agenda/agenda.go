package agenda

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zhanzongyuan/agenda/auth"
	"github.com/zhanzongyuan/agenda/entity"
	"github.com/zhanzongyuan/agenda/validate"
)

var agenda Agenda
var timeLayout string = "2006-01-02 15:04:05 Mon"

type Agenda struct {
	LastUserId    int
	LastMeetingId int

	UserList    []entity.User
	MeetingList []entity.Meeting
	LogList     []entity.Log

	UsernameMap map[string]int

	userDiskFile    string
	meetingDiskFile string
	loginDiskFile   string
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
	agenda.meetingDiskFile = filepath.Join(dataDir, "meeting.json")
	agenda.loginDiskFile = filepath.Join(dataDir, "curUser.txt")

	// Load data
	if err := agenda.Load(); err != nil {
		return err
	}

	return nil
}

// Package Function
func InitConfig(dataDir string) error {
	return agenda.InitConfig(dataDir)
}
func Register(name string, password string, email string, number string) (*entity.User, error) {
	return agenda.Register(name, password, email, number)
}
func CurrentUser() *entity.User {
	return agenda.CurrentUser()
}
func Auth() (*entity.User, error) {
	return agenda.Auth()
}
func Login(name string, password string) (*entity.User, error) {
	return agenda.Login(name, password)
}
func Logout() error {
	return agenda.Logout()
}
func CheckUsers() {
	agenda.CheckUsers()
}
func FindUser(name string) *entity.User {
	return agenda.FindUser(name)
}
func NewMeeting(title string, st time.Time, et time.Time, parsName []string) (*entity.Meeting, error) {
	return agenda.NewMeeting(title, st, et, parsName)
}
func JoinUser(title string, name string) error {
	return agenda.JoinUser(title, name)
}
func MoveoutUser(title string, name string) error {
	return agenda.MoveoutUser(title, name)
}
func CancelMeeting(title string) error {
	return agenda.CancelMeeting(title)
}
func QuitMeeting(title string) error {
	return agenda.QuitMeeting(title)
}
func ClearAllMeetings() error {
	return agenda.ClearAllMeetings()
}
func AllMeetings() ([]entity.Meeting, error) {
	return agenda.AllMeetings()
}
func CheckMeetings(st time.Time, et time.Time) ([]entity.Meeting, error) {
	return agenda.CheckMeetings(st, et)
}
func DeleteUser() error {
	return agenda.DeleteUser()
}
