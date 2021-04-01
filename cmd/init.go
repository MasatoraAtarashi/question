package cmd

import (
	"github.com/spf13/cobra"
)

type Options struct {
	Name string
	Content string
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
	name, err := cmd.PersistentFlags().GetString("name")
	if err != nil {
		return err
	}
	question, err := NewQuestion(name)
	err = question.render()
	if err != nil {
		return err
	}
	return
}

func init() {
	initCmd.PersistentFlags().StringP("name", "n", "", "your name")
	rootCmd.AddCommand(initCmd)
}
