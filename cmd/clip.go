package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
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
		viper.BindPFlag("command", cmd.Flags().Lookup("command"))
		viper.BindPFlag("delimiter", cmd.Flags().Lookup("delimiter"))
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
	command := strings.Join(commands, viper.GetString("delimiter"))
	if viper.GetBool("command") && command != "" {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	return clipboard.WriteAll(command)
}

func init() {
	rootCmd.AddCommand(clipCmd)
	clipCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	clipCmd.Flags().StringP("query", "q", "", `Initial value for query`)
	clipCmd.Flags().BoolP("command", "c", false, `Show the command with the plain text before copying`)
	clipCmd.Flags().StringP("delimiter", "d", "; ", `Use delim as the command delimiter character`)
}
