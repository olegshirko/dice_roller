//go:build !ci

package ui

import (
	"errors"
	"testing"
)

func TestShowFilePicker(t *testing.T) {
	t.Run("successful file selection", func(t *testing.T) {
		// Mock the dialog function
		originalShowFilePicker := ShowFilePicker
		ShowFilePicker = func() ([]string, error) {
			return []string{"/fake/path/to/image.png"}, nil
		}
		defer func() { ShowFilePicker = originalShowFilePicker }()

		files, err := ShowFilePicker()
		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}
		if len(files) != 1 || files[0] != "/fake/path/to/image.png" {
			t.Errorf("expected a single file path, but got: %v", files)
		}
	})

	t.Run("file selection cancelled", func(t *testing.T) {
		// Mock the dialog function to simulate an error
		originalShowFilePicker := ShowFilePicker
		ShowFilePicker = func() ([]string, error) {
			return nil, errors.New("file selection cancelled")
		}
		defer func() { ShowFilePicker = originalShowFilePicker }()

		_, err := ShowFilePicker()
		if err == nil {
			t.Fatal("expected an error, but got none")
		}
		expectedError := "file selection cancelled"
		if err.Error() != expectedError {
			t.Errorf("expected error message '%s', but got '%s'", expectedError, err.Error())
		}
	})
}
