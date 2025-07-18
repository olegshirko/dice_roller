package game

import (
	"github.com/olegshirko/dice_roller/pkg/assets"
	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/olegshirko/dice_roller/pkg/cube"
	"log"
	"math"
	"math/rand"
)

// StateManager управляет состоянием игры (вращение, остановка, выравнивание).
type StateManager struct {
	Cube              *cube.Cube
	AssetManager      *assets.Manager
	IsGrey            [6]bool // Статус "серости" граней
	IsWinner          [6]bool // Статус "победителя" граней
	AngleX, AngleY    float64
	AngleZ            float64
	RotationSpeedX    float64
	RotationSpeedY    float64
	Rotating          bool
	IdleRotating      bool // Новое состояние для вращения при простое
	Snapping          bool
	Aligning          bool
	TargetAngleX      float64
	TargetAngleY      float64
	TargetAngleZ      float64
	WinningFaceIndex  int
	LastWinnerIndex   int
	NeedsToRetireFace bool
	OffsetY           float64 // Смещение для прыжка
	jumpVelocity      float64 // Вертикальная скорость для прыжка
	isJumping         bool    // Флаг активности прыжка
	Shaking           bool    // Флаг для анимации дрожания
	shakeProgress     float64 // Прогресс анимации дрожания
}

// NewStateManager создает новый менеджер состояний.
func NewStateManager(c *cube.Cube, am *assets.Manager) *StateManager {
	sm := &StateManager{
		Cube:             c,
		AssetManager:     am,
		Rotating:         false,
		IdleRotating:     true, // Включаем по умолчанию
		WinningFaceIndex: -1,
		LastWinnerIndex:  -1,
		isJumping:        false,
		jumpVelocity:     0,
		OffsetY:          0,
		RotationSpeedX:   0.005, // Начальная скорость для медленного вращения
		RotationSpeedY:   0.01,
		Shaking:          false,
	}
	// Инициализируем грани пустыми текстурами.
	// Настоящие текстуры будут установлены позже из game.go
	for i := 0; i < 6; i++ {
		sm.Cube.Faces[i].Texture = config.EmptyImage
		sm.IsGrey[i] = true
		sm.IsWinner[i] = false
	}
	return sm
}

// StartRotation инкапсулирует логику запуска вращения куба.
func (sm *StateManager) StartRotation() {
	// Останавливаем фоновое вращение и сбрасываем скорости перед запуском основного цикла
	if sm.IdleRotating {
		sm.IdleRotating = false
		sm.RotationSpeedX = 0
		sm.RotationSpeedY = 0
	}

	if sm.Rotating || sm.Snapping || sm.Aligning {
		return
	}

	// 1. Собираем все валидные (не серые и не выигравшие) грани
	validFaceIndices := []int{}
	for i := 0; i < 6; i++ {
		if !sm.IsGrey[i] && !sm.IsWinner[i] {
			validFaceIndices = append(validFaceIndices, i)
		}
	}

	// 2. Если валидных граней нет, возможно, пора начать новый цикл
	if len(validFaceIndices) == 0 {
		anyActiveFaces := false
		for i := 0; i < 6; i++ {
			if !sm.IsGrey[i] {
				anyActiveFaces = true
				if sm.IsWinner[i] {
					sm.AssetManager.ReplaceFaceTexture(i, &sm.Cube.Faces, &sm.IsGrey)
					if !sm.IsGrey[i] {
						sm.IsWinner[i] = false
						validFaceIndices = append(validFaceIndices, i)
					}
				} else {
					validFaceIndices = append(validFaceIndices, i)
				}
			}
		}
		if !anyActiveFaces {
			log.Println("All faces are grey, no new cycle possible.")
		}
	}

	// 3. Если есть из чего выбирать, выбираем победителя

	// Подсчитываем активные (не серые) грани, чтобы определить, финальный ли это раунд
	activeFaceCount := 0
	for i := 0; i < 6; i++ {
		if !sm.IsGrey[i] {
			activeFaceCount++
		}
	}

	// Эффект "примагничивания" (snap) без вращения,
	// только если это последняя доступная грань в финальном раунде (когда осталось <=2 активных граней).
	if len(validFaceIndices) == 1 && activeFaceCount <= 2 {
		sm.WinningFaceIndex = validFaceIndices[0]
		sm.IsWinner[sm.WinningFaceIndex] = true

		sm.TargetAngleX, sm.TargetAngleY = cube.GetTargetAnglesForFace(sm.WinningFaceIndex)
		sm.AngleZ = 0
		sm.TargetAngleZ = 0

		sm.Snapping = true
		sm.Rotating = false
		sm.Aligning = false

	} else if len(validFaceIndices) > 0 {
		// Полноценное вращение для всех остальных случаев
		sm.AngleZ = 0
		sm.TargetAngleZ = 0

		// Выбираем победителя: если остался один, то он и есть, иначе - случайный
		if len(validFaceIndices) == 1 {
			sm.WinningFaceIndex = validFaceIndices[0]
		} else {
			sm.WinningFaceIndex = validFaceIndices[rand.Intn(len(validFaceIndices))]
		}
		sm.IsWinner[sm.WinningFaceIndex] = true

		sm.TargetAngleX, sm.TargetAngleY = cube.GetTargetAnglesForFace(sm.WinningFaceIndex)

		sm.RotationSpeedX = (rand.Float64() - 0.5) * 0.4
		sm.RotationSpeedY = (rand.Float64() - 0.5) * 0.4
		if math.Abs(sm.RotationSpeedX) < 0.05 && math.Abs(sm.RotationSpeedY) < 0.05 {
			sm.RotationSpeedX = 0.15 + rand.Float64()*0.1
		}

		sm.Rotating = true
		sm.Snapping = false
		sm.Aligning = false
		sm.isJumping = true
		sm.jumpVelocity = -20.0 // Начальная скорость прыжка вверх
	} else {
		// Если нет доступных граней, запускаем анимацию дрожания
		sm.Shaking = true
		sm.shakeProgress = 0
	}
}

