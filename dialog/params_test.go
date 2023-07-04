package dialog

import (
	"reflect"
	"testing"
)

func TestSearchForParams(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"01", args{[]string{"... <domain>"}}, map[string]string{"<domain>": ""}},
		{"02", args{[]string{"... <domain=example.com>"}}, map[string]string{"<domain=example.com>": "example.com"}},
		{"03", args{[]string{"... <domain=example.com> <port=443>"}}, map[string]string{"<domain=example.com>": "example.com", "<port=443>": "443"}},
		{"04", args{[]string{"cat <<EOF > <file=path/to/file>\nEOF"}}, map[string]string{"<file=path/to/file>": "path/to/file"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SearchForParams(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchForParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
