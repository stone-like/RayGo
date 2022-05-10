package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

//TODO:UVマッピングを実装してSphereへの貼り付けを完全なものにする
type CheckersPattern struct {
	*BasePattern
	Color1 Color
	Color2 Color
}

var _ Pattern = CheckersPattern{}

func NewCheckersPattern(c1, c2 Color) CheckersPattern {
	return CheckersPattern{
		NewBasePattern(),
		c1,
		c2,
	}
}

func (cp CheckersPattern) PatternAt(point calc.Tuple4) Color {
	var x float64
	var y float64
	var z float64

	if math.Abs(point[0]) > util.DefaultEpsilon {
		x = point[0]
	}
	if math.Abs(point[1]) > util.DefaultEpsilon {
		y = point[1]
	}
	if math.Abs(point[2]) > util.DefaultEpsilon {
		z = point[2]
	}

	added := math.Floor(x) + math.Floor(y) + math.Floor(z)

	if int(math.Round(added))%2 == 0 {
		return cp.Color1
	}

	return cp.Color2
}

func (cp CheckersPattern) PatternAtShape(world_point calc.Tuple4, shape Shape) (Color, error) {

	return cp.PatternAtShapeOnBase(world_point, shape, cp.PatternAt)
}
