//go:build !ci

package assets

import (
	"github.com/olegshirko/dice_roller/pkg/cube"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)


// TestNewManager проверяет конструктор NewManager.
func TestNewManager(t *testing.T) {
	m := NewManager()
	assert.NotNil(t, m, "NewManager should not return nil")
	assert.NotNil(t, m.AllTextures, "AllTextures should be initialized")
	assert.Empty(t, m.AllTextures, "AllTextures should be empty")
	assert.NotNil(t, m.AvailableTextures, "AvailableTextures should be initialized")
	assert.Empty(t, m.AvailableTextures, "AvailableTextures should be empty")
	assert.NotNil(t, m.loader, "loader should be initialized")
	_, ok := m.loader.(*ebitenTextureLoader)
	assert.True(t, ok, "loader should be of type ebitenTextureLoader")
}

// TestEbitenTextureLoader_Load проверяет реальную загрузку текстуры.
func TestEbitenTextureLoader_Load(t *testing.T) {
	// Создаем временный PNG файл
	tmpFile, err := os.CreateTemp("", "test_*.png")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	err = png.Encode(tmpFile, img)
	assert.NoError(t, err)
	tmpFile.Close()

	loader := &ebitenTextureLoader{}
	ebitenImg := loader.Load(tmpFile.Name())
	assert.NotNil(t, ebitenImg, "Load should return a non-nil image for a valid file")
}

// TestLoadTextureFromFile_Success проверяет успешную загрузку из файла.
func TestLoadTextureFromFile_Success(t *testing.T) {
	// Создаем временный PNG файл
	tmpFile, err := os.CreateTemp("", "test_*.png")
	assert.NoError(t, err)
	path := tmpFile.Name()
	defer os.Remove(path)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.Black)
	err = png.Encode(tmpFile, img)
	assert.NoError(t, err)
	tmpFile.Close()

	ebitenImg := loadTextureFromFile(path)
	assert.NotNil(t, ebitenImg, "Should successfully load a valid image file")
}

// TestLoadTextureFromFile_FileNotExist проверяет случай, когда файл не существует.
func TestLoadTextureFromFile_FileNotExist(t *testing.T) {
	ebitenImg := loadTextureFromFile("non_existent_file.png")
	assert.Nil(t, ebitenImg, "Should return nil for a non-existent file")
}

// TestLoadTextureFromFile_InvalidImage проверяет обработку некорректного формата изображения.
func TestLoadTextureFromFile_InvalidImage(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid_*.txt")
	assert.NoError(t, err)
	path := tmpFile.Name()
	defer os.Remove(path)

	_, err = tmpFile.WriteString("this is not an image")
	assert.NoError(t, err)
	tmpFile.Close()

	ebitenImg := loadTextureFromFile(path)
	assert.Nil(t, ebitenImg, "Should return nil for a file that is not a valid image")
}

// TestSetInitialTextures_WithAvailableTextures проверяет установку начальных текстур.
func TestSetInitialTextures_WithAvailableTextures(t *testing.T) {
	m := NewManager()
	tex1 := ebiten.NewImage(1, 1)
	tex2 := ebiten.NewImage(1, 1)
	m.AvailableTextures = []*ebiten.Image{tex1, tex2}

	var faces [6]cube.Face
	var isGrey [6]bool

	m.SetInitialTextures(&faces, &isGrey)

	assert.Equal(t, tex2, faces[0].Texture, "Face 0 should have the last available texture")
	assert.Equal(t, tex1, faces[1].Texture, "Face 1 should have the second to last available texture")
	assert.False(t, isGrey[0], "Face 0 should not be grey")
	assert.False(t, isGrey[1], "Face 1 should not be grey")
	assert.Empty(t, m.AvailableTextures, "AvailableTextures should be empty after setting initial textures")
}

// TestSetInitialTextures_NoAvailableTextures проверяет установку, когда нет доступных текстур.
func TestSetInitialTextures_NoAvailableTextures(t *testing.T) {
	m := NewManager()
	var faces [6]cube.Face
	var isGrey [6]bool

	m.SetInitialTextures(&faces, &isGrey)

	for i := 0; i < 6; i++ {
		assert.NotNil(t, faces[i].Texture, "Face %d should have a texture", i)
		assert.True(t, isGrey[i], "Face %d should be grey", i)
	}
}

// TestReplaceFaceTexture_WithAvailableTextures проверяет замену текстуры грани.
func TestReplaceFaceTexture_WithAvailableTextures(t *testing.T) {
	m := NewManager()
	newTex := ebiten.NewImage(1, 1)
	m.AvailableTextures = []*ebiten.Image{newTex}

	var faces [6]cube.Face
	var isGrey [6]bool
	faceIndex := 2

	m.ReplaceFaceTexture(faceIndex, &faces, &isGrey)

	assert.Equal(t, newTex, faces[faceIndex].Texture, "Face texture should be replaced")
	assert.False(t, isGrey[faceIndex], "Face should not be grey after replacement")
	assert.Empty(t, m.AvailableTextures, "AvailableTextures should be empty after replacement")
}

// TestReplaceFaceTexture_NoAvailableTextures проверяет замену, когда нет доступных текстур.
func TestReplaceFaceTexture_NoAvailableTextures(t *testing.T) {
	m := NewManager()
	var faces [6]cube.Face
	var isGrey [6]bool
	faceIndex := 3
	faces[faceIndex].Texture = ebiten.NewImage(1, 1) // Изначальная текстура
	isGrey[faceIndex] = false

	m.ReplaceFaceTexture(faceIndex, &faces, &isGrey)

	assert.NotNil(t, faces[faceIndex].Texture, "Face texture should not be nil")
	assert.True(t, isGrey[faceIndex], "Face should become grey")
}

// TestReplaceFaceTexture_InvalidIndex проверяет замену с неверным индексом.
func TestReplaceFaceTexture_InvalidIndex(t *testing.T) {
	m := NewManager()
	var faces [6]cube.Face
	var isGrey [6]bool
	originalFaces := faces
	originalIsGrey := isGrey

	m.ReplaceFaceTexture(10, &faces, &isGrey) // Неверный индекс

	assert.Equal(t, originalFaces, faces, "Faces should not change for invalid index")
	assert.Equal(t, originalIsGrey, isGrey, "isGrey should not change for invalid index")
}
