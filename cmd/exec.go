package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) (err error) {
	var options []string
	commands, err := filter(options)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}
	return run(command, os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
}
