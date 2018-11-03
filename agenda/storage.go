package agenda

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/zhanzongyuan/agenda/entity"
)

// Disk Storage
func (agd *Agenda) Load() error {
	// Load and decode User list
	agd.UsernameMap = make(map[string]int)
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
