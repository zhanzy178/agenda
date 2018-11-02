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

// quitCmd represents the quit command
var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Command to quit a meeting you participated in.",
	Long: `You can specific meeting by its title
to quit a certain meeting you participated in.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse flag
		flag := cmd.Flags()
		title, err := flag.GetString("title")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}

		// Title
		utils.Scan(&title, "[Title]: ", "[Title]: ")

		// ensure
		ensure := ""
		for strings.ToLower(ensure) != "y" && strings.ToLower(ensure) != "n" {
			fmt.Printf("Are you should to quit meeting '%s'?[Y/n]", title)
			fmt.Scanln(&ensure)
			if ensure == "" {
				ensure = "y"
			}
		}

		if strings.ToLower(ensure) == "y" {
			if err := agenda.QuitMeeting(title); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Quit from meeting '%s' successfully!", title)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(quitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quitCmd.PersistentFlags().String("foo", "", "A help for foo")
	quitCmd.Flags().StringP("title", "t", "", "Meeting title you specific")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
