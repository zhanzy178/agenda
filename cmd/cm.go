// Copyright © 2018 Zongyuan Zhan <zhanzy5@mail2.sysu.edu.cn>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/validate"
)

// Time for flag parse
type Time time.Time

func (t *Time) Type() string {
	return fmt.Sprintf("%T", *t)
}
func (t *Time) Set(s string) error {
	ti, err := time.Parse(timeLayout, s)
	if err != nil {
		return err
	}
	*t = Time(ti)
	return nil
}
func (t *Time) String() string {
	return time.Time(*t).Format("2006-01-02 15:04:05 Mon")
}

// global variable
var (
	startTime  Time
	endTime    Time
	nullTime   time.Time
	parsName   []string
	timeLayout string = "2006-1-2 15:04"
)

// cmCmd represents the cm command
var cmCmd = &cobra.Command{
	Use:   "cm",
	Short: "Command to create meeting.",
	Long: `You must login your account before create meeting,
Please input your meeting title, start time, stop time, participator name`,
	Run: func(cmd *cobra.Command, args []string) {
		flag := cmd.Flags()
		curUser, err := agenda.Auth()
		if err != nil {
			log.Fatal(err)
		}
		// Check title
		title, err := flag.GetString("title")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		fmt.Print("[Meeting title]: ")
		if title == "" {
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			title = s.Text()
		} else {
			fmt.Println(title)
		}
		if err := validate.IsTitleValid(title); err != nil {
			log.Fatal(err)
		}

		// Check start time and end time
		st, et := time.Time(startTime), time.Time(endTime)

		if st == nullTime {
			fmt.Printf("[Start time] (e.g. %s): ", time.Now().Format(timeLayout))
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			stStr := s.Text()
			ti, err := time.Parse(timeLayout, stStr)
			if err != nil {
				log.Fatal(err)
			}
			st = ti
		} else {
			fmt.Print("[Start time]: ")
			fmt.Println(st.Format(timeLayout))
		}
		if et == nullTime {
			fmt.Printf("[End time] (e.g. %s): ", time.Now().Format(timeLayout))
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			etStr := s.Text()
			ti, err := time.Parse(timeLayout, etStr)
			if err != nil {
				log.Fatal(err)
			}
			et = ti
		} else {
			fmt.Print("[End time]: ")
			fmt.Println(et.Format(timeLayout))
		}
		if err := validate.IsStartEndTimeValid(st, et); err != nil {
			log.Fatal(err)
		}

		// Check participators
		if len(parsName) == 0 {
			name := "ini"
			s := bufio.NewScanner(os.Stdin)
			for len(name) != 0 {
				fmt.Print("[Input participator] (input nothing and press <enter> to complete): ")
				s.Scan()
				name = s.Text()
				err := validate.IsNameValid(name)
				if len(name) != 0 {
					if err == nil {
						parsName = append(parsName, name)
					} else {
						log.Println(err)
					}
				}
			}
		}
		parsName = append(parsName, curUser.Name)
		fmt.Print("[Participators]: ")
		for _, s := range parsName {
			fmt.Print(s, ", ")
		}
		fmt.Println("")

		// Try to create new meeting
		m, err := agenda.NewMeeting(title, st, et, parsName)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Create meeting succesfully!")
		fmt.Println(m)

	},
}

func init() {
	rootCmd.AddCommand(cmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	cmCmd.Flags().StringP("title", "t", "", "Meeting title")
	cmCmd.Flags().VarP(&startTime, "start", "s", "Meeting start time")
	cmCmd.Flags().VarP(&endTime, "end", "e", "Meeting end time")
	cmCmd.Flags().StringArrayP("participators", "p", parsName, "Meeting participators list")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
