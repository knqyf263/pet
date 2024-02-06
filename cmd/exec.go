package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	commands, err := filter(options, flag.FilterTag)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}
	fmt.Printf("> %s\n", command)

	return run(command, os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
	execCmd.Flags().BoolVarP(&config.Flag.Command, "command", "c", false,
		`Show the command with the plain text before executing`)
}
