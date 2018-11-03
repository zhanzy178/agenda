package agenda

import (
	"testing"
	"time"
	"fmt"
)

// test register
func TestRegister(t *testing.T) {
    if user, err := Register("caiye","1234","caiye@foxmail.com","0"); err != nil {
		t.Error("Register test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("Register test success");
        t.Log("Name: " + user.Name);
        t.Log("Email: " + user.Email);
        t.Log("Number: " + user.Number);
    }
}


// test login
func TestLogin(t *testing.T) {
    if user, err := Login("Caiye","1234"); err != nil {
		t.Error("Login test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("Login test success");
        t.Log("Name: " + user.Name);
        t.Log("Email: " + user.Email);
        t.Log("Number: " + user.Number);
    }
}


// test logout
func TestLogout(t *testing.T) {
    if err := Logout(); err != nil {
		t.Error("Logout test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("Logout test success");
    }
}


// test delete user
func TestDeleteUser(t *testing.T) {
    if err := DeleteUser(); err != nil {
		t.Error("DeleteUser test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("DeleteUser test success");
    }
}


// test create meeting
func TestNewMeeting(t *testing.T) {
    if _, err := NewMeeting("Meeting", time.Now(), time.Now(), []string{"Yvonne"}); err != nil {
		t.Error("NewMeeting test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("NewMeeting test success");
    }
}


// test join user
func TestJoinUser(t *testing.T) {
    if err := JoinUser("Meeting","test"); err != nil {
		t.Error("JoinUser test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("JoinUser test success");
    }
}


// test move out user
func TestMoveoutUser(t *testing.T) {
    if err := MoveoutUser("Meeting","test"); err != nil {
		t.Error("MoveoutUser test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("MoveoutUser test success");
    }
}

// test cancel meeting
func TestCancelMeeting(t *testing.T) {
    if err := CancelMeeting("Meeting"); err != nil {
		t.Error("CancelMeeting test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("CancelMeeting test success");
    }
}


// test quit meeting
func TestQuitMeeting(t *testing.T) {
    if err := QuitMeeting("Meeting"); err != nil {
		t.Error("QuitMeeting test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("QuitMeeting test success");
    }
}


// test clear all meetings
func TestClearAllMeetings(t *testing.T) {
    if err := ClearAllMeetings(); err != nil {
		t.Error("ClearAllMeetings test failed.");
		fmt.Println("error: " , err);
    } else {
        t.Log("ClearAllMeetings test success");
    }
}