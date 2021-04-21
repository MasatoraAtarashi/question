package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/icza/backscanner"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show question logs",
	Run: func(cmd *cobra.Command, args []string) {
		err := runLogCmd(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	},
}

func runLogCmd(cobraCmd *cobra.Command) (err error) {
	bytes, err := getDataFromLogFile()
	if err != nil {
		return
	}
	input := string(bytes)
	scanner := backscanner.New(strings.NewReader(input), len(input))
	var output string
	for {
		line, _, err := scanner.Line()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// ログ・ファイルを加工
		slice := strings.Split(line, " ")
		if len(slice) == 1 {
			continue
		}
		question_hash := slice[0]
		user_name := slice[1]
		date := slice[2]
		time := slice[3]
		subject := slice[len(slice)-1]
		output += "Question: " + question_hash + "\n"
		output += "Author: " + user_name + "\n"
		output += "Date: " + date + " " + time + "\n"
		output += "\n\t" + subject + "\n\n"
	}
	err = execLess(output)
	if err != nil {
		return
	}
	return
}

// ログ・ファイルからデータを取得
func getDataFromLogFile() (bytes []byte, err error) {
	home_dir := os.Getenv("HOME")
	log_path := home_dir + "/.question/logs/HEAD"
	bytes, err = ioutil.ReadFile(log_path)
	return
}

// lessコマンドでログを表示
func execLess(output string) (err error) {
	cmd := exec.Command("less", "-R")
	cmd.Stdin = strings.NewReader(output)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return
}

func init() {
	rootCmd.AddCommand(logCmd)
}
