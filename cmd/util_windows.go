//go:build windows

package cmd

import (
	"os"

	"github.com/knqyf263/pet/cmd/runner"
	"github.com/knqyf263/pet/path"
)

func editFile(command string, filePath path.AbsolutePath, startingLine int) error {
	command += " " + filePath.Get()
	return runner.Run(command, os.Stdin, os.Stdout)
}
