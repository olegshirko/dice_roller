package config

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 720
	CubeSize     = 150
)

var (
	EmptyImage *ebiten.Image
	GreyImage  *ebiten.Image
)

func init() {
	EmptyImage = ebiten.NewImage(3, 3)
	EmptyImage.Fill(color.White)

	GreyImage = ebiten.NewImage(CubeSize, CubeSize)
	GreyImage.Fill(color.Gray{Y: 128})

	rand.Seed(time.Now().UnixNano())
}
