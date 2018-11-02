package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	ftime "github.com/zhanzongyuan/agenda/cmd/flagtime"
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

func ScanFtime(t *ftime.Time, emptyMsg, nemptyMsg string) {
	tt := (*time.Time)(t)

	if *tt == ftime.NullTime {
		fmt.Print(emptyMsg)
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		stStr := s.Text()
		ti, err := time.Parse(ftime.Layout, stStr)
		if err != nil {
			log.Fatal(err)
		}
		*tt = ti
	} else {
		fmt.Print(nemptyMsg)
		fmt.Println(tt.Format(ftime.Layout))
	}
}
