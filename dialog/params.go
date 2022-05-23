package dialog

import (
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	views               = []string{}
	viewsUseMultiValues = map[string]bool{}
	layoutStep          = 3
	curView             = -1
	idxView             = 0

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string
)

func insertParams(command string, params map[string]string) string {
	resultCommand := command
	for k, v := range params {
		resultCommand = strings.Replace(resultCommand, k, v, -1)
	}
	return resultCommand
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
			splitted := strings.Split(p[1], "=")
			if len(splitted) == 1 {
				extracted[p[0]] = ""
			} else {
				extracted[p[0]] = splitted[1]
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
		if _, ok := viewsUseMultiValues[v]; !ok {
			res = strings.Replace(res, "\n", "", -1)
			paramsFilled[v] = strings.TrimSpace(res)
		} else {
			_, cy := view.Cursor()
			l, err := view.Line(cy)
			if err != nil {
				return err
			}
			paramsFilled[v] = strings.TrimSpace(l)
		}
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)
	return gocui.ErrQuit
}
