package scene

import (
	"math"
	"rayGo/calc"
)

type GradientPattern struct {
	*BasePattern
	Start Color
	End   Color
}

var _ Pattern = GradientPattern{}

func NewGradientPattern(c1, c2 Color) GradientPattern {
	return GradientPattern{
		NewBasePattern(),
		c1,
		c2,
	}
}

func (gp GradientPattern) PatternAt(point calc.Tuple4) Color {
	distance := gp.End.Sub(gp.Start).ToTuple4()
	fraction := point[0] - math.Floor(point[0])

	return gp.Start.Add(TupletoColor(calc.MulTupleByScalar(fraction, distance)))
}

func (gp GradientPattern) PatternAtShape(world_point calc.Tuple4, shape Shape) (Color, error) {

	return gp.PatternAtShapeOnBase(world_point, shape, gp.PatternAt)
}
