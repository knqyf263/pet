package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ramiawar/superpet/config"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Load the selected environment",
	Long:  `Load the selected environment`,
	RunE:  activate,
}

func activate(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	envs, err := filterEnv(options, flag.FilterTag)
	if err != nil {
		return err
	}

	// Create new shell and add env vars to it (overwriting old values)
	ex := exec.Command("zsh", "-i")
	ex.Env = os.Environ()
	ex.Env = append(ex.Env, envs...)

	// Attach stdin, stdout to new shell
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	fmt.Println("\nactivating environment...")

	err = ex.Run() // add error checking

	fmt.Println("exited superpet shell")

	return err
}

func init() {
	RootCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	activateCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
}
