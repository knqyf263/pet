package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

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

func _execute(in io.ReadCloser, out io.Writer) (err error) {
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

	// Show final command before executing it
	if !flag.Silent {
		fmt.Fprintf(out, "> %s\n", command)
	}

	return run(command, in, out)
}

func execute(cmd *cobra.Command, args []string) error {
	return _execute(os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().StringVarP(&config.Flag.FilterTag, "tags", "t", "", "Filter by specified tags as comma separated values")
	execCmd.Flags().BoolVarP(&config.Flag.Silent, "silent", "s", false,
		`Suppress the command output`)
}
