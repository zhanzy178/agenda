package validate

import (
	"errors"
	"regexp"
	"time"
)

var (
	MinNameLen     int = 3
	MinPasswordLen int = 4
	EmailRegexp        = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	NumberRegexp       = `^1[34578]\d{9}$`
)

func IsTitleValid(title string) error {
	if len(title) < 3 {
		return errors.New("Title number must be longer than 2!")
	}
	if !(title[0] <= 'Z' && title[0] >= 'A') {
		return errors.New("First letter of title must be capitalized!")
	}
	return nil
}
func IsStartEndTimeValid(start time.Time, end time.Time) error {
	if start.After(end) {
		return errors.New("Meeting start time must be earlier than end time!")
	} else if start.Before(time.Now()) {
		return errors.New("Meeting start time should not be past time!")
	}
	return nil
}

func IsNameValid(name string) error {
	if len(name) < 3 {
		return errors.New("Username number must be longer than 2!")
	}
	if !(name[0] <= 'Z' && name[0] >= 'A') {
		return errors.New("First letter of username must be capitalized!")
	}
	return nil
}

func IsPasswordValid(password string) error {
	if len(password) < MinPasswordLen {
		return errors.New("password number must longer than 3!")
	}
	return nil
}

func IsEmailValid(email string) error {
	validEmail := regexp.MustCompile(EmailRegexp)

	if !validEmail.MatchString(email) {
		return errors.New("Invalid email!")
	}
	return nil
}

func IsNumberValid(number string) error {
	validNumber := regexp.MustCompile(NumberRegexp)

	if len(number) != 0 && !validNumber.MatchString(number) {
		return errors.New("Invalid number!")
	}
	return nil
}
