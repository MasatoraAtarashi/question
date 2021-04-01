package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	Name string
	Content string
}

type UserInput struct {
	Situation string
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

	userInput, err = parseUserInput(string(content))
	if err != nil {
		fmt.Fprint(os.Stdout, fmt.Sprintf("faild parse input. %s\n", err.Error()))
		return
	}

	return
}

func makeTmpFile() (fpath string, err error) {
	home := os.Getenv("HOME")
	fpath = filepath.Join(home, "QUESTION_EDITMSG")
	if !isFileExist(fpath) {
		err = ioutil.WriteFile(fpath, []byte(""), 0644)
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

func parseUserInput(content string) (userInput UserInput, err error) {
	userInput = UserInput{
		Situation: content,
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
