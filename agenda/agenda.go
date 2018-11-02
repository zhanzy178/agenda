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
			agd.UsernameMap = make(map[string]int)
		case "Meeting":
			agd.MeetingList = agd.MeetingList[:0]
		case "Log":
			agd.LogList = agd.LogList[:0]
		}

		// Decoding line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jsonBlob := scanner.Text()
			if len(jsonBlob) == 0 {
				continue
			}
			switch opt {
			case "User":
				agd.UserList = append(agd.UserList, entity.User{})
				json.Unmarshal([]byte(jsonBlob), &agd.UserList[len(agd.UserList)-1])
				tId, name := agd.UserList[len(agd.UserList)-1].Id, agd.UserList[len(agd.UserList)-1].Name
				agd.UsernameMap[name] = 1
				if tId > agd.LastUserId {
					agd.LastUserId = tId
				}
			case "Meeting":
				agd.MeetingList = append(agd.MeetingList, entity.Meeting{})
				json.Unmarshal([]byte(jsonBlob), &agd.MeetingList[len(agd.MeetingList)-1])
				mId := agd.MeetingList[len(agd.MeetingList)-1].Id
				if mId > agd.LastMeetingId {
					agd.LastMeetingId = mId
				}
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
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
	agd.LastUserId++
	user.Id = agd.LastUserId
	agd.UserList = append(agd.UserList, *user)
	user = &agd.UserList[len(agd.UserList)-1]
	agd.UsernameMap[user.Name] = 1
	if err := agd.Sync("User"); err != nil {
		return user, err
	}
	return user, nil
}

func (agd *Agenda) CurrentUser() *entity.User {
	// Check password and pid
	curPid := auth.CurrentBashPid()

	// Check current bash state
	for i := range agd.UserList {
		user := &agd.UserList[i]
		if user.IsLogin() && user.CheckToken(curPid) {
			return user
		}
	}
	return nil
}

