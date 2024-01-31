package dialog

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func generateView(g *gocui.Gui, p *parameter, coords []int, editable bool) error {
	desc := p.name
	fill := p.options[0]

	if StringInSlice(desc, views) {
		return nil
	}

	if v, err := g.SetView(desc, coords[0], coords[1], coords[2], coords[3], 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, fill)
	}
	view, _ := g.View(desc)

	if len(p.options) > 1 {
		view.Title = desc + " (...)"
	} else {
		view.Title = desc
	}
	view.Wrap = false
	view.Autoscroll = true
	view.Editable = editable

	views = append(views, desc)

	idxView++

	return nil
}

// GenerateParamsLayout generates CUI to receive params
func GenerateParamsLayout(params map[string][]string, command string) {
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
	generateView(g,
		&parameter{
			name:    "Command(TAB => Select next, ENTER => Execute command, Cursor up/down => change optional parameter):",
			options: []string{command},
		},
		[]int{maxX / 10, maxY / 10, (maxX / 2) + (maxX / 3), maxY/10 + 5},
		false)

	idx := 0
	for _, p := range parameters {
		generateView(g, p,
			[]int{maxX / 10,
				(maxY / 4) + (idx+1)*layoutStep,
				(maxX / 2) + (maxX / 3),
				(maxY / 4) + 2 + (idx+1)*layoutStep},
			true)
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

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, updateOptionInViewUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, updateOptionInViewDown); err != nil {
		return err
	}

	return nil
}

func updateOptionInView(g *gocui.Gui, ch int) error {
	// Don't handle up/down key for Command view
	if curView == 0 {
		return nil
	}

	// Shitfting one view as Command view is not on the parameters list
	p := parameters[curView-1]
	if len(p.options) < 2 {
		return nil
	}

	// If ch(ange) is -1 --> key down
	//                +1 --> key up
	if ch == -1 {
		p.current--
		if p.current < 0 {
			p.current = len(p.options) - 1
		}

	} else if ch == 1 {
		p.current++
		if p.current > len(p.options)-1 {
			p.current = 0
		}
	}

	view, _ := g.View(views[curView])
	view.Clear()
	view.Write([]byte(p.options[p.current]))
	return nil
}

func updateOptionInViewUp(g *gocui.Gui, _ *gocui.View) error {
	return updateOptionInView(g, 1)
}

func updateOptionInViewDown(g *gocui.Gui, _ *gocui.View) error {
	return updateOptionInView(g, -1)
}

func layout(g *gocui.Gui) error {
	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