// UpdateState обновляет состояние вращения и переходы между фазами.
// Возвращает true, если спин только что завершился.
func (sm *StateManager) UpdateState() (spinFinished bool) {
	spinFinished = false

	if sm.isJumping {
		sm.jumpVelocity += 1.0 // Гравитация
		sm.OffsetY += sm.jumpVelocity

		if sm.OffsetY >= 0 {
			sm.OffsetY = 0
			sm.jumpVelocity = -sm.jumpVelocity * 0.6 // Отскок с затуханием
		}
	}

	if sm.IdleRotating {
		sm.AngleX += sm.RotationSpeedX
		sm.AngleY += sm.RotationSpeedY
	} else if sm.Aligning {
		// ... (код состояния Aligning)
		snapSpeed := 0.1
		sm.AngleZ += (sm.TargetAngleZ - sm.AngleZ) * snapSpeed

		if math.Abs(sm.TargetAngleZ-sm.AngleZ) < 0.001 {
			sm.AngleZ = sm.TargetAngleZ
			sm.Aligning = false

			sm.LastWinnerIndex = sm.WinningFaceIndex
			sm.NeedsToRetireFace = true
			sm.WinningFaceIndex = -1
			spinFinished = true
		}
	} else if sm.Snapping {
		// ... (код состояния Snapping)
		snapSpeed := 0.1
		sm.AngleX += (sm.TargetAngleX - sm.AngleX) * snapSpeed
		sm.AngleY += (sm.TargetAngleY - sm.AngleY) * snapSpeed

		if math.Abs(sm.TargetAngleX-sm.AngleX) < 0.001 && math.Abs(sm.TargetAngleY-sm.AngleY) < 0.001 {
			sm.AngleX = sm.TargetAngleX
			sm.AngleY = sm.TargetAngleY
			sm.Snapping = false
			sm.TargetAngleZ = cube.CalculateAlignmentAngle(sm.WinningFaceIndex, sm.TargetAngleX, sm.TargetAngleY)
			sm.Aligning = true
		}
	} else if sm.Rotating {
		// ... (код состояния Rotating)
		if sm.NeedsToRetireFace {
			sm.LastWinnerIndex = -1
			sm.NeedsToRetireFace = false
		}

		sm.AngleX += sm.RotationSpeedX
		sm.AngleY += sm.RotationSpeedY

		sm.RotationSpeedX *= 0.99
		sm.RotationSpeedY *= 0.99

		// Логируем текущие скорости, чтобы видеть их затухание
		// log.Printf("UpdateState: Rotating... SpeedX: %.4f, SpeedY: %.4f", sm.RotationSpeedX, sm.RotationSpeedY)

		if math.Abs(sm.RotationSpeedX) < 0.01 && math.Abs(sm.RotationSpeedY) < 0.01 {
			sm.Rotating = false
			sm.isJumping = false
			sm.OffsetY = 0
			sm.jumpVelocity = 0
			sm.Snapping = true
		}
	} else if sm.Shaking {
		shakeSpeed := 0.5
		shakeMagnitude := 0.03
		sm.shakeProgress += shakeSpeed

		// Используем синусоиду для создания эффекта дрожания
		offset := math.Sin(sm.shakeProgress) * shakeMagnitude
		sm.AngleX += offset
		sm.AngleY -= offset // Дрожание в противофазе для лучшего эффекта

		// Завершаем анимацию после двух полных циклов синусоиды
		if sm.shakeProgress >= math.Pi*4 {
			sm.Shaking = false
			sm.shakeProgress = 0
			// Возвращаем углы в исходное состояние, чтобы дрожание не смещало куб
			sm.AngleX -= offset
			sm.AngleY += offset
		}
	}
	return
}
