package agenda

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zhanzongyuan/agenda/entity"
)

var agenda Agenda

type User entity.User
type Meeting entity.Meeting
type Log entity.Log
type Agenda struct {
	UserList    []User
	MeetingList []Meeting
	LogList     []Log

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
func (agd *Agenda) Load() error {
	// Load and decode User list
	if err := agd.loadList("User"); err != nil {
		return err
	}

	// Load and decode Meeting list from disk
	if err := agd.loadList("Meeting"); err != nil {
		return err
	}

	// Load and decode OnlineLog list from disk
	if err := agd.loadList("Log"); err != nil {
		return err
	}

	return nil
}
func (agd *Agenda) loadList(opt string) error {
	var filePath string
	switch opt {
	case "User":
		filePath = agd.userDiskFile
	case "Meeting":
		filePath = agd.meetingDiskFile
	case "Log":
		filePath = agd.loginDiskFile
	default:
		return errors.New(fmt.Sprintf("loadList: invalid list opt '%s'", opt))
	}
	// Load and decode list from disk
	if fi, err := os.Lstat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Decoding line by line
		scanner := bufio.Scanner(file)
		for scanner.Scan() {
			jsonBlob := scanner.Text()

			switch opt {
			case "User":
				agd.UserList = append(agd.UserList, User{})
				json.Unmarshal(jsonBlob, &agd.UserList[len(agd.UserList)-1])
			case "Meeting":
				agd.MeetingList = append(agd.MeetingList, Meeting{})
				json.Unmarshal(jsonBlob, &agd.MeetingList[len(agd.MeetingList)-1])
			case "Log":
				agd.LogList = append(agd.LogList, Log{})
				json.Unmarshal(jsonBlob, &agd.LogList[len(agd.LogList)-1])
			}
		}
		log.Printf("%s list loaded.", opt)

		return nil
	}

}
func (agd *Agenda) Sync(opt string) error {
	var filePath string
	switch opt {
	case "User":
		filePath = agd.userDiskFile
	case "Meeting":
		filePath = agd.meetingDiskFile
	case "Log":
		filePath = agd.loginDiskFile
	default:
		return errors.New(fmt.Sprintf("Sync: invalid list opt '%s'", opt))
	}

	// Readinfile
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write to File
	switch opt {
	case "User":
		for _, item := range agd.UserList {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if n, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	case "Meeting":
		for _, item := range agd.MeetingList {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if n, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	case "Log":
		for _, item := range agd.LogList {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if n, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	}

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
	return agenda.Register(name, password, email)
}
func Login(name string, password string) error {
	return agenda.Login(name, password)
}
func Logout(name string) error {
	return agenda.Logout(name)
}
func CheckUsers(name_list []string) {
	return agenda.CheckUsers(name_list)
}
func FindUser(name string) *User {
	return agenda.FindUser(name)
}
func RemoveUser(name string) error {
	return agenda.RemoveUser(name)

}
func NewMeeting(title string, st time.Time, et time.Time, initiator *User) (*Meeting, error) {
	return agenda.NewMeeting(title, st, et, initiator)
}
func FindMeeting(title string) (*Meeting, error) {
	return agenda.FindMeeting(title)
}
