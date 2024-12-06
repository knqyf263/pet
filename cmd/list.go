package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snippets",
	Long:  `Show all snippets`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var snippets snippet.Snippets
	if err := snippets.Load(true); err != nil {
		return err
	}

	if config.Flag.FilterTag != "" {
		snippets.FilterByTags(strings.Split(config.Flag.FilterTag, ","))
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, snippet := range snippets.Snippets {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(snippet.Description, col, "..."), col)
			command := snippet.Command
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.HiGreenString(description), color.HiYellowString(command))
		} else {
			if config.Flag.Debug {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.RedString("   Filename:"), snippet.Filename)
			}
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.HiGreenString("Description:"), snippet.Description)
			if strings.Contains(snippet.Command, "\n") {
				lines := strings.Split(snippet.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.HiYellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.HiYellowString("    Command:"), snippet.Command)
			}
			if snippet.Tag != nil {
				tag := strings.Join(snippet.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.HiCyanString("        Tag:"), tag)
			}
			if snippet.Output != "" {
				output := strings.Replace(snippet.Output, "\n", "\n             ", -1)
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.HiRedString("     Output:"), output)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&config.Flag.OneLine, "oneline", "", false,
		`Display snippets in one line`)
	listCmd.Flags().StringVar(&config.Flag.FilterTag, "t", "", "list by specified tags as comma separated values")
}
