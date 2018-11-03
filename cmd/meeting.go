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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/cmd/utils"
)

// meetingCmd represents the meeting command
var meetingCmd = &cobra.Command{
	Use:   "meeting",
	Short: "Command list meeting table you specific during time interval.",
	Long: `Using this command you can get a table of the meeting that you
initial or you join. You can also setting time interval to filer this table.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime, endTime := time.Time{}, time.Time{}
		// Parse flag
		flag := cmd.Flags()
		showAll, err := flag.GetBool("all")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}

		if !showAll {
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
			if startTime.After(endTime) {
				log.Fatal(errors.New("Start time must be earlier than end time."))
			}

			if _, err := agenda.CheckMeetings(startTime, endTime); err != nil {
				log.Fatal(err)
			}
		} else {
			if _, err := agenda.AllMeetings(); err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(meetingCmd)

	// Here you will define your flags and configuration settings.
	meetingCmd.Flags().StringP("start", "s", "", "Meeting start time")
	meetingCmd.Flags().StringP("end", "e", "", "Meeting end time")
	meetingCmd.Flags().BoolP("all", "a", false, "Flag setting to list all")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// meetingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// meetingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
