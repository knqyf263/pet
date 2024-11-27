package dialog

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name string
		item string
		list []string
		want bool
	}{
		{name: "empty list", item: "a", list: []string{}, want: false},
		{name: "not in list", item: "a", list: []string{"b", "c"}, want: false},
		{name: "in list", item: "a", list: []string{"a", "b", "c"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.item, tt.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
