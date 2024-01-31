package dialog

import (
	"log"
	"regexp"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var (
	views      = []string{}
	parameters = []*parameter{}
	layoutStep = 3
	curView    = -1
	idxView    = 0

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string
)

type parameter struct {
	original string
	name     string
	options  []string
	current  int
}

func insertParams(command string, params map[string]string) string {
	resultCommand := command
	log.Println("in command ", command)
	for _, v := range parameters {
		resultCommand = strings.Replace(resultCommand, v.original, params[v.name], -1)
	}
	log.Println("out command ", resultCommand)
	return resultCommand
}

// SearchForParams returns variables from a command
func SearchForParams(lines []string) map[string][]string {
	re := `<([\S]+?)>`
	if len(lines) == 1 {
		r, _ := regexp.Compile(re)

		params := r.FindAllStringSubmatch(lines[0], -1)
		if len(params) == 0 {
			return nil
		}

		extracted := map[string][]string{}

		for _, p := range params {
			splitted := strings.Split(p[1], "=")
			key := splitted[0]

			// _, param_exists := extracted[key]

			// Set to empty if no value is provided and param is not already set
			// if len(splitted) == 1 && !param_exists {
			if len(splitted) == 1 {
				p := &parameter{
					original: "<" + p[1] + ">",
					name:     key,
					options:  []string{""},
					current:  0,
				}
				parameters = append(parameters, p)
			} else if len(splitted) > 1 {
				// From a list of parameters (divided with "|", get all of them
				pSplit := strings.Split(splitted[1], "|")
				p := &parameter{
					original: "<" + p[1] + ">",
					name:     key,
					options:  pSplit,
					current:  0,
				}
				parameters = append(parameters, p)
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
		if !strings.Contains(v, "Command") {
			res := view.Buffer()
			res = strings.Replace(res, "\n", "", -1)
			paramsFilled[v] = strings.TrimSpace(res)
		}
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)
	return gocui.ErrQuit
}
