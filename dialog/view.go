package dialog

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

func generateView(g *gocui.Gui, desc string, defaultParams []string, coords []int, editable bool) error {
	if StringInSlice(desc, views) {
		return nil
	}

	currentOpt := 0
	maxOpt := len(defaultParams)

	if v, err := g.SetView(desc, coords[0], coords[1], coords[2], coords[3]); err != nil {
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
	view.Title = desc
	view.Wrap = false
	view.Autoscroll = true
	view.Editable = editable

	views = append(views, desc)
	idxView++

	return nil
}

// GenerateParamsLayout generates CUI to receive params
func GenerateParamsLayout(params map[string]string, command string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	maxX, maxY := g.Size()
	generateView(g, "Command(TAB => Select next, ENTER => Execute command):",
		[]string{command}, []int{maxX / 10, maxY / 10, (maxX / 2) + (maxX / 3), maxY/10 + 5}, false)
	idx := 0
	for k, v := range params {
		defaultParams := strings.Split(v, "|")
		generateView(g, k, defaultParams, []int{maxX / 10, (maxY / 4) + (idx+1)*layoutStep,
			maxX/10 + 20, (maxY / 4) + 2 + (idx+1)*layoutStep}, true)
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
