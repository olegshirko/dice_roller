package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// mockTextureLoader is a mock implementation of the textureLoader interface for testing.
type mockTextureLoader struct {
	// failOnLoad can be used to simulate errors during texture loading.
	failOnLoad bool
}

// Load implements the textureLoader interface for the mock.
// It returns a dummy ebiten.Image for .png files and nil for others,
// or if failOnLoad is true.
func (m *mockTextureLoader) Load(path string) *ebiten.Image {
	if m.failOnLoad {
		return nil
	}
	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		// Return a non-nil dummy image to simulate successful loading.
		// The image itself doesn't need to be valid for this test.
		return ebiten.NewImage(1, 1)
	}
	return nil
}

// newTestManager creates a Manager with a mock loader for testing.
func newTestManager(mockLoader textureLoader) *Manager {
	return &Manager{
		AllTextures:       []*ebiten.Image{},
		AvailableTextures: []*ebiten.Image{},
		loader:            mockLoader,
	}
}

func TestLoadFromDirectory_Success(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-assets-success")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Создаем 6 валидных файлов изображений (пустых)
	for i := 1; i <= 6; i++ {
		filePath := filepath.Join(tmpDir, fmt.Sprintf("face%d.png", i))
		if _, err := os.Create(filePath); err != nil {
			t.Fatalf("Failed to create temp file %s: %v", filePath, err)
		}
	}

	// Создаем один не-изображение файл, который должен быть проигнорирован
	nonImagePath := filepath.Join(tmpDir, "document.txt")
	if _, err := os.Create(nonImagePath); err != nil {
		t.Fatalf("Failed to create non-image file: %v", err)
	}

	m := newTestManager(&mockTextureLoader{})
	loaded := m.LoadFromDirectory(tmpDir)

	if !loaded {
		t.Error("LoadFromDirectory() returned false, want true")
	}

	if len(m.AllTextures) != 6 {
		t.Errorf("Expected 6 textures to be loaded, but got %d", len(m.AllTextures))
	}
}

func TestLoadFromDirectory_DirNotFound(t *testing.T) {
	m := newTestManager(&mockTextureLoader{})
	// Пытаемся загрузить из несуществующей директории
	loaded := m.LoadFromDirectory("non-existent-dir")

	if loaded {
		t.Error("LoadFromDirectory() returned true for a non-existent directory, want false")
	}

	if len(m.AllTextures) != 0 {
		t.Errorf("Expected 0 textures, but got %d", len(m.AllTextures))
	}
}

func TestLoadFromDirectory_NoImages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-assets-no-images")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Создаем только не-изображение файлы
	filePath := filepath.Join(tmpDir, "data.json")
	if _, err := os.Create(filePath); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	m := newTestManager(&mockTextureLoader{})
	loaded := m.LoadFromDirectory(tmpDir)

	if loaded {
		t.Error("LoadFromDirectory() returned true when no images are present, want false")
	}

	if len(m.AllTextures) != 0 {
		t.Errorf("Expected 0 textures, but got %d", len(m.AllTextures))
	}
}