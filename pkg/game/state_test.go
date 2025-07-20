package game

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olegshirko/dice_roller/pkg/assets"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/cube"
	"github.com/stretchr/testify/assert"
)

func TestNewStateManager(t *testing.T) {
	c := cube.NewCube()
	am := assets.NewManager()
	sm := NewStateManager(c, am)

	assert.NotNil(t, sm)
	assert.Equal(t, c, sm.Cube)
	assert.Equal(t, am, sm.AssetManager)
	assert.False(t, sm.Rotating, "Rotating should be false on init")
	assert.True(t, sm.IdleRotating, "IdleRotating should be true on init")
	assert.Equal(t, -1, sm.WinningFaceIndex, "WinningFaceIndex should be -1 on init")
	assert.Equal(t, -1, sm.LastWinnerIndex, "LastWinnerIndex should be -1 on init")
	assert.False(t, sm.Shaking, "Shaking should be false on init")

	for i := 0; i < 6; i++ {
		assert.True(t, sm.IsGrey[i], "Face %d should be grey on init", i)
		assert.False(t, sm.IsWinner[i], "Face %d should not be a winner on init", i)
		assert.Equal(t, config.EmptyImage, sm.Cube.Faces[i].Texture, "Face %d should have empty texture on init", i)
	}
}

func TestStartRotation(t *testing.T) {
	// Сценарий 1: Нормальное вращение
	t.Run("Normal Rotation", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		// Делаем несколько граней активными (не серыми)
		sm.IsGrey[0] = false
		sm.IsGrey[1] = false
		sm.IsGrey[2] = false

		sm.StartRotation()

		assert.True(t, sm.Rotating, "Should be in Rotating state")
		assert.False(t, sm.Snapping, "Should not be in Snapping state")
		assert.False(t, sm.IdleRotating, "IdleRotating should be disabled")
		assert.Contains(t, []int{0, 1, 2}, sm.WinningFaceIndex, "Winning face should be one of the valid faces")
		assert.True(t, sm.IsWinner[sm.WinningFaceIndex], "Winning face should be marked as winner")
	})

	// Сценарий 2: Последняя доступная грань
	t.Run("Last Available Face Snap", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		// Все грани неактивны, кроме одной
		for i := range sm.IsGrey {
			sm.IsGrey[i] = true
		}
		sm.IsGrey[3] = false

		sm.StartRotation()

		assert.False(t, sm.Rotating, "Should not be in Rotating state")
		assert.True(t, sm.Snapping, "Should be in Snapping state")
		assert.Equal(t, 3, sm.WinningFaceIndex, "Winning face should be the last available one")
		assert.True(t, sm.IsWinner[3], "The winning face should be marked as winner")
	})

	// Сценарий 3: Нет доступных граней (все серые, но не выигравшие)
	t.Run("No Available Faces Shake", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		// Все грани активны, но уже выиграли
		for i := range sm.IsGrey {
			sm.IsGrey[i] = false
			sm.IsWinner[i] = true
		}
		// ... кроме одной, которая серая
		sm.IsGrey[5] = true

		sm.StartRotation()

		assert.True(t, sm.Shaking, "Should be in Shaking state")
		assert.False(t, sm.Rotating, "Should not be in Rotating state")
	})

	// Сценарий 4: Перезапуск цикла
	t.Run("Restart Cycle when all are winners", func(t *testing.T) {
		am := assets.NewManager()
		// Добавляем "фейковые" текстуры, чтобы было что заменять
		am.AvailableTextures = make([]*ebiten.Image, 6)
		for i := 0; i < 6; i++ {
			am.AvailableTextures[i] = ebiten.NewImage(1, 1)
		}

		sm := NewStateManager(cube.NewCube(), am)
		// Все грани не серые и все выиграли
		for i := 0; i < 6; i++ {
			sm.IsGrey[i] = false
			sm.IsWinner[i] = true
		}

		initialTextureCount := len(am.AvailableTextures)
		sm.StartRotation()

		assert.True(t, sm.Rotating, "Should start rotating again")
		// Проверяем, что текстуры были заменены (их количество уменьшилось)
		assert.Less(t, len(am.AvailableTextures), initialTextureCount, "Available textures should decrease after replacing faces")

		// Проверяем, что статусы IsWinner были сброшены
		for i := 0; i < 6; i++ {
			if i == sm.WinningFaceIndex {
				assert.True(t, sm.IsWinner[i], "New winning face should be marked as winner")
			} else {
				assert.False(t, sm.IsWinner[i], "Non-winning faces should have their winner status reset")
			}
		}
	})
}

func TestUpdateState(t *testing.T) {
	// Тест перехода Rotating -> Snapping
	t.Run("Rotating to Snapping", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		sm.IdleRotating = false // Отключаем, чтобы тестировать именно Rotating
		sm.Rotating = true
		sm.RotationSpeedX = 0.005 // Малая скорость для быстрого перехода
		sm.RotationSpeedY = 0.005

		sm.UpdateState()

		assert.False(t, sm.Rotating, "Should exit Rotating state")
		assert.True(t, sm.Snapping, "Should enter Snapping state")
	})

	// Тест перехода Snapping -> Aligning
	t.Run("Snapping to Aligning", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		sm.IdleRotating = false // Отключаем, чтобы тестировать именно Snapping
		sm.Snapping = true
		sm.WinningFaceIndex = 1
		sm.TargetAngleX, sm.TargetAngleY = cube.GetTargetAnglesForFace(sm.WinningFaceIndex)
		// Устанавливаем углы близко к целевым для быстрого перехода
		sm.AngleX = sm.TargetAngleX - 0.0001
		sm.AngleY = sm.TargetAngleY - 0.0001

		sm.UpdateState()

		assert.False(t, sm.Snapping, "Should exit Snapping state")
		assert.True(t, sm.Aligning, "Should enter Aligning state")
		// Проверяем, что целевой угол для выравнивания был рассчитан
		assert.NotEqual(t, 0, sm.TargetAngleZ, "TargetAngleZ should be calculated for alignment")
	})

	// Тест перехода Aligning -> Finished
	t.Run("Aligning to Finished", func(t *testing.T) {
		am := assets.NewManager()
		sm := NewStateManager(cube.NewCube(), am)
		sm.IdleRotating = false // Отключаем, чтобы тестировать именно Aligning
		sm.Aligning = true
		sm.WinningFaceIndex = 2
		sm.TargetAngleZ = 1.0
		// Устанавливаем угол близко к целевому
		sm.AngleZ = sm.TargetAngleZ - 0.0001

		spinFinished := sm.UpdateState()

		assert.True(t, spinFinished, "UpdateState should return true to indicate spin finished")
		assert.False(t, sm.Aligning, "Should exit Aligning state")
		assert.True(t, sm.NeedsToRetireFace, "NeedsToRetireFace should be true")
		assert.Equal(t, 2, sm.LastWinnerIndex, "LastWinnerIndex should be updated")
		assert.Equal(t, -1, sm.WinningFaceIndex, "WinningFaceIndex should be reset")
	})
}
