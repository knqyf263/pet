package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// MockReadCloser is a mock implementation of io.ReadCloser
type MockReadCloser struct {
	*strings.Reader
}

// Close does nothing for this mock implementation
func (m *MockReadCloser) Close() error {
	return nil
}

func TestScan(t *testing.T) {
	message := "Enter something: "

	input := "test\n" // Simulated user input
	want := "test"    // Expected output
	expectedError := error(nil)

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, false)

	// Check if the input was printed
	got := result

	// Check if the result matches the expected result
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err != expectedError {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestScan_EmptyStringWithAllowEmpty(t *testing.T) {
	message := "Enter something: "

	input := "\n"               // Simulated user input
	want := ""                  // Expected output
	expectedError := error(nil) // Should not error

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, true)

	// Check if the input was printed
	got := result

	// Check if the result is empty
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err != expectedError {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestScan_EmptyStringWithoutAllowEmpty(t *testing.T) {
	message := "Enter something: "

	input := "\n" // Simulated user input
	want := ""    // Expected output
	expectedError := CanceledError()

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, false)

	// Check if the input was printed
	got := result

	// Check if the result matches the expected result
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}
