package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/dialog"
	"github.com/knqyf263/pet/snippet"
)

func editFile(command, file string) error {
	command += " " + file
	return run(command, os.Stdin, os.Stdout)
}

func run(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if len(config.Conf.General.Cmd) > 0 {
		line := append(config.Conf.General.Cmd, command)
		cmd = exec.Command(line[0], line[1:]...)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func filter(options []string, tag string) (commands []string, err error) {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return commands, fmt.Errorf("Load snippet failed: %v", err)
	}

	if len(tag) > 0 {
		var filteredSnippets snippet.Snippets
		for _, snippet := range snippets.Snippets {
			for _, t := range snippet.Tag {
				if tag == t {
					filteredSnippets.Snippets = append(filteredSnippets.Snippets, snippet)
				}
			}
		}
		snippets = filteredSnippets
	}

	// This is a map of section headings to snippet, so that if a heading is picked, we can use the first command as the default
	snippetHeadings := map[string]snippet.SnippetInfo{}
	// This is a map of formatted commands to original commands, so that we can format them however we want
	commandTexts := map[string]string{}
	var text string
	for _, s := range snippets.Snippets {
		var sectionHeader string
		formattedCommands := make([]string, 0)
		// format commands
		for i, command := range s.Commands {
			if strings.ContainsAny(command, "\n") {
				command = strings.Replace(command, "\n", "\\n", -1)
			}
			command = fmt.Sprintf("$ %s", command)
			formattedCommands = append(formattedCommands, command)
			commandTexts[command] = s.Commands[i]
		}

		// format tags
		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf(" #%s", tag)
		}

		// section heading
		headerKey := fmt.Sprintf("[%s] %s", s.Description, tags)
		if config.Flag.Color {
			sectionHeader = fmt.Sprintf("[%s] %s", color.RedString(s.Description), color.BlueString(tags))
		} else {
			sectionHeader = headerKey
		}
		// associate the top-level description with a snippet so we can default to the first option if picked
		// clip in fzf always copies the plaintext version, so we need to store the non-color version as our key
		// TODO see if we can make this more robust with other fuzzy searches
		snippetHeadings[headerKey] = s
		if len(text) > 0 {
			text += "\n"
		}
		text += fmt.Sprintf("%s\n", sectionHeader)

		// add the commands
		text += strings.Join(formattedCommands, "\n")
	}

	var buf bytes.Buffer
	selectCmd := fmt.Sprintf("%s %s",
		config.Conf.General.SelectCmd, strings.Join(options, " "))
	err = run(selectCmd, strings.NewReader(text), &buf)
	if err != nil {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")

	for _, line := range lines {
		var command string
		// first see if they selected a snippet heading
		if snippetInfo, ok := snippetHeadings[line]; ok {
			// default to first command
			command = fmt.Sprint(snippetInfo.Commands[0])
		} else if snippetText, ok := commandTexts[line]; ok {
			command = fmt.Sprint(snippetText)
		} else {
			fmt.Fprintf(color.Output, "\n%s: %s\n", color.HiRedString("unable to select command: "), line)
			continue
		}

		// extract params from the original command
		params := dialog.SearchForParams(command)
		if params != nil {
			dialog.CurrentCommand = command
			dialog.GenerateParamsLayout(params, dialog.CurrentCommand)
			res := []string{dialog.FinalCommand}
			return res, nil
		}

		commands = append(commands, command)
	}
	return commands, nil
}
