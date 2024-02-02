package dialog

import (
	"regexp"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var (
	views      = []string{}
	layoutStep = 4
	curView    = -1
	idxView    = 0

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string

	patternRegex = `<([\S]+?)>`
)

func insertParams(command string, params map[string]string) string {
	r, _ := regexp.Compile(patternRegex)

	matches := r.FindAllStringSubmatch(command, -1)
	if len(matches) == 0 {
		return command
	}

	resultCommand := command
	for _, p := range matches {
		splitted := strings.Split(p[1], "=")
		resultCommand = strings.Replace(resultCommand, p[0], params[splitted[0]], -1)
	}

	return resultCommand
}

// SearchForParams returns variables from a command
func SearchForParams(lines []string) map[string]string {
	if len(lines) == 1 {
		r, _ := regexp.Compile(patternRegex)

		params := r.FindAllStringSubmatch(lines[0], -1)
		if len(params) == 0 {
			return nil
		}

		extracted := map[string]string{}
		for _, p := range params {
			splitted := strings.Split(p[1], "=")
			key := splitted[0]
			_, param_exists := extracted[key]

			// Set to empty if no value is provided and param is not already set
			if len(splitted) == 1 && !param_exists {
				extracted[key] = ""
			} else if len(splitted) > 1 {
				extracted[key] = splitted[1]
			}
		}
		return extracted
	}
	return nil
}

func evaluateParams(g *gocui.Gui, _ *gocui.View) error {
	paramsFilled := map[string]string{}
	for _, v := range views {
		view, _ := g.View(v)
		res := view.Buffer()
		res = strings.Replace(res, "\n", "", -1)
		paramsFilled[v] = strings.TrimSpace(res)
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)
	return gocui.ErrQuit
}
