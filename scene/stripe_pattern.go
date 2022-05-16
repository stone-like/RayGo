package scene

import (
	"math"
	"rayGo/calc"
)

type StripePattern struct {
	*BasePattern
	Color1 Color
	Color2 Color
}

var _ Pattern = StripePattern{}

func NewStripePattern(c1, c2 Color) StripePattern {
	return StripePattern{
		NewBasePattern(),
		c1,
		c2,
	}
}

//math.Floorは与えられたfloat64以下の最大の整数を返す
//ex. 1.001 -> 1
//  1 -> 1
//  0.99 -> 0
//  -0.1 -> -1
//  -0.9999 -> -1
func (sp StripePattern) PatternAt(point calc.Tuple4) Color {
	if int(math.Floor(point[0]))%2 == 0 {
		return sp.Color1
	}

	return sp.Color2
}

func (sp StripePattern) PatternAtShape(world_point calc.Tuple4, shape Shape) (Color, error) {
	// shapeTransInv, err := shape.GetTransform().Inverse()

	// if err != nil {
	// 	return Color{}, err
	// }

	// object_point := shapeTransInv.MulByTuple(world_point)

	// patternTransInv, err := sp.GetTransform().Inverse()
	// if err != nil {
	// 	return Color{}, err
	// }

	// pattern_point := patternTransInv.MulByTuple(object_point)
	// return sp.PatternAt(pattern_point), nil
	return sp.PatternAtShapeOnBase(world_point, shape, sp.PatternAt)
}
