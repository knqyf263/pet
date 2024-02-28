package dialog

import (
	"regexp"
	"strings"

	"github.com/awesome-gocui/gocui"
)

var (
	views = []string{}

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string

	// This matches most encountered patterns
	// Skips match if there is a whitespace at the end ex. <param='my >
	// Ignores <, > characters since they're used to match the pattern
	patternRegex = `<([^<>]*[^\s])>`
)

func insertParams(command string, filledInParams map[string]string) string {
	r := regexp.MustCompile(patternRegex)

	matches := r.FindAllStringSubmatch(command, -1)
	if len(matches) == 0 {
		return command
	}

	resultCommand := command

	// First match is the whole match (with brackets), second is the first group
	// Ex. echo <param='my param'>
	// -> matches[0][0]: <param='my param'>
	// -> matches[0][1]: param='my param'
	for _, p := range matches {
		whole, matchedGroup := p[0], p[1]
		param, _, _ := strings.Cut(matchedGroup, "=")

		// Replace the whole match with the filled-in value of the param
		resultCommand = strings.Replace(resultCommand, whole, filledInParams[param], -1)
	}

	return resultCommand
}

// SearchForParams returns variables from a command
func SearchForParams(command string) [][2]string {
	r := regexp.MustCompile(patternRegex)

	params := r.FindAllStringSubmatch(command, -1)
	if len(params) == 0 {
		return nil
	}

	extracted := map[string]string{}
	ordered_params := [][2]string{}
	for _, p := range params {
		_, matchedGroup := p[0], p[1]
		paramKey, defaultValue, separatorFound := strings.Cut(matchedGroup, "=")
		_, param_exists := extracted[paramKey]

		// Set to empty if no value is provided and param is not already set
		if !separatorFound && !param_exists {
			extracted[paramKey] = ""
		} else if separatorFound {
			// Set to default value instead if it is provided
			extracted[paramKey] = defaultValue
		}

		// Fill in the keys only if seen for the first time to track order
		if !param_exists {
			ordered_params = append(ordered_params, [2]string{paramKey, ""})
		}
	}

	// Fill in the values
	for i, param := range ordered_params {
		pair := [2]string{param[0], extracted[param[0]]}
		ordered_params[i] = pair
	}
	return ordered_params
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
