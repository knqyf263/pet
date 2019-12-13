package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/atotto/clipboard"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "Copy the selected commands",
	Long: `Copy the selected commands to clipboard`,
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
	if debug {
		fmt.Printf("Command: %s\n", command)
	}
	return clipboard.WriteAll(command)
}

func init() {
	rootCmd.AddCommand(clipCmd)
	clipCmd.Flags().StringP("query", "q", "", `Initial value for query`)
	viper.BindPFlag("query", clipCmd.Flags().Lookup("query"))
}
