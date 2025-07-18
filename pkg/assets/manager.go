package assets

import (
	"github.com/olegshirko/dice_roller/internal/utils"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/cube"
	"github.com/olegshirko/dice_roller/pkg/ui"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
)

// Manager управляет всеми ресурсами (текстурами) в игре.
type Manager struct {
	AllTextures       []*ebiten.Image // Все когда-либо загруженные текстуры
	AvailableTextures []*ebiten.Image // Текстуры, доступные для использования
}

// NewManager создает новый менеджер ассетов.
func NewManager() *Manager {
	return &Manager{
		AllTextures:       []*ebiten.Image{},
		AvailableTextures: []*ebiten.Image{},
	}
}

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
		if tex := loadTextureFromFile(filename); tex != nil {
			m.AllTextures = append(m.AllTextures, tex)
		}
	}

	if len(m.AllTextures) > 0 {
		m.prepareAvailableTextures()
	}
}

// LoadFromDirectory загружает все изображения из указанной директории.
func (m *Manager) LoadFromDirectory(dir string) bool {
	log.Printf("Attempting to auto-load textures from directory '%s'", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Directory '%s' not found. Skipping auto-load.", dir)
		} else {
			log.Printf("Could not read directory %s: %v. Skipping auto-load.", dir, err)
		}
		return false
	}

	// Не очищаем текстуры, если они уже были загружены, а добавляем к ним
	// m.AllTextures = nil
	// m.AvailableTextures = nil

	loaded := false
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}

		fullPath := filepath.Join(dir, file.Name())
		if tex := loadTextureFromFile(fullPath); tex != nil {
			m.AllTextures = append(m.AllTextures, tex)
			loaded = true
		}
	}

	if loaded {
		log.Printf("Found and loaded images from '%s'.", dir)
		m.prepareAvailableTextures()
		return true
	}

	log.Printf("No valid images found in directory '%s'.", dir)
	return false
}

// loadTextureFromFile загружает одну текстуру из файла и добавляет на нее метку.
func loadTextureFromFile(path string) *ebiten.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v", path, err)
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding image %s: %v", path, err)
		return nil
	}

	label := filepath.Base(path)
	ext := filepath.Ext(label)
	label = label[:len(label)-len(ext)]

	labeledImg := utils.AddLabelToImage(img, label)
	log.Printf("Loaded and labeled texture from %s", path)
	return labeledImg
}

// prepareAvailableTextures копирует все загруженные текстуры в пул доступных и перемешивает их.
func (m *Manager) prepareAvailableTextures() {
	m.AvailableTextures = make([]*ebiten.Image, len(m.AllTextures))
	copy(m.AvailableTextures, m.AllTextures)
	rand.Shuffle(len(m.AvailableTextures), func(i, j int) {
		m.AvailableTextures[i], m.AvailableTextures[j] = m.AvailableTextures[j], m.AvailableTextures[i]
	})
	log.Printf("Loaded %d textures. Available pool created and shuffled.", len(m.AvailableTextures))
}

// SetInitialTextures устанавливает начальные текстуры на грани куба.
func (m *Manager) SetInitialTextures(faces *[6]cube.Face, isGrey *[6]bool) {
	if len(m.AvailableTextures) == 0 {
		log.Println("No available textures to set on start.")
		for i := 0; i < 6; i++ {
			(*faces)[i].Texture = config.EmptyImage
			(*isGrey)[i] = true
		}
		return
	}

	log.Println("Setting initial textures on cube faces...")
	for i := 0; i < 6; i++ {
		if len(m.AvailableTextures) > 0 {
			newTex := m.AvailableTextures[len(m.AvailableTextures)-1]
			m.AvailableTextures = m.AvailableTextures[:len(m.AvailableTextures)-1]
			(*faces)[i].Texture = newTex
			(*isGrey)[i] = false
		} else {
			(*faces)[i].Texture = config.GreyImage
			(*isGrey)[i] = true
			log.Printf("Available textures ran out. Face %d is set to grey.", i)
		}
	}
	log.Printf("Set initial active textures. %d textures remaining available.", len(m.AvailableTextures))
}

// ReplaceFaceTexture заменяет текстуру на указанной грани.
func (m *Manager) ReplaceFaceTexture(faceIndex int, faces *[6]cube.Face, isGrey *[6]bool) {
	if faceIndex < 0 || faceIndex >= 6 {
		return
	}

	if len(m.AvailableTextures) > 0 {
		newTexture := m.AvailableTextures[len(m.AvailableTextures)-1]
		m.AvailableTextures = m.AvailableTextures[:len(m.AvailableTextures)-1]
		(*faces)[faceIndex].Texture = newTexture
		(*isGrey)[faceIndex] = false
	} else {
		(*isGrey)[faceIndex] = true
	}
}
