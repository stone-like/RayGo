package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

//Planeはxzに広がるものとしている、yz平面とかは対応していない
type Plane struct {
	*BaseShape
}

func NewPlane() Plane {
	return Plane{
		NewBaseShape(),
	}
}

func (p Plane) calcLocalNormal(localPoint calc.Tuple4) calc.Tuple4 {
	return calc.NewVector(0, 1, 0)
}

func (p Plane) NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error) {
	return p.ShapeNormalAt(worldPoint, p.calcLocalNormal)
}

func (p Plane) calcLocalIntersect(r Ray) (Intersections, error) {
	//rayのy要素がEPSILON以下ならxz平面に広がるだけのPlaneとは交差しない
	if math.Abs(r.Direction[1]) < util.EPSILON {
		return Intersections{}, nil
	}

	t := -r.Origin[1] / r.Direction[1]

	intersection := Intersection{
		Time:   t,
		Object: p,
	}

	return Intersections{
		Count:         1,
		Intersections: []*Intersection{&intersection},
	}, nil
}

func (p Plane) Intersect(r Ray) (Intersections, error) {
	return p.ShapeIntersect(r, p.calcLocalIntersect)
}
