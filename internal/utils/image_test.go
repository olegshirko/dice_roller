//go:build !ci

package utils

import (
	"image"
	"testing"
)

func TestAddLabelToImage(t *testing.T) {
	// Создаем "пустышку" изображения для теста
	dummyImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	label := "Test Label"

	// Проверяем, что функция не паникует
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic: %v", r)
		}
	}()

	// Вызываем тестируемую функцию
	resultImg := AddLabelToImage(dummyImg, label)

	// Проверяем, что результат не nil
	if resultImg == nil {
		t.Error("AddLabelToImage returned nil")
	}
}
