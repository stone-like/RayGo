package scene

import (
	"rayGo/calc"
)

type Group struct {
	*BaseShape
	Children []Shape
}

var _ Shape = &Group{}

func NewGroup() *Group {

	return &Group{
		NewBaseShape(),
		[]Shape{},
	}

}

func (g *Group) AddChildren(ss ...Shape) {
	g.Children = append(g.Children, ss...)

	for _, shape := range ss {
		shape.SetParent(g)
	}
}

func (g *Group) calcLocalNormal(localPoint calc.Tuple4) calc.Tuple4 {
	return calc.Tuple4{}
}

func (g *Group) NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error) {
	return g.ShapeNormalAt(worldPoint, g.calcLocalNormal)
}

func (g *Group) calcLocalIntersect(r Ray) (Intersections, error) {

	if len(g.Children) == 0 {
		return Intersections{}, nil
	}

	var xs []*Intersection

	for _, child := range g.Children {
		section, err := child.Intersect(r)

		if err != nil {
			return Intersections{}, err
		}

		xs = append(xs, section.Intersections...)
	}

	return AggregateIntersection(xs...), nil
}

func (g *Group) Intersect(r Ray) (Intersections, error) {
	return g.ShapeIntersect(r, g.calcLocalIntersect)
}

func (g *Group) GetMaterial() *Material {
	return g.Material
}

func (g *Group) SetMaterial(m *Material) {
	g.Material = m
}
