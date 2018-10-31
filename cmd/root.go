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
	"os"
	"path/filepath"

	"github.com/zhanzongyuan/agenda/agenda"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dataDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agenda",
	Short: "Create meeting schedule for you.",
	Long:  `Agenda is an useful CLI program for everyone to manage meeting.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// If we don't definite Run(), then command will be not runnable and execute unrunnable will get errhelp.
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// .agenda.yaml (define your agenda data dir)
	// E.g.
	// agenda_data_dir: directory/path/to/your/data/file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agenda.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
			return
		}
		// Search config in home directory with name ".agenda" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".agenda")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Prepare data directory.
	if dataDir = viper.GetString("agendaDataRoot"); dataDir == "" {
		// Use default data directory '$HOME/.agenda'
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
			return
		}
		dataDir = filepath.Join(home, ".agenda")
	}
	fi, err := os.Lstat(dataDir)
	if err != nil || !fi.Mode().IsDir() {
		// Directory is not exist, mkdir one
		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Create data directory: %s\n", dataDir)
	}
	log.Printf("Data directory prepared: %s\n", dataDir)

	// Config data directory for agenda
	if err := agenda.InitConfig(dataDir); err != nil {
		log.Fatal(err)
		return
	}

}
