package agenda

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/zhanzongyuan/agenda/entity"
	"github.com/zhanzongyuan/agenda/validate"
)

// Meeting Management
func (agd *Agenda) NewMeeting(title string, st time.Time, et time.Time, parsName []string) (*entity.Meeting, error) {
	user, err := agd.Auth()
	if err != nil {
		return nil, err
	}
	meeting := &entity.Meeting{
		Title:             title,
		SponsorName:       user.Name,
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

func (agd *Agenda) CheckParticipatorsNameList(nameList []string) ([]string, error) {
	if len(nameList) == 0 {
		return nil, errors.New("Empty participators list")
	}
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
			nameList = append(nameList[:i], nameList[i+1:]...)
		}
	}
	return nameList, nil
}
func (agd *Agenda) IsMeetingValid(m *entity.Meeting) error {
	// Check username
	pm, err := agd.CheckParticipatorsNameList(m.ParticipatorsName)
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
func (agd *Agenda) JoinUser(title, name string) error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	// Check Name
	if agd.UsernameMap[name] == 0 {
		return errors.New(fmt.Sprintf("User '%s' is not exist", name))
	}

	// Get user meeting that it participate in
	ml, _ := agd.UserMeeting(name, false)

	// Check Meeting
	for i, m := range agd.MeetingList {
		if m.Title == title {
			// No access
			if curUser.Name != m.SponsorName {
				return errors.New(fmt.Sprintf("You have not access to the meeting '%s'", title))
			}

			// if user has in meeting or conflict with other meeting
			for _, meeting := range ml {
				if meeting.Title == title {
					return errors.New("User have been in this meeting!")
				} else if meeting.Conflict(m) {
					return errors.New(fmt.Sprintf("Conflict with meeting '%s'!", meeting.Title))
				}
			}

			// Add user
			agd.MeetingList[i].ParticipatorsName = append(agd.MeetingList[i].ParticipatorsName, name)
			if err := agd.Sync("Meeting"); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Meeting '%s' is not exist", title))
}

func (agd *Agenda) MoveoutUser(title, name string) error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	// Check Name
	if agd.UsernameMap[name] == 0 {
		return errors.New(fmt.Sprintf("User '%s' is not exist", curUser.Name))
	}

	for i := range agd.MeetingList {
		meeting := &agd.MeetingList[i]
		if meeting.Title == title {
			// No access
			if curUser.Name != meeting.SponsorName {
				return errors.New(fmt.Sprintf("You have not access to the meeting '%s'", title))
			}

			// Check name in participater
			index, join := -1, false
			for j, p := range meeting.ParticipatorsName {
				if p == name {
					index = j
					join = true
					break
				}
			}
			if !join {
				return errors.New("User is not in this meeting")
			} else {
				meeting.ParticipatorsName = append(meeting.ParticipatorsName[:index], meeting.ParticipatorsName[index+1:]...)
				if len(meeting.ParticipatorsName) == 0 {
					agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
				}
				if err := agd.Sync("Meeting"); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}
func (agd *Agenda) CancelMeeting(title string) error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	for i := range agd.MeetingList {
		meeting := &agd.MeetingList[i]
		if meeting.Title == title {
			// No access
			if curUser.Name != meeting.SponsorName {
				return errors.New(fmt.Sprintf("You have not access to the meeting '%s'", title))
			}

			// Cancel
			agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
			if err := agd.Sync("Meeting"); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Meeting '%s' exist!", title))

}
func (agd *Agenda) QuitMeeting(title string) error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	for i := range agd.MeetingList {
		meeting := &agd.MeetingList[i]
		if meeting.Title == title {
			index, pa := -1, false
			for j, p := range meeting.ParticipatorsName {
				if p == curUser.Name {
					index, pa = j, true
					break
				}
			}
			if !pa {
				return errors.New(fmt.Sprintf("Current user '%s' is not in this meeting '%s'", curUser.Name, title))
			} else {
				meeting.ParticipatorsName = append(meeting.ParticipatorsName[:index], meeting.ParticipatorsName[index+1:]...)
				if len(meeting.ParticipatorsName) == 0 {
					agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
				}
				if err := agd.Sync("Meeting"); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf("Meeting '%s' is not exist", title))
}
func (agd *Agenda) ClearAllMeetings() error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	i := 0
	for i < len(agd.MeetingList) {
		if agd.MeetingList[i].SponsorName == curUser.Name {
			agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
		} else {
			i++
		}
	}
	if err := agd.Sync("Meeting"); err != nil {
		return err
	}
	return nil
}
func (agd *Agenda) AllMeetings() ([]entity.Meeting, error) {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return nil, err
	}

	ml, _ := agd.UserMeeting(curUser.Name, true)

	fmt.Println(agd.TableFormatMeetingList(ml))
	return agd.MeetingList, nil
}
func (agd *Agenda) CheckMeetings(st time.Time, et time.Time) ([]entity.Meeting, error) {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return nil, err
	}
	ml, _ := agd.UserMeeting(curUser.Name, true)
	filtedMeeting := []entity.Meeting{}
	for i := range ml {
		meeting := &ml[i]
		if (st.Before(meeting.StartTime) && et.After(meeting.StartTime)) ||
			(st.Before(meeting.EndTime) && et.After(meeting.EndTime)) {
			filtedMeeting = append(filtedMeeting, *meeting)
		}
	}
	fmt.Println(agd.TableFormatMeetingList(filtedMeeting))
	return filtedMeeting, nil
}
func (agd *Agenda) UserMeeting(username string, sponor bool) ([]entity.Meeting, error) {
	ml := []entity.Meeting{}
	for _, m := range agd.MeetingList {
		if sponor && m.SponsorName == username {
			ml = append(ml, m)
			continue
		}
		for _, n := range m.ParticipatorsName {
			if n == username {
				ml = append(ml, m)
				break
			}
		}
	}
	return ml, nil
}
func (agd *Agenda) TableFormatMeetingList(ml []entity.Meeting) string {
	str := ""
	idW, titleW, sprW, parsW, startW, endW :=
		len("Id"), len("Title"), len("Sponsor"), len("Participators"), len("Since"), len("To")
	for i := range ml {
		meeting := &ml[i]
		idL, titleL, sprL, startL, endL :=
			len(strconv.FormatInt(int64(meeting.Id), 10)), len(meeting.Title), len(meeting.SponsorName),
			len(meeting.StartTime.Format(timeLayout)), len(meeting.EndTime.Format(timeLayout))
		if idW < idL {
			idW = idL
		}
		if titleW < titleL {
			titleW = titleL
		}
		if sprW < sprL {
			sprW = sprL
		}
		if startW < startL {
			startW = startL
		}
		if endW < endL {
			endW = endL
		}
		for _, p := range meeting.ParticipatorsName {
			if len(p)+2 > parsW {
				parsW = len(p) + 2
			}
		}
	}
	idW += 2
	titleW += 2
	sprW += 2
	startW += 2
	endW += 2
	parsW += 2

	outputFormat := ""
	outputFormat += "%-" + strconv.FormatInt(int64(idW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(titleW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(sprW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(startW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(endW), 10) + "s |"
	outputFormat += "%-" + strconv.FormatInt(int64(parsW), 10) + "s\n"

	totalW := idW + titleW + sprW + startW + endW + parsW + 21
	title := " Meeting Table "
	padW := (totalW - len(title)) / 2
	if padW < 0 {
		padW = 0
	}
	str += fmt.Sprint(strings.Repeat("-", padW))
	str += fmt.Sprint(title)
	str += fmt.Sprintln(strings.Repeat("-", padW))
	str += fmt.Sprintf(outputFormat, "Id", "Title", "Sponsor", "Since", "To", "Participators")
	str += fmt.Sprintln(strings.Repeat("-", totalW))
	for _, meeting := range ml {
		for i, p := range meeting.ParticipatorsName {
			if i == 0 {
				str += fmt.Sprintf(outputFormat,
					strconv.FormatInt(int64(meeting.Id), 10),
					meeting.Title, meeting.SponsorName,
					meeting.StartTime.Format(timeLayout),
					meeting.EndTime.Format(timeLayout),
					p+", ",
				)
			} else {
				str += fmt.Sprintf(outputFormat, "", "", "", "", "", p+", ")
			}
		}
		str += fmt.Sprintln(strings.Repeat("-", totalW))
	}
	return str
}
func (agd *Agenda) TableFormatUserList(ul []entity.User) string {
	str := ""
	idW, nameW, emailW, numberW, lastLogW, loginW :=
		len("Id"), len("Name"), len("Email"), len("Number"), len("Last-Log"), len("Offline")
	for i := range ul {
		user := &ul[i]
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

	totalW := idW + nameW + emailW + numberW + loginW + lastLogW + 21
	title := " User Table "
	padW := (totalW - len(title)) / 2
	if padW < 0 {
		padW = 0
	}
	str += fmt.Sprint(strings.Repeat("-", padW))
	str += fmt.Sprint(title)
	str += fmt.Sprintln(strings.Repeat("-", padW))
	str += fmt.Sprintf(outputFormat, "Id", "Name", "E-mail", "Number", "State", "Last-Log")
	str += fmt.Sprintln(strings.Repeat("-", totalW))
	for _, u := range ul {
		state := "Offline"
		if u.IsLogin() {
			state = "Online"
		}
		str += fmt.Sprintf(outputFormat, strconv.FormatInt(int64(u.Id), 10),
			u.Name, u.Email, u.Number, state, u.LastLog.Format(timeLayout))
		str += fmt.Sprintln(strings.Repeat("-", totalW))
	}
	return str
}
