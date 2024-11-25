package path

import "testing"

func TestExpandAbsolutePathDoesNothing(t *testing.T) {
	test_path := "/var/tmp/"
	want := "/var/tmp/"

	got, err := expandPath(test_path)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandPathWithEmptyInputErrors(t *testing.T) {
	test_path := ""

	_, err := expandPath(test_path)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}
