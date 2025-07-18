package utils

import (
	"image"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// AddLabelToImage добавляет текстовую метку на изображение.
func AddLabelToImage(img image.Image, label string) *ebiten.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, img, image.Point{}, draw.Src)

	// Настройка шрифта и текста
	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  newImg,
		Src:  image.White,
		Face: face,
	}

	// Вычисляем позицию для центрирования текста
	textWidth := drawer.MeasureString(label)
	x := (fixed.I(bounds.Dx()) - textWidth) / 2
	y := fixed.I(bounds.Dy() - face.Metrics().Height.Ceil()*2) // Смещаем текст вниз

	// Рисуем "тень" или контур для лучшей читаемости
	outlineColor := image.Black
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			drawer.Dot = fixed.Point26_6{X: x + fixed.I(dx), Y: y + fixed.I(dy)}
			drawer.Src = outlineColor
			drawer.DrawString(label)
		}
	}

	// Рисуем основной текст
	drawer.Dot = fixed.Point26_6{X: x, Y: y}
	drawer.Src = image.White
	drawer.DrawString(label)

	return ebiten.NewImageFromImage(newImg)
}
