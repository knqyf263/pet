package dialog

import (
	"testing"

	"github.com/go-test/deep"
)

func TestSearchForParams(t *testing.T) {
	command := "<a=1> <b> hello"

	params := map[string]string{
		"a": "1",
		"b": "",
	}

	got := SearchForParams([]string{command})

	for key, value := range params {
		if got[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}

	for key, value := range got {
		if params[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}
}

func TestSearchForParams_WithNoParams(t *testing.T) {
	command := "no params"

	got := SearchForParams([]string{command})

	if got != nil {
		t.Fatalf("wanted nil, got '%v'", got)
	}
}

func TestSearchForParams_WithMultipleParams(t *testing.T) {
	command := "<a=1> <b> <c=3>"

	params := map[string]string{
		"a": "1",
		"b": "",
		"c": "3",
	}

	got := SearchForParams([]string{command})

	for key, value := range params {
		if got[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}

	for key, value := range got {
		if params[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}
}

func TestSearchForParams_WithEmptyCommand(t *testing.T) {
	command := ""

	got := SearchForParams([]string{command})

	if got != nil {
		t.Fatalf("wanted nil, got '%v'", got)
	}
}

func TestSearchForParams_WithNewline(t *testing.T) {
	command := "<a=1> <b> hello\n<c=3>"

	params := map[string]string{
		"a": "1",
		"b": "",
		"c": "3",
	}

	got := SearchForParams([]string{command})

	for key, value := range params {
		if got[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}

	for key, value := range got {
		if params[key] != value {
			t.Fatalf("wanted param '%s' to equal '%s', got '%s'", key, value, got[key])
		}
	}
}

func TestSearchForParams_InvalidParamFormat(t *testing.T) {
	command := "<a=1 <b> hello"
	want := map[string]string{
		"b": "",
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_ConfusingBrackets(t *testing.T) {
	command := "cat <<EOF > <file=path/to/file>\nEOF"
	want := map[string]string{
		"file": "path/to/file",
	}
	got := SearchForParams([]string{command})
	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKey(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := map[string]string{
		"a": "3",
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := map[string]string{
		"a": "3",
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4>"
	want := map[string]string{
		"a": "3",
		"b": "4",
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat(t *testing.T) {
	command := "<a=1> <a=2 <a=3>"
	want := map[string]string{
		"a": "3",
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3 \n<b=4>"
	want := map[string]string{
		"a": "2",
		"b": "4",
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines2(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4"
	want := map[string]string{
		"a": "3",
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}
