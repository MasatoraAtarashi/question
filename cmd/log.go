package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func runLogCmd(cmd *cobra.Command) (err error) {
	home_dir := os.Getenv("HOME")
	log_path := home_dir + "/.question/logs/HEAD"
	f, err := os.Open(log_path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		slice := strings.Split(scanner.Text(), " ")
		fmt.Printf("Question: " + slice[0] + "\n")
		fmt.Printf("Author: " + slice[1] + "\n")
		fmt.Printf("Date: " + slice[2] + " " + slice[3] + "\n")
		fmt.Printf("\n\t" + slice[len(slice)-1] + "\n\n")
	}
	return
}

func init() {
	rootCmd.AddCommand(logCmd)
}
