package cmd

import (
	"fmt"
	"strings"

	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var delimiter string
var current_cmd_name string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search snippets",
	Long:  `Search snippets interactively (default filtering tool: peco)`,
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) (err error) {
	current_cmd_name = cmd.Name()
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", flag.Query))
	}
	options = append(options, fmt.Sprintf("--layout %s", flag.Layout))
	commands, err := filter(options)
	if err != nil {
		return err
	}

	fmt.Print(strings.Join(commands, flag.Delimiter))
	if terminal.IsTerminal(1) {
		fmt.Print("\n")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	searchCmd.Flags().StringVarP(&config.Flag.Layout, "layout", "", "top-down",
		`layout to be used. 'top-down' or 'bottom-up'`)
	searchCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}
