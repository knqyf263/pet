package dialog

import (
	"fmt"
	"testing"

	"github.com/awesome-gocui/gocui"
	"github.com/go-test/deep"
)

type mockGui struct {
	setViewFn        func(name string, x0, y0, x1, y1 int, overlaps byte) (*gocui.View, error)
}

func (m *mockGui) SetView(name string, x0, y0, x1, y1 int, overlaps byte) (*gocui.View, error) {
	if m.setViewFn != nil {
		return m.setViewFn(name, x0, y0, x1, y1, overlaps)
	}

	return &gocui.View{}, nil
}

func TestCreateView(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			views = []string{}
		})

		want := &gocui.View{
			Title:      "dummy view",
			Wrap:       true,
			Autoscroll: true,
		}
		wantViews := []string{"dummy view"}
		expectedError := error(nil)

		g := &mockGui{}

		got, err := createView(g, "dummy view", [4]int{}, false)
		gotViews := views

		if diff := deep.Equal(want, got); diff != nil {
			t.Fatal(diff)
		}

		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if diff := deep.Equal(wantViews, gotViews); diff != nil {
			t.Fatal(diff)
		}
	})

	t.Run("success - does nothing if view already exists", func(t *testing.T) {
		t.Cleanup(func() {
			views = []string{}
		})

		views = []string{"dummy view"}

		var want *gocui.View
		wantViews := []string{"dummy view"}
		expectedError := error(nil)

		g := &mockGui{}

		got, err := createView(g, "dummy view", [4]int{}, false)
		gotViews := views

		if diff := deep.Equal(want, got); diff != nil {
			t.Fatal(diff)
		}

		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if diff := deep.Equal(wantViews, gotViews); diff != nil {
			t.Fatal(diff)
		}
	})

	t.Run("error - returns error if gui returns error", func(t *testing.T) {
		t.Cleanup(func() {
			views = []string{}
		})

		var want *gocui.View
		wantViews := []string{}
		expectedError := fmt.Errorf("dummy error")

		g := &mockGui{
			setViewFn: func(name string, x0, y0, x1, y1 int, overlaps byte) (*gocui.View, error) {
				return nil, expectedError
			},
		}

		got, err := createView(g, "dummy view", [4]int{}, false)
		gotViews := views

		if diff := deep.Equal(want, got); diff != nil {
			t.Fatal(diff)
		}

		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if diff := deep.Equal(wantViews, gotViews); diff != nil {
			t.Fatal(diff)
		}
	})
}
