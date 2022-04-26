package scene

import (
	"rayGo/calc"
)

//Origin -> Point,Direction -> Vector
type Ray struct {
	Origin    calc.Tuple4
	Direction calc.Tuple4
}

func NewRay(origin, direction calc.Tuple4) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
	}
}

//return Point
func (r Ray) Position(t float64) calc.Tuple4 {
	return calc.AddTuple(r.Origin, calc.MulTupleByScalar(t, r.Direction))
}

func (r Ray) Transform(mat calc.Mat4x4) Ray {
	return Ray{
		Origin:    mat.MulByTuple(r.Origin),
		Direction: mat.MulByTuple(r.Direction),
	}
}
