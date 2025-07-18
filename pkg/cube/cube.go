package cube

import (
	"math"

	"github.com/olegshirko/dice_roller/pkg/config"

	"github.com/hajimehoshi/ebiten/v2"
)

// Point3D представляет собой вершину в 3D пространстве.
type Point3D struct {
	X, Y, Z float64
}

// Face представляет грань куба.
type Face struct {
	Indices [4]int        // Индексы вершин, образующих грань
	Texture *ebiten.Image // Текстура грани
	UVs     [4][2]float32 // UV-координаты для каждой вершины
}

// Cube содержит геометрию куба.
type Cube struct {
	Vertices [8]Point3D
	Faces    [6]Face
}

// NewCube создает новый экземпляр куба с определенными вершинами и гранями.
func NewCube() *Cube {
	// Определяем 8 вершин куба
	vertices := [8]Point3D{
		{-config.CubeSize / 2, -config.CubeSize / 2, -config.CubeSize / 2}, // 0
		{config.CubeSize / 2, -config.CubeSize / 2, -config.CubeSize / 2},  // 1
		{config.CubeSize / 2, config.CubeSize / 2, -config.CubeSize / 2},   // 2
		{-config.CubeSize / 2, config.CubeSize / 2, -config.CubeSize / 2},  // 3
		{-config.CubeSize / 2, -config.CubeSize / 2, config.CubeSize / 2},  // 4
		{config.CubeSize / 2, -config.CubeSize / 2, config.CubeSize / 2},   // 5
		{config.CubeSize / 2, config.CubeSize / 2, config.CubeSize / 2},    // 6
		{-config.CubeSize / 2, config.CubeSize / 2, config.CubeSize / 2},   // 7
	}

	// Определяем 6 граней куба
	faces := [6]Face{
		{Indices: [4]int{0, 1, 2, 3}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Задняя
		{Indices: [4]int{5, 4, 7, 6}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Передняя
		{Indices: [4]int{1, 5, 6, 2}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Правая
		{Indices: [4]int{4, 0, 3, 7}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Левая
		{Indices: [4]int{3, 2, 6, 7}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Верхняя
		{Indices: [4]int{4, 5, 1, 0}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Нижняя
	}

	return &Cube{
		Vertices: vertices,
		Faces:    faces,
	}
}

// GetTargetAnglesForFace вычисляет целевые углы X и Y для ориентации на грань.
func GetTargetAnglesForFace(faceIndex int) (float64, float64) {
	pi_2 := math.Pi / 2
	switch faceIndex {
	case 0:
		return 0, math.Pi
	case 1:
		return 0, 0
	case 2:
		return 0, -pi_2
	case 3:
		return 0, pi_2
	case 4:
		return pi_2, 0
	case 5:
		return -pi_2, 0
	default:
		return 0, 0
	}
}

// CalculateAlignmentAngle вычисляет угол, на который нужно довернуть куб для выравнивания.
func CalculateAlignmentAngle(faceIndex int, targetAngleX, targetAngleY float64) float64 {
	var upVector Point3D
	switch faceIndex {
	case 0, 1, 2, 3:
		upVector = Point3D{X: 0, Y: -1, Z: 0}
	case 4:
		upVector = Point3D{X: 0, Y: 0, Z: 1}
	case 5:
		upVector = Point3D{X: 0, Y: 0, Z: -1}
	default:
		return 0
	}

	cosX, sinX := math.Cos(targetAngleX), math.Sin(targetAngleX)
	cosY, sinY := math.Cos(targetAngleY), math.Sin(targetAngleY)

	rotatedY := Point3D{
		X: upVector.X*cosY - upVector.Z*sinY,
		Y: upVector.Y,
		Z: upVector.X*sinY + upVector.Z*cosY,
	}
	rotatedX := Point3D{
		X: rotatedY.X,
		Y: rotatedY.Y*cosX - rotatedY.Z*sinX,
		Z: rotatedY.Y*sinX + rotatedY.Z*cosX,
	}

	currentAngle := math.Atan2(rotatedX.Y, rotatedX.X)
	targetScreenAngle := -math.Pi / 2
	alignmentAngle := targetScreenAngle - currentAngle

	for alignmentAngle <= -math.Pi {
		alignmentAngle += 2 * math.Pi
	}
	for alignmentAngle > math.Pi {
		alignmentAngle -= 2 * math.Pi
	}
	return alignmentAngle
}
