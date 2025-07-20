//go:build !ci

package assets

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olegshirko/dice_roller/pkg/ui"
	"github.com/sqweek/dialog"
	"github.com/stretchr/testify/assert"
)

// TestLoadTextures_Success проверяет успешный сценарий загрузки текстур.
func TestLoadTextures_Success(t *testing.T) {
	// Подменяем функцию выбора файлов
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return []string{"/fake/path/to/texture1.png", "/fake/path/to/texture2.png"}, nil
	}

	// Создаем менеджер с мок-загрузчиком
	manager := newTestManager(&mockTextureLoader{})

	// Вызываем тестируемую функцию
	manager.LoadTextures()

	// Проверяем результат
	assert.Len(t, manager.AllTextures, 2, "Должно быть загружено 2 текстуры")
	assert.Len(t, manager.AvailableTextures, 2, "Должно быть 2 доступные текстуры")
}

// TestLoadTextures_Cancelled проверяет сценарий отмены выбора файла.
func TestLoadTextures_Cancelled(t *testing.T) {
	// Подменяем функцию выбора файлов, чтобы она возвращала ошибку отмены
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return nil, dialog.ErrCancelled
	}

	manager := NewManager()
	manager.AllTextures = []*ebiten.Image{ebiten.NewImage(1, 1)} // Предварительно заполняем

	manager.LoadTextures()

	// Текстуры не должны быть очищены
	assert.Len(t, manager.AllTextures, 1, "Текстуры не должны были измениться после отмены")
}

// TestLoadTextures_GenericError проверяет сценарий с общей ошибкой при выборе файла.
func TestLoadTextures_GenericError(t *testing.T) {
	// Подменяем функцию выбора файлов, чтобы она возвращала произвольную ошибку
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return nil, errors.New("some generic error")
	}

	manager := NewManager()
	manager.AllTextures = []*ebiten.Image{ebiten.NewImage(1, 1)} // Предварительно заполняем

	manager.LoadTextures()

	// Текстуры не должны быть очищены
	assert.Len(t, manager.AllTextures, 1, "Текстуры не должны были измениться при ошибке")
}

// TestLoadTextures_NoFilesSelected проверяет сценарий, когда файлы не были выбраны.
func TestLoadTextures_NoFilesSelected(t *testing.T) {
	// Подменяем функцию выбора файлов, чтобы она возвращала пустой срез
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return []string{}, nil
	}

	manager := NewManager()
	manager.AllTextures = []*ebiten.Image{ebiten.NewImage(1, 1)} // Предварительно заполняем

	manager.LoadTextures()

	// Текстуры не должны быть очищены
	assert.Len(t, manager.AllTextures, 1, "Текстуры не должны были измениться, если файлы не выбраны")
}

// TestLoadTextures_ClearsOldTextures проверяет, что старые текстуры очищаются при успешной новой загрузке.
func TestLoadTextures_ClearsOldTextures(t *testing.T) {
	// Подменяем функцию выбора файлов
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return []string{"/fake/new.png"}, nil
	}

	// Создаем менеджер с мок-загрузчиком и "старыми" текстурами
	manager := newTestManager(&mockTextureLoader{})
	manager.AllTextures = []*ebiten.Image{ebiten.NewImage(1, 1), ebiten.NewImage(1, 1)}
	manager.AvailableTextures = []*ebiten.Image{ebiten.NewImage(1, 1)}

	manager.LoadTextures()

	// Проверяем, что старые текстуры заменены новыми
	assert.Len(t, manager.AllTextures, 1, "Старые текстуры должны быть заменены одной новой")
	assert.Len(t, manager.AvailableTextures, 1, "Доступные текстуры должны быть заменены одной новой")
}

// TestLoadTextures_LoaderReturnsNil проверяет, что текстуры, которые не удалось загрузить, не добавляются.
func TestLoadTextures_LoaderReturnsNil(t *testing.T) {
	// Подменяем функцию выбора файлов
	originalShowFilePicker := ui.ShowFilePicker
	defer func() { ui.ShowFilePicker = originalShowFilePicker }()
	ui.ShowFilePicker = func() ([]string, error) {
		return []string{"/fake/good.png", "/fake/bad.png"}, nil
	}

	// Создаем менеджер с мок-загрузчиком, который иногда возвращает nil
	manager := newTestManager(&mockTextureLoader{failOnLoad: true})

	manager.LoadTextures()

	// Только одна текстура должна была быть успешно загружена
	assert.Len(t, manager.AllTextures, 0, "Только успешно загруженные текстуры должны быть добавлены")
	assert.Len(t, manager.AvailableTextures, 0, "Только успешно загруженные текстуры должны быть доступны")
}