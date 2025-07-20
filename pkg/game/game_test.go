//go:build !ci

package game

import (
	"testing"

	"github.com/olegshirko/dice_roller/pkg/assets"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/cube"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRenderer is a mock implementation of the Renderer.
type MockRenderer struct {
	mock.Mock
}

func (m *MockRenderer) DrawCube(screen *ebiten.Image, c *cube.Cube, angleX, angleY, angleZ, offsetY float64) {
	m.Called(screen, c, angleX, angleY, angleZ, offsetY)
}

// MockAssetManager is a mock implementation of the AssetManager.
type MockAssetManager struct {
	mock.Mock
}

func (m *MockAssetManager) LoadTextures() {
	m.Called()
}

func (m *MockAssetManager) SetInitialTextures(faces *[6]cube.Face, isGrey *[6]bool) {
	m.Called(faces, isGrey)
}

func (m *MockAssetManager) ReplaceFaceTexture(faceIndex int, faces *[6]cube.Face, isGrey *[6]bool) {
	m.Called(faceIndex, faces, isGrey)
}

func TestGame_Update(t *testing.T) {
	assetManager := &assets.Manager{}
	game := NewGame(assetManager)

	// Test initial state
	assert.NotNil(t, game)
	assert.NotNil(t, game.Cube)
	assert.NotNil(t, game.StateManager)
	assert.NotNil(t, game.Renderer)

	// Simulate key press
	// This is tricky to test without a running game loop.
	// We will focus on the state changes that Update triggers.
	initialAngleX := game.StateManager.AngleX
	initialAngleY := game.StateManager.AngleY

	// To properly test Update, we would need to simulate key presses.
	// For now, we just call Update and check that it doesn't panic and that the state is updated.
	err := game.Update()
	assert.NoError(t, err)

	// Since IdleRotating is true by default, the angles should change.
	assert.NotEqual(t, initialAngleX, game.StateManager.AngleX, "AngleX should change due to idle rotation")
	assert.NotEqual(t, initialAngleY, game.StateManager.AngleY, "AngleY should change due to idle rotation")
}

func TestGame_Draw(t *testing.T) {
	assetManager := &assets.Manager{}
	game := NewGame(assetManager)

	mockRenderer := new(MockRenderer)
	game.Renderer = mockRenderer // Inject mock renderer

	// Create a dummy screen
	screen := ebiten.NewImage(100, 100)

	// Set up the mock expectation
	mockRenderer.On("DrawCube", screen, game.Cube, game.StateManager.AngleX, game.StateManager.AngleY, game.StateManager.AngleZ, game.StateManager.OffsetY)

	// Call the method
	game.Draw(screen)

	// Assert that the expected method was called
	mockRenderer.AssertExpectations(t)
}

func TestGame_Layout(t *testing.T) {
	game := &Game{}
	width, height := game.Layout(800, 600)
	assert.Equal(t, config.ScreenWidth, width)
	assert.Equal(t, config.ScreenHeight, height)
}
