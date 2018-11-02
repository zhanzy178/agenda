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
	"crypto/sha256"
	"fmt"
	"log"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/validate"
	"golang.org/x/crypto/ssh/terminal"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Command login your account on Agenda system.",
	Long: `To run this command, you should input your
username and password. If you account is login in other bash,
then your operation will force others to logout. If you have logined in
current bash, you will only get warning.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Check if one user login in current bash
		if user := agenda.CurrentUser(); user != nil {
			log.Printf("Current bash has one user('%s'), If you want to login, please logout the user in current bash first.\n", user.Name)
			fmt.Println(user)
			return
		}

		// Parse flag username
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}

		// Input username
		if len(username) == 0 {
			fmt.Print("[Username]:")
			_, err := fmt.Scan(&username)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("[Username]:", username)
		}
		if err := validate.IsNameValid(username); err != nil {
			log.Fatal(err)
		}
		// Input password
		fmt.Print("[Password]:")
		bytePass, err := terminal.ReadPassword(int(syscall.Stdin))
		password := string(bytePass)
		fmt.Println("")
		if err := validate.IsPasswordValid(password); err != nil {
			log.Fatal(err)
		}
		sha := sha256.New()
		sha.Write(bytePass)
		password = fmt.Sprintf("%x", sha.Sum(nil))

		// Login
		user, err := agenda.Login(username, password)
		if err != nil {
			log.Fatal(err)
		}
		if user != nil {
			fmt.Println(user)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username to login")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
