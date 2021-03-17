package cmd

import (
	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Edit config file",
	Long:  `Edit config file (default: opened by vim)`,
	RunE:  configure,
}

func configure(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	return editFile(editor, configFile, 0)
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
