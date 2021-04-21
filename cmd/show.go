package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/icza/backscanner"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [<object>]",
	Short: "Show the specified question",
	Run: func(cmd *cobra.Command, args []string) {
		err := runShowCmd(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	},
}

func runShowCmd(cobraCmd *cobra.Command, args []string) (err error) {
	if len(args) >= 2 {
		return errors.New("Invalid argument\n")
	}
	if len(args) == 1 {
		showQuestion(args[0])
	} else {
		showLatestQuestion()
	}
	return nil
}

// 最新の質問を表示する
func showLatestQuestion() (err error) {
	qId, err := getLatestQuestionId()
	if err != nil {
		return
	}

	path, err := getQuestionPath(qId)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	output := "Question: " + qId + "\n\n" + string(bytes)
	err = execLess(output)
	return
}

// 最新の質問のハッシュを取得する
func getLatestQuestionId() (qId string, err error) {
	home_dir := os.Getenv("HOME")
	log_path := home_dir + "/.question/logs/HEAD"
	bytes, err := ioutil.ReadFile(log_path)
	input := string(bytes)
	scanner := backscanner.New(strings.NewReader(input), len(input))
	// 改行をスキップ
	scanner.Line()
	line, _, err := scanner.Line()
	slice := strings.Split(line, " ")
	qId = slice[0]
	return
}

// 指定されたハッシュをidに持つ質問へのファイルパスを取得する
func getQuestionPath(qId string) (path string, err error) {
	home_dir := os.Getenv("HOME")
	path = home_dir + "/.question/objects/" + qId
	return
}

// 指定されたハッシュをidに持つ質問文を表示する
func showQuestion(qId string) (err error) {
	path, err := getQuestionPath(qId)
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	output := "Question: " + qId + "\n\n" + string(bytes)
	err = execLess(output)
	return
}

func init() {
	rootCmd.AddCommand(showCmd)
}
