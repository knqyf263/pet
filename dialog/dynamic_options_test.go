//go:build !windows

package dialog_test

import (
	"slices"
	"testing"

	"github.com/knqyf263/pet/dialog"
)

const nameOfThisFile = "dynamic_options_test.go"

func TestEvaluator(t *testing.T) {
	param := "ls"

	e, err := dialog.DynamicOptions(param)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(e) == 0 {
		t.Log("Expected at least one result, but got none.")
		t.FailNow()
	}

	if !slices.Contains(e, nameOfThisFile) {
		t.Log("it should have the current file in this path")
		t.FailNow()
	}
}
