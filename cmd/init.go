package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		runInitCmd(cmd, args)
	},
}

func runInitCmd(cmd *cobra.Command, args []string) (err error) {
	fmt.Print("init called by ")
	name, err := cmd.PersistentFlags().GetString("name")
	if err != nil {
		return err
	}
	if name != "" {
		fmt.Println(name)
	}
	return
}

func init() {
	initCmd.PersistentFlags().StringP("name", "n", "", "your name")
	rootCmd.AddCommand(initCmd)
}
