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
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhanzongyuan/agenda/agenda"
	"github.com/zhanzongyuan/agenda/validate"
	"golang.org/x/crypto/ssh/terminal"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Command register your account for agenda system.",
	Long: `This is a command for agenda account register. 
You can make use of this system to help you to manage meeting.
You can register with command format in example.
E.g.

	agenda register -uYourName -pYourPassword -eYourEmail -nYourNumber

`,
	Run: func(cmd *cobra.Command, args []string) {
		flag := cmd.Flags()
		// Register User

		// Check username
		username, err := flag.GetString("username")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		fmt.Println("Please input your username with first letter capitalized.")
		fmt.Print("[Username]: ")
		if len(username) == 0 {
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			username = s.Text()
		} else {
			fmt.Println(username)
		}
		if err := validate.IsNameValid(username); err != nil {
			log.Fatal(err)
		}

		// Check password
		password, err := flag.GetString("password")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		if len(password) == 0 {
			fmt.Print("[Password]:")
			bytePass, _ := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println("")
			password = string(bytePass)
			if err := validate.IsPasswordValid(password); err != nil {
				log.Fatal(err)
			}
			fmt.Print("[Check Password]:")
			bytePassAg, _ := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println("")
			if passwordAg := string(bytePassAg); passwordAg != password {
				log.Fatal(errors.New("The second password is not same as the first one!"))
			}
		} else {
			log.Println("Use command flag password")
			if err := validate.IsPasswordValid(password); err != nil {
				log.Fatal(err)
			}
		}
		sha := sha256.New()
		sha.Write([]byte(password))
		password = fmt.Sprintf("%x", sha.Sum(nil))

		// Check email
		email, err := flag.GetString("email")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		fmt.Print("[E-mail]: ")
		if len(email) == 0 {
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			email = s.Text()
		} else {
			fmt.Println(email)
		}
		if err := validate.IsEmailValid(email); err != nil {
			log.Fatal(err)
		}

		// Check Number
		number, err := flag.GetString("number")
		if err != nil {
			cmd.Help()
			log.Fatal(err)
		}
		fmt.Print("[Number]: ")
		if len(number) == 0 {
			s := bufio.NewScanner(os.Stdin)
			s.Scan()
			number = s.Text()
		} else {
			fmt.Println(number)
		}
		if err := validate.IsNumberValid(number); len(number) != 0 && err != nil {
			log.Fatal(err)
		}
		// Register
		user, err := agenda.Register(username, password, email, number)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Println("Your account register successfully:")
		fmt.Println(*user)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.
	// TODO: Add Flag here
	registerCmd.Flags().StringP("username", "u", "", "Username for register, it must be unique and not null")
	registerCmd.Flags().StringP("password", "p", "", "Password for register")
	registerCmd.Flags().StringP("email", "e", "", "Email for your account, it must be valid in format 'xxx@xx.xx'")
	registerCmd.Flags().StringP("number", "n", "", "Number for your account. Default null.")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
