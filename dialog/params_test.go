package dialog

import (
	"testing"

	"github.com/go-test/deep"
)

func TestSearchForParams(t *testing.T) {
	command := "<a=1> <b> hello"

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
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

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
		{"c", "3"},
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
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

	want := [][2]string{
		{"a", "1"},
		{"b", ""},
		{"c", "3"},
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_InvalidParamFormat(t *testing.T) {
	command := "<a=1 <b> hello"
	want := [][2]string{
		{"b", ""},
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_ConfusingBrackets(t *testing.T) {
	command := "cat <<EOF > <file=path/to/file>\nEOF"
	want := [][2]string{
		{"file", "path/to/file"},
	}
	got := SearchForParams([]string{command})
	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKey(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues(t *testing.T) {
	command := "<a=1> <a=2> <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4>"
	want := [][2]string{
		{"a", "3"},
		{"b", "4"},
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat(t *testing.T) {
	command := "<a=1> <a=2 <a=3>"
	want := [][2]string{
		{"a", "3"},
	}
	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines(t *testing.T) {
	command := "<a=1> <a=2> <a=3 \n<b=4>"
	want := [][2]string{
		{"a", "2"},
		{"b", "4"},
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestSearchForParams_MultipleParamsSameKeyDifferentValues_InvalidFormat_MultipleLines2(t *testing.T) {
	command := "<a=1> <a=2> <a=3>\n<b=4"
	want := [][2]string{
		{"a", "3"},
	}

	got := SearchForParams([]string{command})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatal(diff)
	}
}

func TestInsertParams(t *testing.T) {
	command := "<a=1> <a> <b> hello"

	params := map[string]string{
		"a": "test",
		"b": "case",
	}

	got := insertParams(command, params)
	want := "test test case hello"
	if want != got {
		t.Fatalf("wanted '%s', got '%s'", want, got)
	}
}

func TestInsertParams_unique_parameters(t *testing.T) {
	command := "curl -X POST \"<host=http://localhost:9200>/<index>\" -H 'Content-Type: application/json'"

	params := map[string]string{
		"host":  "localhost:9200",
		"index": "test",
	}

	got := insertParams(command, params)
	want := "curl -X POST \"localhost:9200/test\" -H 'Content-Type: application/json'"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInsertParams_complex(t *testing.T) {
	command := "something <host=http://localhost:9200>/<test>/_delete_by_query/<host>"

	params := map[string]string{
		"host": "localhost:9200",
		"test": "case",
	}

	got := insertParams(command, params)
	want := "something localhost:9200/case/_delete_by_query/localhost:9200"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
