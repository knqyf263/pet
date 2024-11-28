//go:build !windows

package cmd

import (
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/path"
)

func run(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if len(config.Conf.General.Cmd) > 0 {
		line := append(config.Conf.General.Cmd, command)
		cmd = exec.Command(line[0], line[1:]...)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func editFile(command string, filePath path.AbsolutePath, startingLine int) error {
	command += " +" + strconv.Itoa(startingLine) + " " + filePath.Get()
	return run(command, os.Stdin, os.Stdout)
}
