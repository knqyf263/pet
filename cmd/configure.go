package cmd

import (
	"github.com/spf13/viper"
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
	editor := viper.GetString("general.editor")
	return editFile(editor, viper.ConfigFileUsed())
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
