package flagtimen

import (
	"fmt"
	"time"
)

type Time time.Time

var (
	Layout   string = "2006-1-2 15:04"
	NullTime time.Time
)

func (t *Time) Type() string {
	return fmt.Sprintf("%T", *t)
}
func (t *Time) Set(s string) error {
	ti, err := time.Parse(Layout, s)
	if err != nil {
		return err
	}
	*t = Time(ti)
	return nil
}
func (t *Time) String() string {
	return time.Time(*t).Format(Layout)
}
