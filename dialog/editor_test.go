package dialog

import (
	"testing"
)

func TestWordOffset(t *testing.T) {
	tests := []struct {
		name string
		line string
		cx   int
		dir  int
		want int
	}{
		{
			name: "backward from end of a single word",
			line: "gamma",
			cx:   5,
			dir:  -1,
			want: -5,
		},
		{
			name: "backward stops at the start of the previous word",
			line: "alpha beta gamma",
			cx:   16,
			dir:  -1,
			want: -5,
		},
		{
			name: "backward skips the separators before the word",
			line: "alpha beta   ",
			cx:   13,
			dir:  -1,
			want: -7,
		},
		{
			name: "backward at the start of the line does not move",
			line: "alpha",
			cx:   0,
			dir:  -1,
			want: 0,
		},
		{
			name: "forward stops at the end of the next word",
			line: "alpha beta gamma",
			cx:   0,
			dir:  1,
			want: 5,
		},
		{
			name: "forward skips the separators before the word",
			line: "alpha   beta",
			cx:   5,
			dir:  1,
			want: 7,
		},
		{
			name: "forward at the end of the line does not move",
			line: "alpha",
			cx:   5,
			dir:  1,
			want: 0,
		},
		{
			name: "punctuation splits words",
			line: "antonio-gemini-api-key",
			cx:   22,
			dir:  -1,
			want: -3,
		},
		{
			name: "counts runes, not bytes",
			line: "olá mundo",
			cx:   9,
			dir:  -1,
			want: -5,
		},
		{
			name: "cursor past the end of the line is clamped",
			line: "alpha",
			cx:   99,
			dir:  -1,
			want: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wordOffset([]rune(tt.line), tt.cx, tt.dir)
			if got != tt.want {
				t.Errorf("wordOffset(%q, %d, %d) = %d, want %d",
					tt.line, tt.cx, tt.dir, got, tt.want)
			}
		})
	}
}
