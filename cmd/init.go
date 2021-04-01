package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	template = `#################################
## 以下の8項目を入力してください。
#################################
# 1. 概要
## (例): ○○を実行すると、○○というエラーになる問題で困っています。


# 2. 発生している問題
## エラーメッセージやキャプチャを入力してください。


# 3. 発生している問題を再現する手順
## (例): 
## (1) XXXX.cgiをhttp://xxxx からダウンロードする。
## (2) 管理ファイル名$adminの値をadmin.datからadmin.txtに変更する。


# 4. 期待していた結果


# 5. 参考資料


# 6. 問題解決のために自分自身で行ったこと
## (例): 
## (1) 入力を○○ではなく××にしてみた
## →上記と同じ結果になった


# 7. 詳細なログ


# 8. 環境設定情報
## 【マシン, メモリ量, 関連周辺機器, OS, 利用ソフト, バージョンなど】を箇条書きにしてください。


`
)

type Options struct {
	Name    string
	Content string
}

type UserInput struct {
	Overview    string
	Probrem     string
	Procedure   string
	Expected    string
	Reference   string
	TriedAction string
	Log         string
	Env         string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		runInitCmd(cmd, args)
	},
}

func runInitCmd(cmd *cobra.Command, args []string) (err error) {
	userInput, err := getUserInput()

	if err != nil {
		return
	}

	name, err := getName(cmd)
	if err != nil {
		return
	}

	question, err := NewQuestion(userInput, name)
	err = question.render()
	if err != nil {
		return
	}

	return
}

func getUserInput() (userInput UserInput, err error) {
	// make tmp file
	fpath, err := makeTmpFile()
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed make edit file. %s\n", err.Error()))
		return
	}
	defer deleteFile(fpath)

	// open text editor
	err = openEditor("vim", fpath)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed open text editor. %s\n", err.Error()))
		return
	}

	// read edit file
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed read content. %s\n", err.Error()))
		return
	}

	userInput, err = parseUserInput(content)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("faild parse input. %s\n", err.Error()))
		return
	}

	return
}

func makeTmpFile() (fpath string, err error) {
	home := os.Getenv("HOME")
	fpath = filepath.Join(home, "QUESTION_EDITMSG")
	if err != nil {
		return
	}
	if !isFileExist(fpath) {
		err = ioutil.WriteFile(fpath, []byte(template), 0644)
		if err != nil {
			return
		}
	}
	return
}

func isFileExist(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil || !os.IsNotExist(err)
}

func deleteFile(fpath string) error {
	return os.Remove(fpath)
}

func openEditor(program string, fpath string) error {
	c := exec.Command(program, fpath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func parseUserInput(content []byte) (userInput UserInput, err error) {
	var overview = ""
	var probrem = ""
	var procedure = ""
	var expected = ""
	var reference = ""
	var tried_action = ""
	var log = ""
	var env = ""
	var variables = []*string{nil, &overview, &probrem, &procedure, &expected, &reference, &tried_action, &log, &env}
	var i = 0

	reader := bytes.NewReader(content)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			i++
			if i >= len(variables) {
				break
			}
			continue
		}
		*variables[i] += line + "\n"
	}
	userInput = UserInput{
		Overview:    overview,
		Probrem:     probrem,
		Procedure:   procedure,
		Expected:    expected,
		Reference:   reference,
		TriedAction: tried_action,
		Log:         log,
		Env:         env,
	}
	return
}

// get name from option or config
func getName(cmd *cobra.Command) (name string, err error) {
	name, err = cmd.PersistentFlags().GetString("name")
	if err != nil {
		return
	}
	if name == "" {
		name = config.Question.Name
	}
	return
}

func init() {
	// options
	initCmd.PersistentFlags().StringP("name", "n", "", "your name")

	rootCmd.AddCommand(initCmd)
}
