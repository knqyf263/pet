package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("color", cmd.Flags().Lookup("color"))
		viper.BindPFlag("query", cmd.Flags().Lookup("query"))
		viper.BindPFlag("command", cmd.Flags().Lookup("command"))
	},
	RunE: execute,
}

func execute(cmd *cobra.Command, args []string) (err error) {

	var options []string
	if viper.GetString("query") != "" {
		options = append(options, fmt.Sprintf("--query %s", viper.GetString("query")))
	}

	commands, err := filter(options)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if viper.GetBool("command") && command != "" {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	return run(command, os.Stdin, os.Stdout)
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	execCmd.Flags().StringP("query", "q", "", `Initial value for query`)
	execCmd.Flags().BoolP("command", "c", false, `Show the command with the plain text before executing`)
}
