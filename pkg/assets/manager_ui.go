package assets

import (
	"github.com/olegshirko/dice_roller/pkg/ui"
	"github.com/sqweek/dialog"
	"log"
)

// LoadTextures предлагает пользователю выбрать один или несколько файлов текстур.
func (m *Manager) LoadTextures() {
	log.Println("Opening file dialog to select textures...")
	filenames, err := ui.ShowFilePicker()
	if err != nil {
		if err == dialog.ErrCancelled {
			log.Println("Texture selection cancelled.")
		} else {
			log.Printf("Error selecting file(s): %v", err)
		}
		return
	}

	if len(filenames) == 0 {
		log.Println("No files were selected.")
		return
	}

	// Очищаем старые текстуры только если выбраны новые
	m.AllTextures = nil
	m.AvailableTextures = nil

	for _, filename := range filenames {
		if tex := m.loader.Load(filename); tex != nil {
			m.AllTextures = append(m.AllTextures, tex)
		}
	}

	if len(m.AllTextures) > 0 {
		m.prepareAvailableTextures()
	}
}