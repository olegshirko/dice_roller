//go:build !ci

package graphics

import (
	"testing"

	"github.com/olegshirko/dice_roller/pkg/cube"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderer(t *testing.T) {
	renderer := NewRenderer()
	assert.NotNil(t, renderer, "NewRenderer() should not return nil")
}

func TestDrawCube(t *testing.T) {
	// Создаем "дымовой" тест, чтобы убедиться, что отрисовка не вызывает панику.
	// Этот тест не проверяет фактический результат рендеринга, а только то,
	// что код выполняется без ошибок.

	// 1. Создаем Renderer
	renderer := NewRenderer()

	// 2. Создаем "экран" для рисования
	screen := ebiten.NewImage(100, 100)

	// 3. Создаем куб
	c := cube.NewCube()

	// 4. Инициализируем текстуры для граней, чтобы избежать паники
	//    при доступе к nil-указателю face.Texture.
	for i := range c.Faces {
		c.Faces[i].Texture = ebiten.NewImage(1, 1)
	}

	// 5. Вызываем DrawCube и проверяем, что паники не произошло
	assert.NotPanics(t, func() {
		renderer.DrawCube(screen, c, 0, 0, 0, 0)
	}, "DrawCube should not panic")
}