package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search snippets",
	Long:  `Search snippets interactively (default filtering tool: peco)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("color", cmd.Flags().Lookup("color"))
		viper.BindPFlag("query", cmd.Flags().Lookup("query"))
		viper.BindPFlag("delimiter", cmd.Flags().Lookup("delimiter"))
	},
	RunE: search,
}

func search(cmd *cobra.Command, args []string) (err error) {
	var options []string
	if viper.GetString("query") != "" {
		options = append(options, fmt.Sprintf("--query %s", viper.GetString("query")))
	}

	commands, err := filter(options)
	if err != nil {
		return err
	}

	fmt.Print(strings.Join(commands, viper.GetString("delimiter")))
	if terminal.IsTerminal(1) {
		fmt.Print("\n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolP("color", "", false, `Enable colorized output (only fzf)`)
	searchCmd.Flags().StringP("query", "q", "", `Initial value for query`)
	searchCmd.Flags().StringP("delimiter", "d", "; ", `Use delim as the command delimiter character`)
}
