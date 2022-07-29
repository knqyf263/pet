package dialog

import "testing"

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
