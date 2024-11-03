package snippet

import (
	"testing"

	"github.com/go-test/deep"
)

func TestTomlRegEx(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "match - filename", path: "snippet.toml", want: true},
		{name: "match - absolute path", path: "/home/username/.config/pet/config.toml", want: true},
		{name: "match - absolute path with home alias", path: "~/.config/pet/snippet2.toml", want: true},
		{name: "match - relative path", path: "../../some/directory/best.toml", want: true},
		{name: "mismatch - filename", path: "file.yaml", want: false},
		{name: "mismatch - absolute path", path: "/home/username/.config/pet/config.json", want: false},
		{name: "mismatch - absolute path with home alias", path: "~/.config/pet/snippet2.xml", want: false},
		{name: "mismatch - relative path", path: "../../some/directory/unrelated.html", want: false},
		{name: "mismatch - extension with dot", path: ".toml", want: false},
		{name: "mismatch - extension without dot", path: "toml", want: false},
		{name: "mismatch - similar extension prefix", path: "file.tomlx", want: false},
		{name: "mismatch - similar extension suffix", path: "file.xtoml", want: false},
		{name: "mismatch - similar extension infix", path: "file.xtomlx", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tomlRegEx.MatchString(tt.path)

			if got != tt.want {
				t.Errorf("Expected result %v, but got %v", tt.want, got)
			}
		})
	}
}

func TestGetFiles(t *testing.T) {
	t.Run("success - returns list of toml files", func(t *testing.T) {
		got := getFiles("testdata")
		want := []string{
			"testdata/01-snippet.toml",
			"testdata/03-snippet.toml",
			"testdata/04-subdir/05-snippet.toml",
		}

		if diff := deep.Equal(want, got); diff != nil {
			t.Fatal(diff)
		}
	})
}
