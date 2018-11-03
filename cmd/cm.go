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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/cmd/utils"
	"github.com/zhanzongyuan/agenda/validate"
)

// cmCmd represents the cm command
var cmCmd = &cobra.Command{
	Use:   "cm",
	Short: "Command to create meeting.",
	Long: `You must login your account before create meeting,
Please input your meeting title, start time, stop time, participator name`,
	Run: func(cmd *cobra.Command, args []string) {
		title := ""
		startTime, endTime := time.Time{}, time.Time{}
		parsName := []string{}
		// Parse flag
		flag := cmd.Flags()
		title, err := flag.GetString("title")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		stStr, err := flag.GetString("start")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		if len(stStr) != 0 {
			startTime, err = time.Parse(utils.Layout, stStr)
			if err != nil {
				log.Fatal(err)
			}
		}
		etStr, err := flag.GetString("end")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		if len(etStr) != 0 {
			endTime, err = time.Parse(utils.Layout, etStr)
			if err != nil {
				log.Fatal(err)
			}
		}
		participateStr, err := flag.GetString("participators")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		parseStr := strings.Split(participateStr, ",")
		for _, s := range parseStr {
			if len(s) != 0 {
				parsName = append(parsName, s)
			}
		}

		// Auth
		_, err = agenda.Auth()
		if err != nil {
			log.Fatal(err)
		}

		// Check title
		fmt.Println("Please input meeting title with first letter capitalized.")
		utils.Scan(&title, "[Meeting title]: ", "[Meeting title]: ")
		if err := validate.IsTitleValid(title); err != nil {
			log.Fatal(err)
		}

		// Check start time and end time
		utils.ScanFtime(
			&startTime,
			fmt.Sprintf("[Start time] (e.g. %s): ", time.Now().Format(utils.Layout)),
			"[Start time]: ",
		)
		utils.ScanFtime(
			&endTime,
			fmt.Sprintf("[End time] (e.g. %s): ", time.Now().Format(utils.Layout)),
			"[End time]: ",
		)
		if err := validate.IsStartEndTimeValid(startTime, endTime); err != nil {
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
		fmt.Print("[Participators]: ")
		for _, s := range parsName {
			fmt.Print(s, ", ")
		}
		fmt.Println("")

		// Try to create new meeting
		m, err := agenda.NewMeeting(title, startTime, endTime, parsName)
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
	cmCmd.Flags().StringP("start", "s", "", "Meeting start time")
	cmCmd.Flags().StringP("end", "e", "", "Meeting end time")
	cmCmd.Flags().StringP("participators", "p", "", "Meeting participators list")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
