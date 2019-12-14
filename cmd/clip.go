package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "Copy the selected commands",
	Long:  `Copy the selected commands to clipboard`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("color", cmd.Flags().Lookup("color"))
		viper.BindPFlag("query", cmd.Flags().Lookup("query"))
	},
	RunE: clip,
}

func clip(cmd *cobra.Command, args []string) (err error) {
	var options []string
	if viper.GetString("query") != "" {
		options = append(options, fmt.Sprintf("--query %s", viper.GetString("query")))
	}

	commands, err := filter(options)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if viper.GetBool("debug") {
		fmt.Printf("Command: %s\n", command)
	}
	return clipboard.WriteAll(command)
}

func init() {
	rootCmd.AddCommand(clipCmd)
	clipCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	clipCmd.Flags().StringP("query", "q", "", `Initial value for query`)
}
