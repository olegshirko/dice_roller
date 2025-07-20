package ui

import (
	"github.com/sqweek/dialog"
)

var ShowFilePicker = func() ([]string, error) {
	filename, err := dialog.File().Filter("Image files", "png", "jpg", "jpeg").Title("Select a texture").Load()
	if err != nil {
		return nil, err
	}
	return []string{filename}, nil
}
