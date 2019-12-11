package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
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
	if viper.GetString("query") != "" {
		options = append(options, fmt.Sprintf("--query %s", viper.GetString("query")))
	}

	commands, err := filter(options)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if debug {
		fmt.Printf("Command: %s\n", command)
	}
	if viper.GetBool("command") {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	return run(command, os.Stdin, os.Stdout)
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	execCmd.Flags().StringP("query", "q", "", `Initial value for query`)
	execCmd.Flags().BoolP("command", "c", false, `Show the command with the plain text before executing`)
	viper.BindPFlag("color", execCmd.Flags().Lookup("color"))
	viper.BindPFlag("query", execCmd.Flags().Lookup("query"))
	viper.BindPFlag("command", execCmd.Flags().Lookup("command"))
}
