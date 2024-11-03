package dialog

import (
	"fmt"
	"testing"

	"github.com/awesome-gocui/gocui"
	"github.com/go-test/deep"
)

type mockGui struct {
	setViewFn        func(name string, x0, y0, x1, y1 int, overlaps byte) (*gocui.View, error)
	setCurrentViewFn func(name string) (*gocui.View, error)
}

func (m *mockGui) SetView(name string, x0, y0, x1, y1 int, overlaps byte) (*gocui.View, error) {
	if m.setViewFn != nil {
		return m.setViewFn(name, x0, y0, x1, y1, overlaps)
	}

	return &gocui.View{}, nil
}

func (m *mockGui) SetCurrentView(name string) (*gocui.View, error) {
	if m.setCurrentViewFn != nil {
		return m.setCurrentViewFn(name)
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

func TestNextView(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			name        string
			curView     int
			views       []string
			want        string
			wantCurView int
		}{
			{name: "single view", curView: -1, views: []string{"dummy view"}, want: "dummy view", wantCurView: 0},
			{name: "multiple views", curView: 0, views: []string{"dummy view", "dummy view 2"}, want: "dummy view 2", wantCurView: 1},
			{name: "loop back to first view", curView: 1, views: []string{"dummy view", "dummy view 2"}, want: "dummy view", wantCurView: 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Cleanup(func() {
					views = []string{}
					curView = -1
				})

				views = tt.views
				curView = tt.curView

				want := tt.want
				expectedError := error(nil)

				g := &mockGui{
					setCurrentViewFn: func(name string) (*gocui.View, error) {
						got := name

						if diff := deep.Equal(want, got); diff != nil {
							t.Fatal(diff)
						}

						return &gocui.View{}, nil
					},
				}

				err := nextView(g)

				if err != expectedError {
					t.Errorf("Expected error %v, but got %v", expectedError, err)
				}

				if diff := deep.Equal(tt.wantCurView, curView); diff != nil {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("error - returns error if gui returns error", func(t *testing.T) {
		t.Cleanup(func() {
			views = []string{}
		})

		views = []string{"dummy view"}

		wantCurView := -1
		expectedError := fmt.Errorf("dummy error")

		g := &mockGui{
			setCurrentViewFn: func(name string) (*gocui.View, error) {
				return nil, expectedError
			},
		}

		err := nextView(g)

		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if diff := deep.Equal(wantCurView, curView); diff != nil {
			t.Fatal(diff)
		}
	})
}
