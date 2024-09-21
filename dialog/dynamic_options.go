package dialog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	runner "github.com/knqyf263/pet/cmd/runner"
)

// DynamicOptions run a command, split and returns an array of strings. It expect to
// receive command substitions.
//
// Example:
//
//	DynamicOptions("$(fd -tf --hidden --no-ignore --max-depth=1 .)")
func DynamicOptions(what string) ([]string, error) {
	if what[:2] != "$(" || what[len(what)-1] != ')' {
		return nil, fmt.Errorf("no evaluated command found: %v", what)
	}

	var w bytes.Buffer
	err := runner.Run(what[2:len(what)-1], os.Stdin, &w)
	if err != nil {
		log.Fatal(what)
		return nil, err
	}

	return strings.Split(strings.TrimSpace(w.String()), "\n"), nil
}
