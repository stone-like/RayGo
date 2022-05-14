package scene

import (
	"rayGo/calc"
	"rayGo/util"
)

type Triangle struct {
	*BaseShape
	P1        calc.Tuple4
	P2        calc.Tuple4
	P3        calc.Tuple4
	E1        calc.Tuple4
	E2        calc.Tuple4
	NormalVec calc.Tuple4
}

func NewTriangle(p1, p2, p3 calc.Tuple4) Triangle {

	e1 := calc.SubTuple(p2, p1)
	e2 := calc.SubTuple(p3, p1)
	normal_vec := calc.CrossTuple(e2, e1).Normalize()

	return Triangle{
		NewBaseShape(),
		p1,
		p2,
		p3,
		e1,
		e2,
		normal_vec,
	}
}

var _ Shape = Triangle{}

func (tri Triangle) calcLocalNormal(localPoint calc.Tuple4, hit Intersection) calc.Tuple4 {
	return tri.NormalVec
}

func (tri Triangle) NormalAt(worldPoint calc.Tuple4, hit Intersection) (calc.Tuple4, error) {
	return tri.ShapeNormalAt(worldPoint, hit, tri.calcLocalNormal)
}

//borrow from Moller-Trumbore
func (tri Triangle) calcLocalIntersect(r Ray) (Intersections, error) {
	dir_cross_e2 := calc.CrossTuple(r.Direction, tri.E2)

	det := calc.DotTuple(tri.E1, dir_cross_e2)

	if util.IsNearlyEqualZero(det) {
		return Intersections{}, nil
	}

	f := 1.0 / det

	p1_to_origin := calc.SubTuple(r.Origin, tri.P1)

	u := f * calc.DotTuple(p1_to_origin, dir_cross_e2)

	if u < 0 || 1 < u {
		return Intersections{}, nil
	}

	origin_cross_e1 := calc.CrossTuple(p1_to_origin, tri.E1)
	v := f * calc.DotTuple(r.Direction, origin_cross_e1)

	if v < 0 || 1 < (u+v) {
		return Intersections{}, nil
	}

	t := f * calc.DotTuple(tri.E2, origin_cross_e1)

	return AggregateIntersection(&Intersection{
		t, tri, 0, 0,
	}), nil
}

func (tri Triangle) Intersect(r Ray) (Intersections, error) {
	return tri.ShapeIntersect(r, tri.calcLocalIntersect)
}

func (tri Triangle) GetMaterial() *Material {
	return tri.Material
}

func (tri Triangle) SetMaterial(m *Material) {
	tri.Material = m
}
