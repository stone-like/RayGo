package scene

import (
	"rayGo/calc"
	"sort"
)

type Shape interface {
	Intersect(r Ray) (Intersections, error)
	NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error)
	GetMaterial() Material
}

type BaseShape struct {
	Transform calc.Mat4x4
	Material  Material
}

func NewBaseShape() *BaseShape {
	return &BaseShape{
		Transform: calc.Ident4x4,
		Material:  DefaultMaterial(),
	}
}

func (b *BaseShape) GetTransform() calc.Mat4x4 {
	return b.Transform
}

func (b *BaseShape) SetTransform(mat calc.Mat4x4) {
	b.Transform = mat
}

func (b *BaseShape) GetMaterial() Material {
	return b.Material
}

func (b *BaseShape) SetMaterial(m Material) {
	b.Material = m
}

type Intersection struct {
	Time   float64
	Object Shape
}

type Intersections struct {
	Intersections []*Intersection
	Count         int
}

func AggregateIntersection(sections ...*Intersection) Intersections {

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Time < sections[j].Time
	})

	return Intersections{
		Count:         len(sections),
		Intersections: sections,
	}

}

//hitはintersectionsの中で最初の正のtimeをもつintersection
func GenerateHit(intersections Intersections) *Intersection {
	for _, intersection := range intersections.Intersections {
		if intersection.Time >= 0 {
			return intersection
		}
	}

	//正のtimeのintersectionがなかったらnil
	return nil
}
