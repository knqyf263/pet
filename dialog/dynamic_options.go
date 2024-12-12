package dialog

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	runner "github.com/knqyf263/pet/cmd/runner"
)

var ReDynamicOptions = regexp.MustCompile(`\$\((.*)\)`)

// DynamicOptions runs a command, split by lines and returns an array of
// strings.
//
// Although it uses the same syntax as some shells (like bash), it doesn't
// evaluate outputs internally. It'll catch the outermost command, and run it
// directly in your default shell.
// This means that if you plan to do several command substitutions nested,
// this function will take the outermost string and send it to your shell,
// delegating the inner substitutions to be handled by your shell.
//
// Therefore, if you're some shell that supports it
// (such as bash; ash; zsh; dash) you can keep the same syntax, and your
// shell will do the rest, otherwise, you may need to change it a bit, or
// perhaps doesn't even support it (like cmd).
//
// Example:
//
//	// When running under a shell that supports command substitution
//	DynamicOptions(`$(fd -tf --hidden --no-ignore --max-depth=1 .)`)
//	DynamicOptions(`prefix $(echo "$(ls)") suffix`)
//
//	// When running in a csh (supports, but needs a bit change)
//	DynamicOptions("$(echo \"`ls`\")")
//
//	// When running in cmd
//	DynamicOptions(`$(dir)`) // it doesn't support command substitution
func DynamicOptions(rawcmd string) ([]string, error) {
	indx := ReDynamicOptions.FindIndex([]byte(rawcmd))
	if len(indx) == 0 {
		return nil, fmt.Errorf("dynamic options: no command")
	}
	indx = []int{indx[0] + 2, indx[1] - 1} // skip $()

	cmd := strings.TrimSpace(string(rawcmd[indx[0]:indx[1]]))

	var w bytes.Buffer
	err := runner.Run(cmd, os.Stdin, &w)
	if err != nil {
		return nil, err
	}

	prefix, suffix := rawcmd[:indx[0]-2], rawcmd[indx[1]+1:]

	res := make([]string, 0)
	for _, line := range strings.Split(w.String(), "\n") {
		res = append(res, fmt.Sprintf("%s%s%s", prefix, line, suffix))
	}

	return res, nil
}

// hasDynamicOptions returns true and the command, if any, otherwise, returns
// false and the original argument.
//
// See [github.com/knqyf263/pet/dialog.DynamicOptions].
func hasDynamicOptions(opt string) bool { return ReDynamicOptions.Match([]byte(opt)) }
