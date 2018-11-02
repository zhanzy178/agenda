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
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/cmd/utils"
)

// moveoutCmd represents the moveout command
var moveoutCmd = &cobra.Command{
	Use:   "moveout",
	Short: "Command move out user from meeting participators",
	Long: `If you are the sponsor of a meeting and
login in the bash, ou can move out participator from the meeting.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flag
		flag := cmd.Flags()
		title, err := flag.GetString("title")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		username, err := flag.GetString("username")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}

		// Title
		utils.Scan(&title, "[Title]: ", "[Title]: ")

		// Username
		utils.Scan(&username, "[Username]: ", "[Username]: ")

		// ensure
		ensure := ""
		for strings.ToLower(ensure) != "y" && strings.ToLower(ensure) != "n" {
			fmt.Printf("Are you should to move out user '%s'?[Y/n]", username)
			fmt.Scanln(&ensure)
			if ensure == "" {
				ensure = "y"
			}
		}

		if strings.ToLower(ensure) == "y" {
			if err := agenda.MoveoutUser(title, username); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("User '%s' is moved out from meeting '%s' successfully!", username, title)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(moveoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveoutCmd.PersistentFlags().String("foo", "", "A help for foo")
	moveoutCmd.Flags().StringP("username", "u", "", "User you decide to move out from meeting")
	moveoutCmd.Flags().StringP("title", "t", "", "Meeting title you specific")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
