package dialog

import (
	"bytes"
	"os"
	"strings"

	runner "github.com/knqyf263/pet/cmd/runner"
)

// DynamicOptions runs a command, split by lines and returns an array of
// strings. It expect to receive the command directly.
//
// Example:
//
//	DynamicOptions("fd -tf --hidden --no-ignore --max-depth=1 .")
func DynamicOptions(rawcmd string) ([]string, error) {
	var w bytes.Buffer
	err := runner.Run(rawcmd, os.Stdin, &w)
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(w.String()), "\n"), nil
}

// hasDynamicOptions returns true and the command, if any, otherwise, returns
// false and the original argument.
//
// See [github.com/knqyf263/pet/dialog.DynamicOptions].
func hasDynamicOptions(opt string) (bool, string) {
	if opt[:2] == "$(" && opt[len(opt)-1] == ')' {
		return true, opt[2 : len(opt)-1]
	}
	return false, opt
}
