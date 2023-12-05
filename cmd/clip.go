package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "Copy the selected commands",
	Long:  `Copy the selected commands to clipboard`,
	RunE:  clip,
}

func clip(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", flag.Query))
	}

	commands, err := filter(options, flag.FilterTag)
	if err != nil {
		return err
	}
	command := strings.Join(commands, flag.Delimiter)
	if flag.Command && command != "" {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	return clipboard.WriteAll(command)
}

func init() {
	RootCmd.AddCommand(clipCmd)
	clipCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	clipCmd.Flags().BoolVarP(&config.Flag.Command, "command", "", false,
		`Display snippets in one line`)
	clipCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
	clipCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
	clipCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf) (not working)`)
}
