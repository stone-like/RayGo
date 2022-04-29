package scene

import (
	"math"
	"rayGo/calc"
)

type RingPattern struct {
	*BasePattern
	Color1 Color
	Color2 Color
}

var _ Pattern = RingPattern{}

func NewRingPattern(c1, c2 Color) RingPattern {
	return RingPattern{
		NewBasePattern(),
		c1,
		c2,
	}
}

func (rp RingPattern) PatternAt(point calc.Tuple4) Color {

	pow1, pow2 := math.Pow(point[0], 2), math.Pow(point[2], 2)
	root := math.Sqrt(pow1 + pow2)

	if int(root)%2 == 0 {
		return rp.Color1
	}

	return rp.Color2

}

func (rp RingPattern) PatternAtShape(world_point calc.Tuple4, shape Shape) (Color, error) {

	return rp.PatternAtShapeOnBase(world_point, shape, rp.PatternAt)
}
