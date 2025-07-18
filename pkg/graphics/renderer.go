package graphics

import (
	"ebit-hello/pkg/config"
	"ebit-hello/pkg/cube"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer struct {
	// Пока рендерер не имеет состояния, но может получить его в будущем (например, кэши).
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// DrawCube отрисовывает куб на экране.
func (r *Renderer) DrawCube(screen *ebiten.Image, c *cube.Cube, angleX, angleY, angleZ, offsetY float64) {
	screen.Fill(color.Transparent)
	ebitenutil.DebugPrint(screen, "Press 'L' to load textures, 'S' to spin")

	cosX, sinX := math.Cos(angleX), math.Sin(angleX)
	cosY, sinY := math.Cos(angleY), math.Sin(angleY)
	cosZ, sinZ := math.Cos(angleZ), math.Sin(angleZ)

	type RotatedPoint struct {
		cube.Point3D
		ProjX, ProjY float64
	}
	rotatedPoints := make([]RotatedPoint, len(c.Vertices))

	for i, v := range c.Vertices {
		// Вращение вокруг оси Y
		rotatedY := cube.Point3D{
			X: v.X*cosY - v.Z*sinY,
			Y: v.Y,
			Z: v.X*sinY + v.Z*cosY,
		}
		// Вращение вокруг оси X
		rotatedX := cube.Point3D{
			X: rotatedY.X,
			Y: rotatedY.Y*cosX - rotatedY.Z*sinX,
			Z: rotatedY.Y*sinX + rotatedY.Z*cosX,
		}
		// Финальный доворот для выравнивания
		finalRotated := cube.Point3D{
			X: rotatedX.X*cosZ - rotatedX.Y*sinZ,
			Y: rotatedX.X*sinZ + rotatedX.Y*cosZ,
			Z: rotatedX.Z,
		}

		scale := 1.5
		rotatedPoints[i] = RotatedPoint{
			Point3D: finalRotated,
			ProjX:   finalRotated.X*scale + config.ScreenWidth/2,
			ProjY:   finalRotated.Y*scale + config.ScreenHeight/2 + offsetY,
		}
	}

	type faceToSort struct {
		face     cube.Face
		averageZ float64
	}
	sortedFaces := make([]faceToSort, 0, len(c.Faces))

	for _, face := range c.Faces {
		avgZ := (rotatedPoints[face.Indices[0]].Z +
			rotatedPoints[face.Indices[1]].Z +
			rotatedPoints[face.Indices[2]].Z +
			rotatedPoints[face.Indices[3]].Z) / 4.0

		// Back-face culling
		v0 := rotatedPoints[face.Indices[0]].Point3D
		v1 := rotatedPoints[face.Indices[1]].Point3D
		v2 := rotatedPoints[face.Indices[2]].Point3D
		u := cube.Point3D{X: v1.X - v0.X, Y: v1.Y - v0.Y, Z: v1.Z - v0.Z}
		v := cube.Point3D{X: v2.X - v0.X, Y: v2.Y - v0.Y, Z: v2.Z - v0.Z}
		normal := cube.Point3D{
			X: u.Y*v.Z - u.Z*v.Y,
			Y: u.Z*v.X - u.X*v.Z,
			Z: u.X*v.Y - u.Y*v.X,
		}
		if normal.Z > 0 {
			sortedFaces = append(sortedFaces, faceToSort{face: face, averageZ: avgZ})
		}
	}

	// Сортировка граней для правильного отображения (Painter's algorithm)
	sort.Slice(sortedFaces, func(i, j int) bool {
		return sortedFaces[i].averageZ < sortedFaces[j].averageZ
	})

	for _, fts := range sortedFaces {
		face := fts.face
		p0 := rotatedPoints[face.Indices[0]]
		p1 := rotatedPoints[face.Indices[1]]
		p2 := rotatedPoints[face.Indices[2]]
		p3 := rotatedPoints[face.Indices[3]]

		texWidth, texHeight := face.Texture.Size()

		v0 := ebiten.Vertex{DstX: float32(p0.ProjX), DstY: float32(p0.ProjY), SrcX: face.UVs[0][0] * float32(texWidth), SrcY: face.UVs[0][1] * float32(texHeight), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}
		v1 := ebiten.Vertex{DstX: float32(p1.ProjX), DstY: float32(p1.ProjY), SrcX: face.UVs[1][0] * float32(texWidth), SrcY: face.UVs[1][1] * float32(texHeight), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}
		v2 := ebiten.Vertex{DstX: float32(p2.ProjX), DstY: float32(p2.ProjY), SrcX: face.UVs[2][0] * float32(texWidth), SrcY: face.UVs[2][1] * float32(texHeight), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}
		v3 := ebiten.Vertex{DstX: float32(p3.ProjX), DstY: float32(p3.ProjY), SrcX: face.UVs[3][0] * float32(texWidth), SrcY: face.UVs[3][1] * float32(texHeight), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}

		op := &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.FillAll,
		}

		screen.DrawTriangles([]ebiten.Vertex{v0, v1, v2}, []uint16{0, 1, 2}, face.Texture, op)
		screen.DrawTriangles([]ebiten.Vertex{v0, v2, v3}, []uint16{0, 1, 2}, face.Texture, op)
	}
}
