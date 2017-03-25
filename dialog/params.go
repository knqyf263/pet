package dialog

import (
	"github.com/jroimartin/gocui"
	"regexp"
	"strings"
)

var re string = "<(.+?)>"

var (
	Params        = map[string]string{}
	Views         = []string{}
	Layout_step   = 3
	Params_filled = map[string]string{}
	curView       = -1
	idxView       = 0

	Current_command string
	Final_command   string
)

func insertParams(command string, params map[string]string) string {
	var result_command string = command
	for k, v := range params {
		result_command = strings.Replace(result_command, k, v, -1)
	}
	return result_command
}

func SearchForParams(lines []string) map[string]string {
	if len(lines) == 1 {
		first_line_only := lines[0]
		r, _ := regexp.Compile(re)

		params := r.FindAllStringSubmatch(first_line_only, -1)
		if len(params) == 0 {
			return nil
		} else {
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
	}
	return nil
}

func evaluateParams(g *gocui.Gui, _ *gocui.View) error {
	for _, v := range Views {
		view, _ := g.View(v)
		res := view.Buffer()
		res = strings.Replace(res, "\n", "", -1)
		Params_filled[v] = string(res)
	}
	Final_command = insertParams(Current_command, Params_filled)
	return gocui.ErrQuit
}
