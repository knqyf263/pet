package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"

	"github.com/knqyf263/pet/dialog"
)

func autoSync(file string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	fi, err := os.Stat(file)
	if err != nil {
		return err
	}

	client := githubClient()
	gist, _, err := client.Gists.Get(context.Background(), config.Conf.Gist.GistID)
	if err != nil {
		return err
	}
	local := fi.ModTime().UTC()
	remote := gist.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload()
	case remote.After(local):
		return update(gist)
	default:
		return nil
	}
}

func update(gist *github.Gist) error {
	var (
		content     = ""
		snippetFile = config.Conf.General.SnippetFile
		filename    = config.Conf.Gist.FileName
	)
	for _, file := range gist.Files {
		if *file.Filename == filename {
			content = *file.Content
		}
	}
	if content == "" {
		return fmt.Errorf("%s is empty", filename)
	}

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}
	body, err := snippets.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to update
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(snippetFile, []byte(content), os.ModePerm)
}

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
		t := fmt.Sprintf("[%s]: %s", s.Description, s.Command)

		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf(" #%s", tag)
		}
		t += tags

		snippetTexts[t] = s
		if config.Flag.Color {
			t = fmt.Sprintf("[%s]: %s",
				color.RedString(s.Description), s.Command)
		}
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

	params := dialog.SearchForParams(lines)
	if params != nil {
		snippetInfo := snippetTexts[lines[0]]
		dialog.CurrentCommand = snippetInfo.Command
		dialog.GenerateParamsLayout(params, dialog.CurrentCommand)
		res := []string{dialog.FinalCommand}
		return res, nil
	}
	for _, line := range lines {
		snippetInfo := snippetTexts[line]
		commands = append(commands, fmt.Sprint(snippetInfo.Command))
	}
	return commands, nil
}
