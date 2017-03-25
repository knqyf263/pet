package dialog

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func generateView(g *gocui.Gui, desc string, fill string, coords []int, editable bool) error {
	if StringInSlice(desc, Views) {
		return nil
	}
	if v, err := g.SetView(desc, coords[0], coords[1], coords[2], coords[3]); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, fill)
	}
	view, _ := g.View(desc)
	view.Title = desc
	view.Wrap = true
	view.Autoscroll = true
	view.Editable = editable

	Views = append(Views, desc)
	curView = len(Views) - 1
	idxView += 1

	return nil
}

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
		command, []int{maxX / 10, maxY / 10, (maxX / 2) + (maxX / 3), maxY/10 + 5}, false)
	idx := 0
	for k, v := range params {
		if v == "" {
			generateView(g, k, "", []int{maxX / 10, (maxY / 4) + (idx+1)*Layout_step,
				maxX/10 + 20, (maxY / 4) + 2 + (idx+1)*Layout_step}, true)
			idx += 1
		} else {
			generateView(g, k, v, []int{maxX / 10, (maxY / 4) + (idx+1)*Layout_step,
				maxX/10 + 20, (maxY / 4) + 2 + (idx+1)*Layout_step}, true)
			idx += 1
		}
	}

	initKeybindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func nextView(g *gocui.Gui) error {
	next := curView + 1
	if next > len(Views)-1 {
		next = 0
	}

	if _, err := g.SetCurrentView(Views[next]); err != nil {
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
