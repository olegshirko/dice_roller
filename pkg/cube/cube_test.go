package cube

import (
	"math"
	"testing"

	"github.com/olegshirko/dice_roller/pkg/config"
	"github.com/stretchr/testify/assert"
)

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestNewCube(t *testing.T) {
	c := NewCube()

	assert.NotNil(t, c, "NewCube should not return nil")
	assert.Len(t, c.Vertices, 8, "Cube should have 8 vertices")
	assert.Len(t, c.Faces, 6, "Cube should have 6 faces")

	expectedVertices := [8]Point3D{
		{-config.CubeSize / 2, -config.CubeSize / 2, -config.CubeSize / 2},
		{config.CubeSize / 2, -config.CubeSize / 2, -config.CubeSize / 2},
		{config.CubeSize / 2, config.CubeSize / 2, -config.CubeSize / 2},
		{-config.CubeSize / 2, config.CubeSize / 2, -config.CubeSize / 2},
		{-config.CubeSize / 2, -config.CubeSize / 2, config.CubeSize / 2},
		{config.CubeSize / 2, -config.CubeSize / 2, config.CubeSize / 2},
		{config.CubeSize / 2, config.CubeSize / 2, config.CubeSize / 2},
		{-config.CubeSize / 2, config.CubeSize / 2, config.CubeSize / 2},
	}
	assert.Equal(t, expectedVertices, c.Vertices, "Vertices should be initialized correctly")

	expectedFaces := [6]Face{
		{Indices: [4]int{0, 1, 2, 3}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Back
		{Indices: [4]int{5, 4, 7, 6}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Front
		{Indices: [4]int{1, 5, 6, 2}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Right
		{Indices: [4]int{4, 0, 3, 7}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Left
		{Indices: [4]int{3, 2, 6, 7}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Top
		{Indices: [4]int{4, 5, 1, 0}, UVs: [4][2]float32{{0, 0}, {1, 0}, {1, 1}, {0, 1}}}, // Bottom
	}

	for i, face := range c.Faces {
		assert.Equal(t, expectedFaces[i].Indices, face.Indices, "Face %d indices should be correct", i)
		assert.Equal(t, expectedFaces[i].UVs, face.UVs, "Face %d UVs should be correct", i)
	}
}

func TestGetTargetAnglesForFace(t *testing.T) {
	pi_2 := math.Pi / 2
	cases := []struct {
		name      string
		faceIndex int
		expectedX float64
		expectedY float64
	}{
		{"Face 0 (Back)", 0, 0, math.Pi},
		{"Face 1 (Front)", 1, 0, 0},
		{"Face 2 (Right)", 2, 0, -pi_2},
		{"Face 3 (Left)", 3, 0, pi_2},
		{"Face 4 (Top)", 4, pi_2, 0},
		{"Face 5 (Bottom)", 5, -pi_2, 0},
		{"Invalid face index (-1)", -1, 0, 0},
		{"Invalid face index (6)", 6, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			x, y := GetTargetAnglesForFace(tc.faceIndex)
			assert.True(t, almostEqual(x, tc.expectedX), "face %d: expected X: %f, got %f", tc.faceIndex, tc.expectedX, x)
			assert.True(t, almostEqual(y, tc.expectedY), "face %d: expected Y: %f, got %f", tc.faceIndex, tc.expectedY, y)
		})
	}
}

func TestCalculateAlignmentAngle(t *testing.T) {
	cases := []struct {
		name          string
		faceIndex     int
		targetAngleX  float64
		targetAngleY  float64
		expectedAngle float64
	}{
		{"Face 0, no rotation", 0, 0, 0, 0},
		{"Face 1, 45 deg rotation", 1, math.Pi / 4, math.Pi / 4, 0},
		{"Face 4, top face", 4, 0, 0, -math.Pi / 2},
		{"Face 5, bottom face", 5, 0, 0, -math.Pi / 2},
		{"Invalid face index", -1, 0, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			angle := CalculateAlignmentAngle(tc.faceIndex, tc.targetAngleX, tc.targetAngleY)
			assert.True(t, almostEqual(angle, tc.expectedAngle), "test '%s': expected angle %f, got %f", tc.name, tc.expectedAngle, angle)
		})
	}
}
