package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Layout   string = "2006-1-2 15:04"
	NullTime time.Time
)

func Scan(target *string, emptyMsg, nemptyMsg string) {
	if *target == "" {
		fmt.Print(emptyMsg)
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		*target = s.Text()
	} else {
		fmt.Println(nemptyMsg, *target)
	}
}

func ScanFtime(t *time.Time, emptyMsg, nemptyMsg string) {

	if *t == NullTime {
		fmt.Print(emptyMsg)
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		stStr := s.Text()
		ti, err := time.Parse(Layout, stStr)
		if err != nil {
			log.Fatal(err)
		}
		*t = ti
	} else {
		fmt.Print(nemptyMsg)
		fmt.Println(t.Format(Layout))
	}
}
