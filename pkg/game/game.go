package game

import (
	"github.com/olegshirko/dice_roller/pkg/assets"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/cube"
	"github.com/olegshirko/dice_roller/pkg/graphics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Cube         *cube.Cube
	AssetManager *assets.Manager
	StateManager *StateManager
	Renderer     *graphics.Renderer
}

// NewGame создает новую игру.
func NewGame(assetManager *assets.Manager) *Game {
	c := cube.NewCube()
	sm := NewStateManager(c, assetManager)
	r := graphics.NewRenderer()

	g := &Game{
		Cube:         c,
		AssetManager: assetManager,
		StateManager: sm,
		Renderer:     r,
	}

	// Устанавливаем начальные текстуры, если они были загружены
	assetManager.SetInitialTextures(&g.Cube.Faces, &g.StateManager.IsGrey)

	return g
}

// Update выполняется каждый такт (tick).
func (g *Game) Update() error {
	// Обработка пользовательского ввода
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		go func() {
			g.AssetManager.LoadTextures()
			g.AssetManager.SetInitialTextures(&g.Cube.Faces, &g.StateManager.IsGrey)
			g.StateManager.LastWinnerIndex = -1
		}()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.StateManager.StartRotation()
	}

	// Обновляем состояние игры (вращение, и т.д.)
	g.StateManager.UpdateState()

	return nil
}

// Draw выполняется каждый кадр (frame).
func (g *Game) Draw(screen *ebiten.Image) {
	g.Renderer.DrawCube(screen, g.Cube, g.StateManager.AngleX, g.StateManager.AngleY, g.StateManager.AngleZ, g.StateManager.OffsetY)
}

// Layout принимает логические размеры экрана и возвращает физические размеры.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
