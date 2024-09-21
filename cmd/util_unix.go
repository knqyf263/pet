//go:build !windows

package cmd

import (
	"os"
	"strconv"

	"github.com/knqyf263/pet/cmd/runner"
	"github.com/knqyf263/pet/path"
)

func editFile(command string, filePath path.AbsolutePath, startingLine int) error {
	command += " +" + strconv.Itoa(startingLine) + " " + filePath.Get()
	return runner.Run(command, os.Stdin, os.Stdout)
}
