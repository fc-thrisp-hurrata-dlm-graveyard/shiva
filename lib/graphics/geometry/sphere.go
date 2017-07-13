package geometry

import (
	"github.com/Laughs-In-Flowers/shiva/lib/graphics"
	"github.com/Laughs-In-Flowers/shiva/lib/math"

	glm "math"
)

type Sphere struct {
	Geometry
	Radius         float64
	WidthSegments  int
	HeightSegments int
	PhiStart       float64
	PhiLength      float64
	ThetaStart     float64
	ThetaLength    float64
}

func NewSphere(
	radius float64,
	wSegments, hSegments int,
	phiStart, phiLength, thetaStart, thetaLength float64,
) *Sphere {
	s := &Sphere{
		New(),
		radius,
		wSegments, hSegments,
		phiStart, phiLength, thetaStart, thetaLength,
	}
	thetaEnd := thetaStart + thetaLength
	vertexCount := (wSegments + 1) * (hSegments + 1)

	positions := math.NewAF32(vertexCount*3, vertexCount*3)
	normals := math.NewAF32(vertexCount*3, vertexCount*3)
	uvs := math.NewAF32(vertexCount*2, vertexCount*2)
	indices := math.NewAU32(0, vertexCount)

	index := 0
	vertices := make([][]uint32, 0)
	var normal math.Vector = math.Vec3(0, 0, 0)

	for y := 0; y <= hSegments; y++ {
		verticesRow := make([]uint32, 0)
		v := float64(y) / float64(hSegments)
		for x := 0; x <= wSegments; x++ {
			u := float64(x) / float64(wSegments)
			px := -radius * glm.Cos(phiStart+u*phiLength) * glm.Sin(thetaStart+v*thetaLength)
			py := radius * glm.Cos(thetaStart+v*thetaLength)
			pz := radius * glm.Sin(phiStart+u*phiLength) * glm.Sin(thetaStart+v*thetaLength)
			normal.Set(1, float32(px))
			normal.Set(2, float32(py))
			normal.Set(3, float32(pz))
			normal.Normalize()

			positions.Set(index*3, float32(px), float32(py), float32(pz))
			//normals.SetVector3(index*3, &normal)
			//uvs.Set(index*2, float32(u), float32(v))
			verticesRow = append(verticesRow, uint32(index))
			index++
		}
		vertices = append(vertices, verticesRow)
	}

	for y := 0; y < hSegments; y++ {
		for x := 0; x < wSegments; x++ {
			v1 := vertices[y][x+1]
			v2 := vertices[y][x]
			v3 := vertices[y+1][x]
			v4 := vertices[y+1][x+1]
			if y != 0 || thetaStart > 0 {
				indices.Append(v1, v2, v4)
			}
			if y != hSegments-1 || thetaEnd < glm.Pi {
				indices.Append(v2, v3, v4)
			}
		}
	}

	s.SetIndices(indices)
	s.AddVBO(graphics.NewBuff().AddAttrib("VertexPosition", 3).SetBuffer(positions))
	s.AddVBO(graphics.NewBuff().AddAttrib("VertexNormal", 3).SetBuffer(normals))
	s.AddVBO(graphics.NewBuff().AddAttrib("VertexTexcoord", 2).SetBuffer(uvs))
	return s
}
