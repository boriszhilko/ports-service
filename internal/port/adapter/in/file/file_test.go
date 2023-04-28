package file

import (
	"os"
	"testing"
)

func TestNewInput(t *testing.T) {
	filePath := "test.txt"
	os.Create(filePath)
	defer os.Remove(filePath)

	input, err := NewInput(filePath)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if input == nil {
		t.Fatalf("Expected a non-nil Input, but got nil")
	}

	if input.GetData() == nil {
		t.Fatalf("Expected a non-nil io.Reader, but got nil")
	}

	if err := input.Stop(); err != nil {
		t.Fatalf("Expected no error from Stop, but got %v", err)
	}
}
