package dialog

import (
	"github.com/awesome-gocui/gocui"
)

// paramEditor edits a parameter field. It matches gocui's default editor except
// for word-wise motions, which the default editor does not implement: terminals
// encode Alt+Left/Alt+Right as ESC+b/ESC+f, and gocui's simpleEditor falls
// through to writing the bare rune for any Alt-modified key, so pressing those
// keys inserted a literal "b"/"f" into the field.
func paramEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if mod == gocui.ModAlt {
		switch {
		case ch == 'b', key == gocui.KeyArrowLeft:
			v.MoveCursor(cursorWordOffset(v, -1), 0)
			return
		case ch == 'f', key == gocui.KeyArrowRight:
			v.MoveCursor(cursorWordOffset(v, 1), 0)
			return
		}
	}

	gocui.DefaultEditor.Edit(v, key, ch, mod)
}

// cursorWordOffset resolves the current line and cursor column into the offset
// wordOffset expects. The view's cursor column counts runes, not bytes.
func cursorWordOffset(v *gocui.View, dir int) int {
	cx, cy := v.Cursor()

	line, err := v.Line(cy)
	if err != nil {
		return 0
	}

	return wordOffset([]rune(line), cx, dir)
}

// wordOffset returns how far the cursor has to move along line to reach the next
// word boundary in the direction dir (-1 backward, 1 forward): skip separators,
// then skip the word itself. This mirrors readline's backward-word/forward-word,
// which is the convention the keys being handled here come from.
func wordOffset(line []rune, cx, dir int) int {
	if cx > len(line) {
		cx = len(line)
	}
	if cx < 0 {
		cx = 0
	}

	i := cx
	if dir < 0 {
		for i > 0 && !isWordRune(line[i-1]) {
			i--
		}
		for i > 0 && isWordRune(line[i-1]) {
			i--
		}
	} else {
		for i < len(line) && !isWordRune(line[i]) {
			i++
		}
		for i < len(line) && isWordRune(line[i]) {
			i++
		}
	}

	return i - cx
}

func isWordRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
	case r >= 'A' && r <= 'Z':
	case r >= '0' && r <= '9':
	default:
		return false
	}

	return true
}
