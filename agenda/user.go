package agenda

import (
	"errors"
	"fmt"
	"log"

	"github.com/zhanzongyuan/agenda/auth"
	"github.com/zhanzongyuan/agenda/entity"
)

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

func (agd *Agenda) Login(name, password string) (*entity.User, error) {
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
					agd.LogList = append(agd.LogList[:l], agd.LogList[l+1:]...)
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
	user, err := agd.Auth()
	if err != nil {
		return err
	}
	user.Logout()
	l := 0
	curPid := auth.CurrentBashPid()
	for l < len(agd.LogList) {
		if agd.LogList[l].UserId == user.Id || agd.LogList[l].Token == curPid {
			agd.LogList = append(agd.LogList[:l], agd.LogList[l+1:]...)
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
	fmt.Println(agd.TableFormatUserList(agd.UserList))
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
func (agd *Agenda) DeleteUser() error {
	// Check user online
	curUser, err := agd.Auth()
	if err != nil {
		return err
	}

	// Delete Log state
	i := 0
	for i < len(agd.LogList) {
		if agd.LogList[i].UserId == curUser.Id {
			agd.LogList = append(agd.LogList[:i], agd.LogList[i+1:]...)
		} else {
			i++
		}
	}

	// Delete Meeting
	i = 0
	for i < len(agd.MeetingList) {
		if agd.MeetingList[i].SponsorName == curUser.Name {
			agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
		} else {
			pars := &agd.MeetingList[i].ParticipatorsName
			for i, p := range *pars {
				if p == curUser.Name {
					*pars = append((*pars)[:i], (*pars)[i+1:]...)
					break
				}
			}
			if len(*pars) == 0 {
				agd.MeetingList = append(agd.MeetingList[:i], agd.MeetingList[i+1:]...)
			} else {
				i++
			}
		}
	}

	// Delete user
	for j := range agd.UserList {
		if agd.UserList[j].Name == curUser.Name {
			agd.UserList = append(agd.UserList[:j], agd.UserList[j+1:]...)
			break
		}
	}

	if err = agd.Sync("User"); err != nil {
		return err
	}
	if err = agd.Sync("Meeting"); err != nil {
		return err
	}
	if err = agd.Sync("Log"); err != nil {
		return err
	}
	return nil
}
