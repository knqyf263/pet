package dialog

import (
	"fmt"
	"log"
	"regexp"

	"github.com/awesome-gocui/gocui"
)

var (
	layoutStep = 3
	curView    = -1

	// This is for matching multiple default values in parameters
	parameterMultipleValueRegex = `(\|_.*?_\|)`
)

func generateSingleParameterView(g *gocui.Gui, desc string, defaultParam string, coords []int, editable bool) error {
	if StringInSlice(desc, views) {
		return nil
	}
	if v, err := g.SetView(desc, coords[0], coords[1], coords[2], coords[3], 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, defaultParam)
	}
	view, _ := g.View(desc)
	view.Title = desc
	view.Wrap = true
	view.Autoscroll = true
	view.Editable = editable

	views = append(views, desc)

	return nil
}

func generateMultipleParameterView(g *gocui.Gui, desc string, defaultParams []string, coords []int, editable bool) error {
	if StringInSlice(desc, views) {
		return nil
	}

	currentOpt := 0
	maxOpt := len(defaultParams)

	if v, err := g.SetView(desc, coords[0], coords[1], coords[2], coords[3], 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetKeybinding(v.Name(), gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			if maxOpt == 0 {
				return nil
			}
			next := currentOpt + 1
			if next >= maxOpt {
				next = currentOpt
			}
			v.Clear()
			fmt.Fprint(v, defaultParams[next])
			currentOpt = next
			return nil
		})
		g.SetKeybinding(v.Name(), gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			if maxOpt == 0 {
				return nil
			}
			prev := currentOpt - 1
			if prev < 0 {
				prev = currentOpt
			}
			v.Clear()
			fmt.Fprint(v, defaultParams[prev])
			currentOpt = prev
			return nil
		})
		g.SetKeybinding(v.Name(), gocui.KeyCtrlK, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			v.Clear()
			return nil
		})

		fmt.Fprint(v, defaultParams[currentOpt])
	}
	
	view, _ := g.View(desc)
	viewTitle := desc
	// Adjust view title to hint the user about the available 
	// options if there are more than one
	if maxOpt > 1 {
		viewTitle = desc + " (UP/DOWN => Select default value)"
	}
	
	view.Title = viewTitle
	view.Wrap = false
	view.Autoscroll = true
	view.Editable = editable
	views = append(views, desc)
	return nil
}

// GenerateParamsLayout generates CUI to receive params
func GenerateParamsLayout(params [][2]string, command string) {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	maxX, maxY := g.Size()
	leftX := (maxX / 2) - (maxX / 3)
	rightX := (maxX / 2) + (maxX / 3)

	generateSingleParameterView(g, "Command(TAB => Select next, ENTER => Execute command):",
		command, []int{leftX, maxY / 10, rightX, maxY/10 + 5}, false)
	idx := 0

	// Create a view for each param
	for _, pair := range params {
		// Unpack parameter key and value
		parameterKey, parameterValue := pair[0], pair[1]

		// Check value for multiple defaults
		r := regexp.MustCompile(parameterMultipleValueRegex)
		matches := r.FindAllStringSubmatch(parameterValue, -1)

		if len(matches) > 0 {
			// Extract the default values and generate multiple params view
			parameters := []string{}
			for _, p := range matches {
				_, matchedGroup := p[0], p[1]
				// Remove the separators
				matchedGroup = matchedGroup[2 : len(matchedGroup)-2]
				parameters = append(parameters, matchedGroup)
			}
			generateMultipleParameterView(
				g, parameterKey, parameters, []int{
					leftX,
					(maxY / 4) + (idx+1)*layoutStep,
					rightX,
					(maxY / 4) + 2 + (idx+1)*layoutStep},
				true)
		} else {
			// Generate single param view using the single value
			generateSingleParameterView(g, parameterKey, parameterValue,
				[]int{
					leftX,
					(maxY / 4) + (idx+1)*layoutStep,
					rightX,
					(maxY / 4) + 2 + (idx+1)*layoutStep},
				true)
		}
		idx++
	}

	initKeybindings(g)

	curView = 0
	if idx > 0 {
		curView = 1
	}
	g.SetCurrentView(views[curView])

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func nextView(g *gocui.Gui) error {
	next := curView + 1
	if next > len(views)-1 {
		next = 0
	}

	if _, err := g.SetCurrentView(views[next]); err != nil {
		return err
	}

	curView = next
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, evaluateParams); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextView(g)
		}); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
