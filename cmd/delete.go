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
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Command will delete your current account.",
	Long: `Please be careful about this operation.
If you delete your account, all meetings you initial will be cancel.`,
	Run: func(cmd *cobra.Command, args []string) {
		user := agenda.CurrentUser()
		username := user.Name
		if user == nil {
			log.Fatal(errors.New("You are not login!"))
		}

		ensure := ""
		for strings.ToLower(ensure) != "y" && strings.ToLower(ensure) != "n" {
			fmt.Printf("All your meetings will be canceled.\nAre you should to delete current user '%s'?[Y/n]", user.Name)
			fmt.Scanln(&ensure)
			if ensure == "" {
				ensure = "y"
			}
		}
		if strings.ToLower(ensure) == "y" {
			if err := agenda.DeleteUser(); err != nil {
				log.Fatal(err)
			}
			log.Printf("User '%s' is deleted from agenda system now!", username)
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
