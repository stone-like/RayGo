package scene

import (
	"rayGo/calc"
)

func CreateIntersection(t float64, object Shape) Intersection {
	return Intersection{
		Time:   t,
		Object: object,
	}
}

func GlassSphere() Sphere {
	s := NewSphere(1)

	m := DefaultMaterial()
	m.Transparency = 1.0
	m.RefractiveIndex = 1.5
	s.SetMaterial(m)

	return s
}

type TestPattern struct {
	*BasePattern
}

var _ Pattern = TestPattern{}

func NewTestPattern() TestPattern {
	return TestPattern{
		NewBasePattern(),
	}
}

func (tp TestPattern) PatternAt(point calc.Tuple4) Color {
	return TupletoColor(point)
}

func (tp TestPattern) PatternAtShape(world_point calc.Tuple4, shape Shape) (Color, error) {

	return tp.PatternAtShapeOnBase(world_point, shape, tp.PatternAt)
}
