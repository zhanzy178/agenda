package agenda

import (
	"os"
	"time"

	"github.com/zhanzongyuan/agenda/entity"
)

type User entity.User
type Meeting entity.Meeting
type Agenda struct {
	UserList    []User
	MeetingList []Meeting

	userDiskFile    os.File
	meetingDiskFile os.File
	loginDiskFile   os.File
}

// Disk Storage

// User Management
func Register(name string, password string, email string) (*User, error) {

}
func Login(name string, password string) error {

}
func Logout(name string) error {

}
func CheckUsers(name_list []string) {

}
func FindUser(name string) *User {

}
func RemoveUser(name string) error {

}

// Meeting Management
func NewMeeting(title string, st time.Time, et time.Time, initiator *User) (*Meeting, error) {

}
func FindMeeting(title string) (*Meeting, error) {

}
