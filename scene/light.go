package scene

import (
	"math"
	"rayGo/calc"
)

type Light struct {
	Position  calc.Tuple4
	Intensity Color
}

func NewLight(p calc.Tuple4, color Color) Light {
	return Light{
		Position:  p,
		Intensity: color,
	}
}

//shapeのmaterial,position
//eyeVec(rayのnegate)
//shapeのpositionに対してのnormalVec
//light.Positionは光源の位置、shape.Positionは図形の位置

//Borrow from Phone Reflection Model
//さすがにtuple絡みは関数じゃなくてメソッドを使ってChainさせた方が良さそうな感じ
func (l *Light) Lighting(m Material, position, eye_vec, normal_vec calc.Tuple4) Color {

	effective_color := m.Color.Mul(l.Intensity).ToTuple4()

	light_vec := calc.SubTuple(l.Position, position).Normalize()

	ambient := TupletoColor(calc.MulTupleByScalar(m.Ambient, effective_color))

	light_dot_normal := calc.DotTuple(light_vec, normal_vec)

	var diffuse Color
	var specular Color

	if light_dot_normal < 0 {
		diffuse = Black
		specular = Black

		return ambient.Add(diffuse).Add(specular)
	}

	diffuse = TupletoColor(calc.MulTupleByScalar(light_dot_normal, calc.MulTupleByScalar(m.Diffuse, effective_color)))

	reflect_vec := calc.Reflect(calc.NegTuple(light_vec), normal_vec)
	refelect_dot_eye := calc.DotTuple(reflect_vec, eye_vec)

	if refelect_dot_eye <= 0 {
		specular = Black
	} else {
		factor := math.Pow(refelect_dot_eye, m.Shininess)
		specular = TupletoColor(calc.MulTupleByScalar(factor, calc.MulTupleByScalar(m.Specular, l.Intensity.ToTuple4())))
	}

	return ambient.Add(diffuse).Add(specular)
}
