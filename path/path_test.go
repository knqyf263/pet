package path

import "testing"

func TestExpandPathWithEmptyInputErrors(t *testing.T) {
	test_path := ""

	_, err := expandPath(test_path)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}
