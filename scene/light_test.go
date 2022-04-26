package scene

import (
	"math"
	"rayGo/calc"
	"testing"
)

func BuildLightTestSturct() []struct {
	normal_vec calc.Tuple4
	eye_vec    calc.Tuple4
	light      Light
	ans        Color
} {
	return []struct {
		normal_vec calc.Tuple4
		eye_vec    calc.Tuple4
		light      Light
		ans        Color
	}{
		{
			calc.NewVector(0, 0, -1),
			calc.NewVector(0, 0, -1),
			NewLight(calc.NewPoint(0, 0, -10), NewColor(1, 1, 1)),
			NewColor(1.9, 1.9, 1.9),
		},
		{
			calc.NewVector(0, 0, -1),
			calc.NewVector(0, math.Sqrt(2)/2, -math.Sqrt(2)/2),
			NewLight(calc.NewPoint(0, 0, -10), NewColor(1, 1, 1)),
			NewColor(1.0, 1.0, 1.0),
		},
		{
			calc.NewVector(0, 0, -1),
			calc.NewVector(0, 0, -1),
			NewLight(calc.NewPoint(0, 10, -10), NewColor(1, 1, 1)),
			NewColor(0.7364, 0.7364, 0.7364),
		},
		{
			calc.NewVector(0, 0, -1),
			calc.NewVector(0, -math.Sqrt(2)/2, -math.Sqrt(2)/2),
			NewLight(calc.NewPoint(0, 10, -10), NewColor(1, 1, 1)),
			NewColor(1.6364, 1.6364, 1.6364),
		},
		{
			calc.NewVector(0, 0, -1),
			calc.NewVector(0, 0, -1),
			NewLight(calc.NewPoint(0, 0, 10), NewColor(1, 1, 1)),
			NewColor(0.1, 0.1, 0.1),
		},
	}
}

func TestLighting(t *testing.T) {
	m := DefaultMaterial()
	pos := calc.NewPoint(0, 0, 0)

	for _, target := range BuildLightTestSturct() {

		res := target.light.Lighting(m, pos, target.eye_vec, target.normal_vec)

		ret := colorCompare(target.ans, res)
		if ret == false {
			t.Error(target.ans, res)
		}
	}

}