func (agd *Agenda) Login(name string, password string) (*entity.User, error) {
	// Check password and pid
	curPid := auth.CurrentBashPid()

	// Login auth check
	authLogin := false
	var user *entity.User
	for i := range agd.UserList {
		user = &agd.UserList[i]
		if user.Auth(name, password) {
			authLogin = true
			// Other login this user.
			// Warning: other bash login this user may lost authorization, remove other Log from list
			l := 0
			for l < len(agd.LogList) {
				if agd.LogList[l].UserId == user.Id {
					if l != len(agd.LogList)-1 {
						agd.LogList = append(agd.LogList[:l], agd.LogList[l+1:]...)
					} else {
						agd.LogList = agd.LogList[:l]
					}
					log.Println("Warning: Other bash login this user may lost authorization!")
				} else {
					l++
				}
			}

			// Login and record current bash pid
			user.Login()
			user.UpdateToken(curPid)
			agd.LogList = append(agd.LogList, entity.Log{UserId: user.Id, Token: curPid, LastLogDate: user.LastLog})
			if err := agd.Sync("User"); err != nil {
				return nil, err
			}
			if err := agd.Sync("Log"); err != nil {
				return nil, err
			}
			break
		}
	}
	if authLogin {
		return user, nil
	} else {
		return nil, errors.New("Invalid password or username")
	}
}
func (agd *Agenda) Logout() error {
	user := agd.CurrentUser()
	if user == nil {
		log.Println("There is not user login in current bash!")
		return nil
	}
	user.Logout()
	l := 0
	curPid := auth.CurrentBashPid()
	for l < len(agd.LogList) {
		if agd.LogList[l].UserId == user.Id || agd.LogList[l].Token == curPid {
			if l != len(agd.LogList)-1 {
				agd.LogList = append(agd.LogList[:l], agd.LogList[l+1:]...)
			} else {
				agd.LogList = agd.LogList[:l]
			}
		} else {
			l++
		}
	}
	if err := agd.Sync("User"); err != nil {
		return err
	}
	if err := agd.Sync("Log"); err != nil {
		return err
	}
	log.Printf("Logout user '%s' successfully!\n", user.Name)
	fmt.Println(user)
	return nil
}
func (agd *Agenda) Auth() (*entity.User, error) {
	user := agd.CurrentUser()
	if user == nil {
		return nil, errors.New("You are not login!")
	} else {
		log.Printf("Current user '%s'\n", user.Name)
		return user, nil
	}
}
func (agd *Agenda) CheckUsers() {
	if _, err := agd.Auth(); err != nil {
		log.Fatal(err)
		return
	}
	idW, nameW, emailW, numberW, lastLogW, loginW :=
		len("Id"), len("Name"), len("Email"), len("Number"), len("Last-Log"), len("Offline")
	for i := range agd.UserList {
		user := &agd.UserList[i]
		idL, nameL, emailL, numberL, lastLogL :=
			len(strconv.FormatInt(int64(user.Id), 10)), len(user.Name), len(user.Email),
			len(user.Number), len(user.LastLog.Format(timeLayout))
		if idW < idL {
			idW = idL
		}
		if nameW < nameL {
			nameW = nameL
		}
		if emailW < emailL {
			emailW = emailL
		}
		if numberW < numberL {
			numberW = numberL
		}
		if lastLogW < lastLogL {
			lastLogW = lastLogL
		}
	}
	idW += 2
	nameW += 2
	emailW += 2
	numberW += 2
	loginW += 2
	lastLogW += 2
	outputFormat := ""
	outputFormat += "%-" + strconv.FormatInt(int64(idW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(nameW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(emailW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(numberW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(loginW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(lastLogW), 10) + "s\n"

	fmt.Printf(outputFormat, "Id", "Name", "E-mail", "Number", "State", "Last-Log")
	fmt.Println(strings.Repeat("-", idW+nameW+emailW+numberW+loginW+lastLogW+21))
	for _, u := range agd.UserList {
		state := "Offline"
		if u.IsLogin() {
			state = "Online"
		}
		fmt.Printf(outputFormat, strconv.FormatInt(int64(u.Id), 10),
			u.Name, u.Email, u.Number, state, u.LastLog.Format(timeLayout))
	}

}
func (agd *Agenda) FindUser(name string) *entity.User {
	if _, err := agd.Auth(); err != nil {
		log.Fatal(err)
		return nil
	}
	for i := range agd.UserList {
		user := &agd.UserList[i]
		if user.Name == name {
			return user
		}
	}

	return nil
}
func (agd *Agenda) RemoveUser(name string) error {
	return nil
}

// Meeting Management
func (agd *Agenda) NewMeeting(title string, st time.Time, et time.Time, parsName []string) (*entity.Meeting, error) {
	user, err := agd.Auth()
	if err != nil {
		return nil, err
	}
	meeting := &entity.Meeting{
		Title:             title,
		InitiatorName:     user.Name,
		ParticipatorsName: parsName,
		StartTime:         st,
		EndTime:           et,
	}
	if err := agd.IsMeetingValid(meeting); err != nil {
		return nil, err
	}
	agd.MeetingList = append(agd.MeetingList, *meeting)
	meeting = &agd.MeetingList[len(agd.MeetingList)-1]
	agd.LastMeetingId++
	meeting.Id = agd.LastMeetingId

	if err := agd.Sync("Meeting"); err != nil {
		return nil, err
	}
	return meeting, nil
}
func (agd *Agenda) FindMeeting(title string) (*entity.Meeting, error) {
	return nil, nil
}
func (agd *Agenda) CheckNameList(nameList []string) ([]string, error) {
	duplicateMap := make(map[string]int)
	i := 0
	for i < len(nameList) {
		n := nameList[i]
		if agd.UsernameMap[n] == 0 {
			return nil, errors.New(fmt.Sprintf("User '%s' is not exist", n))
		}
		if duplicateMap[n] == 0 {
			duplicateMap[n] = 1
			i++
		} else if duplicateMap[n] == 1 {
			log.Printf("Username '%s' duplicate!", n)
			if i == len(nameList)-1 {
				nameList = nameList[:i]
			} else {
				nameList = append(nameList[:i], nameList[i+1:]...)
			}
		}
	}
	return nameList, nil
}
func (agd *Agenda) IsMeetingValid(m *entity.Meeting) error {
	// Check username
	pm, err := agd.CheckNameList(m.ParticipatorsName)
	if err != nil {
		return err
	}
	m.ParticipatorsName = pm

	// Check Time
	if err := validate.IsStartEndTimeValid(m.StartTime, m.EndTime); err != nil {
		return err
	}

	// Check Title
	if err := validate.IsTitleValid(m.Title); err != nil {
		return err
	}

	// Check meeting conflict
	for i := range agd.MeetingList {
		meeting := &agd.MeetingList[i]
		if m.Title == meeting.Title {
			return errors.New(fmt.Sprintf("Duplicate meeting title '%s'!", m.Title))
		}
		if meeting.Conflict(*m) {
			for _, p := range meeting.ParticipatorsName {
				for _, mp := range m.ParticipatorsName {
					if p == mp {
						return errors.New(fmt.Sprintf("Conflict with meeting '%s'!", meeting.Title))
					}
				}
			}
		}
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
func RemoveUser(name string) error {
	return agenda.RemoveUser(name)
}
func NewMeeting(title string, st time.Time, et time.Time, parsName []string) (*entity.Meeting, error) {
	return agenda.NewMeeting(title, st, et, parsName)
}
func FindMeeting(title string) (*entity.Meeting, error) {
	return agenda.FindMeeting(title)
}
