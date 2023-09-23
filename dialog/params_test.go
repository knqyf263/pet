package dialog

import (
	"testing"
)

func TestNormalizeCommand(t *testing.T) {
	output := normalizeCommand("cmd <param=1>")
	expected := "cmd <param>"
	assertEquals(t, output, expected)

	output = normalizeCommand("cmd <param with spaces=value with spaces>")
	expected = "cmd <param with spaces>"
	assertEquals(t, output, expected)

	output = normalizeCommand(`echo "param1: <param1>, param1: <param1=1>, param2: <param2>, param2: <param2=2>"`)
	expected = `echo "param1: <param1>, param1: <param1>, param2: <param2>, param2: <param2>"`
	assertEquals(t, output, expected)
}

func assertEquals(t *testing.T, output string, expected string) {
	if output != expected {
		t.Fatalf("\nExpected\n%s\bto be equal to\n%s", expected, output)
	}
}
