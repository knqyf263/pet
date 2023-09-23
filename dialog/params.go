package dialog

import (
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	views      = []string{}
	layoutStep = 3
	curView    = -1
	idxView    = 0

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string
)

func insertParams(command string, params map[string]string) string {
	resultCommand := normalizeCommand(command)
	for k, v := range params {
		resultCommand = strings.Replace(resultCommand, normalizeParam(k), v, -1)
	}
	return resultCommand
}

// strips param of default value
func normalizeParam(param string) string {
    splits := strings.SplitN(param, "=", 2)

    if len(splits) == 1 {
        return splits[0]
    }

    return splits[0] + ">"
}

// returns the command with all the params stripped of default values
func normalizeCommand(command string) string {
    r, _ := regexp.Compile(`<([\S].+?[\S])>`)

    params := r.FindAllString(command, -1)

    for _, p := range params {
        command = strings.Replace(command, p, normalizeParam(p), -1)
    }

    return command
}

// SearchForParams returns variables from a command
func SearchForParams(lines []string) map[string]string {
	re := `<([\S].+?[\S])>`
	if len(lines) == 1 {
		r, _ := regexp.Compile(re)

		params := r.FindAllStringSubmatch(lines[0], -1)
		if len(params) == 0 {
			return nil
		}

		extracted := map[string]string{}
		for _, p := range params {
            normalizedP := normalizeParam(p[0])
			splitted := strings.Split(p[1], "=")

            found := false
            for key := range extracted {
                if normalizeParam(key) == normalizedP {
                    found = true
                    break
                }
            }

            if !found {
                if len(splitted) == 1 {
                    extracted[p[0]] = ""
                } else {
                    extracted[p[0]] = splitted[1]
                }
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
