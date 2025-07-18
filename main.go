package main

import (
	"github.com/olegshirko/dice_roller/pkg/assets"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Rotating 3D Cube")

	assetManager := assets.NewManager()
	assetManager.LoadFromDirectory("img")

	g := game.NewGame(assetManager)

	if err := ebiten.RunGame(g); err != nil {
		if err != ebiten.Termination {
			log.Fatal(err)
		}
	}
	log.Println("Game finished.")
}
