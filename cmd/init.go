package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

const (
	QUESTION_TEMPLATE = `## 以下の項目を入力してください。

# 1. タイトル(必須)
## わからないこと・解決したいことを入力してください。


# 2. 実現したいこと(必須)


# 3. 発生している問題・エラーメッセージ(必須)


# 4. 発生している問題を再現する手順(必須)


# 5. 該当のソースコード(オプション)


# 6. 試したこと(必須)


# 7. 参考資料(オプション)


# 8. 詳細なログ(オプション)


# 9. 環境設定情報(オプション)
## 【マシン, メモリ量, 関連周辺機器, OS, 利用ソフト, バージョンなど】を箇条書きにしてください。


`
)

var requireUserInputLabels = map[string]string{
	"Subject":     "タイトル",
	"Ideal":       "実現したいこと",
	"Problem":     "発生している問題",
	"Procedure":   "発生している問題を再現する手順",
	"TriedAction": "試したこと",
}

type Options struct {
	Name    string
	Content string
}

type RequireUserInput struct {
	Subject     string
	Ideal       string
	Problem     string
	Procedure   string
	TriedAction string
}

type UserInput struct {
	RequireUserInput RequireUserInput
	Source           string
	Reference        string
	Log              string
	Env              string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create question",
	Run: func(cmd *cobra.Command, args []string) {
		err := runInitCmd(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	},
}

func runInitCmd(cmd *cobra.Command, args []string) (err error) {
	// ユーザからの入力を受け取る
	userInput, err := getUserInput()
	if err != nil {
		return
	}

	// ユーザ名を取得
	name, err := getName(cmd)
	if err != nil {
		return
	}

	// struct Questionにマッピング
	question, err := NewQuestion(userInput, name)
	if err != nil {
		return
	}

	// 質問文を生成
	result, err := question.Execute()
	if err != nil {
		return
	}

	// 生成された質問文をユーザに確認してもらう
	result, err = confirmQuestion(result)
	if err != nil {
		return
	}

	fmt.Printf(result)
	return
}

func getUserInput() (userInput UserInput, err error) {
	// make tmp file
	fpath, err := makeTmpFile(QUESTION_TEMPLATE)
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

func makeTmpFile(msg string) (fpath string, err error) {
	home := os.Getenv("HOME")
	fpath = filepath.Join(home, "QUESTION_EDITMSG")
	if err != nil {
		return
	}
	if !isFileExist(fpath) {
		err = ioutil.WriteFile(fpath, []byte(msg), 0644)
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
	var subject = ""
	var ideal = ""
	var problem = ""
	var procedure = ""
	var source = ""
	var tried_action = ""
	var reference = ""
	var log = ""
	var env = ""
	var variables = []*string{nil, &subject, &ideal, &problem, &procedure, &source, &tried_action, &reference, &log, &env}
	var i = 0

	reader := bytes.NewReader(content)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			if i > 0 {
				*variables[i] = strings.Trim(*variables[i], "\n")
			}
			i++
			if i >= len(variables) {
				break
			}
			continue
		}
		*variables[i] += line + "\n"
	}
	requireUserInput := RequireUserInput{
		Subject:     subject,
		Ideal:       ideal,
		Problem:     problem,
		Procedure:   procedure,
		TriedAction: tried_action,
	}
	valid, empty_vars, err := validateUserInput(&requireUserInput)
	if !valid {
		fmt.Fprintf(os.Stderr, "Abort!: %sは必須項目です", empty_vars)
		os.Exit(1)
	}
	userInput = UserInput{
		RequireUserInput: requireUserInput,
		Source:           source,
		Reference:        reference,
		Log:              log,
		Env:              env,
	}
	return
}

func validateUserInput(requireUserInput *RequireUserInput) (valid bool, empty_vars []string, err error) {
	valid = true
	u := reflect.TypeOf(*requireUserInput)
	elem := reflect.ValueOf(requireUserInput).Elem()
	cnt := elem.NumField()
	for i := 0; i < cnt; i++ {
		value := elem.Field(i).String()
		if value == "" {
			valid = false
			empty_vars = append(empty_vars, requireUserInputLabels[u.Field(i).Name])
		}
	}
	return
}

func confirmQuestion(draft string) (result string, err error) {
	const instruction = "###### 質問文を確認してください。 ######\n\n"
	msg := instruction + draft
	fpath, err := makeTmpFile(msg)
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("failed make edit file. %s\n", err.Error()))
		return
	}
	defer deleteFile(fpath)

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
	result = string(content)
	result = strings.TrimLeft(result, instruction)
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
	if name == "" {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		name = user.Username
	}
	return
}

func init() {
	// options
	initCmd.PersistentFlags().StringP("name", "n", "", "your name")

	rootCmd.AddCommand(initCmd)
}
