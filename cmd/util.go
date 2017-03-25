package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"

	"github.com/knqyf263/pet/dialog"
)

func editFile(command, file string) error {
	command += " " + file
	return run(command, os.Stdin, os.Stdout)
}

func run(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func filter(options []string) (commands []string, err error) {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return commands, fmt.Errorf("Load snippet failed: %v", err)
	}

	snippetTexts := map[string]snippet.SnippetInfo{}
	var text string
	for _, s := range snippets.Snippets {
		t := fmt.Sprintf("[%s] %s", s.Description, s.Command)
		snippetTexts[t] = s
		text += t + "\n"
	}

	var buf bytes.Buffer
	selectCmd := fmt.Sprintf("%s %s",
		config.Conf.General.SelectCmd, strings.Join(options, " "))
	err = run(selectCmd, strings.NewReader(text), &buf)
	if err != nil {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	dialog.Params = dialog.SearchForParams(lines)
	if dialog.Params == nil {
		for _, line := range lines {
			snippetInfo := snippetTexts[line]
			commands = append(commands, fmt.Sprint(snippetInfo.Command))
		}
		return commands, nil
	} else {
		snippetInfo := snippetTexts[lines[0]]
		dialog.Current_command = snippetInfo.Command
		dialog.GenerateParamsLayout(dialog.Params, dialog.Current_command)
		res := []string{dialog.Final_command}
		return res, nil
	}
}
