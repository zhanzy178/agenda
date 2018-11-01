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

type Agenda struct {
	LastId int

	UserList    []entity.User
	MeetingList []entity.Meeting
	LogList     []entity.Log

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
	agenda.meetingDiskFile = filepath.Join(dataDir, "meeting,json")
	agenda.loginDiskFile = filepath.Join(dataDir, "curUser.txt")

	// Load data
	if err := agenda.Load(); err != nil {
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
	if _, err := os.Lstat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Clear list
		switch opt {
		case "User":
			agd.UserList = agd.UserList[:0]
		case "Meeting":
			agd.MeetingList = agd.MeetingList[:0]
		case "Log":
			agd.LogList = agd.LogList[:0]
		}

		// Decoding line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jsonBlob := scanner.Text()

			switch opt {
			case "User":
				agd.UserList = append(agd.UserList, entity.User{})
				json.Unmarshal([]byte(jsonBlob), &agd.UserList[len(agd.UserList)-1])
				tId := agd.UserList[len(agd.UserList)-1].Id
				if tId > agd.LastId {
					agd.LastId = tId
				}
			case "Meeting":
				agd.MeetingList = append(agd.MeetingList, entity.Meeting{})
				json.Unmarshal([]byte(jsonBlob), &agd.MeetingList[len(agd.MeetingList)-1])
			case "Log":
				agd.LogList = append(agd.LogList, entity.Log{})
				json.Unmarshal([]byte(jsonBlob), &agd.LogList[len(agd.LogList)-1])
			}
		}
		log.Printf("%s list loaded.", opt)

	}

	return nil
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
			if _, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	case "Meeting":
		for _, item := range agd.MeetingList {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if _, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	case "Log":
		for _, item := range agd.LogList {
			b, err := json.Marshal(item)
			if err != nil {
				return err
			}
			if _, err := file.WriteString(string(b) + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

// User Management
func (agd *Agenda) Register(name string, password string, email string, number string) (*entity.User, error) {
	for _, u := range agd.UserList {
		if u.Name == name {
			return nil, errors.New(fmt.Sprintf("'%s' is exist!\n", name))
		}
	}
	user, err := entity.NewUser(0, name, password, email, number)
	if err != nil {
		return nil, err
	}
	agd.LastId++
	user.Id = agd.LastId
	agd.UserList = append(agd.UserList, *user)
	user = &agd.UserList[len(agd.UserList)-1]
	if err := agd.Sync("User"); err != nil {
		return user, err
	}
	return user, nil
}
func (agd *Agenda) Login(name string, password string) error {
	return nil
}
func (agd *Agenda) Logout(name string) error {
	return nil
}
func (agd *Agenda) CheckUsers(name_list []string) {
}
func (agd *Agenda) FindUser(name string) *entity.User {
	return nil
}
func (agd *Agenda) RemoveUser(name string) error {
	return nil
}

// Meeting Management
func (agd *Agenda) NewMeeting(title string, st time.Time, et time.Time, initiator *entity.User) (*entity.Meeting, error) {
	return nil, nil
}
func (agd *Agenda) FindMeeting(title string) (*entity.Meeting, error) {
	return nil, nil
}

// Package Function
func InitConfig(dataDir string) error {
	return agenda.InitConfig(dataDir)
}
func Register(name string, password string, email string, number string) (*entity.User, error) {
	return agenda.Register(name, password, email, number)
}
func Login(name string, password string) error {
	return agenda.Login(name, password)
}
func Logout(name string) error {
	return agenda.Logout(name)
}
func CheckUsers(name_list []string) {
	agenda.CheckUsers(name_list)
}
func FindUser(name string) *entity.User {
	return agenda.FindUser(name)
}
func RemoveUser(name string) error {
	return agenda.RemoveUser(name)

}
func NewMeeting(title string, st time.Time, et time.Time, initiator *entity.User) (*entity.Meeting, error) {
	return agenda.NewMeeting(title, st, et, initiator)
}
func FindMeeting(title string) (*entity.Meeting, error) {
	return agenda.FindMeeting(title)
}
