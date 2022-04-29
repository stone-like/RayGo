package scene

import (
	"math"
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
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
	s := NewSphere(1)
	m := DefaultMaterial()
	s.SetMaterial(m)
	pos := calc.NewPoint(0, 0, 0)
	in_shadow := false
	for _, target := range BuildLightTestSturct() {

		res, err := target.light.Lighting(m, pos, target.eye_vec, target.normal_vec, in_shadow, s)
		require.Nil(t, err)

		ret := colorCompare(target.ans, res)
		if ret == false {
			t.Error(target.ans, res)
		}
	}

}

func Test_Ignore_Diffuse_And_Specular_When_in_Shadow(t *testing.T) {
	eye_vec := calc.NewVector(0, 0, -1)
	normal_vec := calc.NewVector(0, 0, -1)
	light := NewLight(calc.NewPoint(0, 0, -10), NewColor(1, 1, 1))
	in_shadow := true

	s := NewSphere(1)
	m := DefaultMaterial()
	s.SetMaterial(m)

	pos := calc.NewPoint(0, 0, 0)
	res, err := light.Lighting(m, pos, eye_vec, normal_vec, in_shadow, s)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.1, 0.1, 0.1), res))

}

func Test_Lighting_with_Pattern(t *testing.T) {
	m := DefaultMaterial()
	m.SetPattern(NewStripePattern(White, Black))
	m.Ambient = 1
	m.Diffuse = 0
	m.Specular = 0

	s := NewSphere(1)
	s.SetMaterial(m)

	eye_vec := calc.NewVector(0, 0, -1)
	normal_vec := calc.NewVector(0, 0, -1)

	light := NewLight(calc.NewPoint(0, 0, -10), NewColor(1, 1, 1))

	c1, err := light.Lighting(m, calc.NewPoint(0.9, 0, 0), eye_vec, normal_vec, false, s)
	require.Nil(t, err)

	c2, err := light.Lighting(m, calc.NewPoint(1.1, 0, 0), eye_vec, normal_vec, false, s)
	require.Nil(t, err)

	require.True(t, colorCompare(White, c1))
	require.True(t, colorCompare(Black, c2))

}
