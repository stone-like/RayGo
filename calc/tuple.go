package calc

import "math"

type Tuple4 [4]float64

func (t Tuple4) Magnitude() float64 {

	content := math.Pow(t[0], 2) + math.Pow(t[1], 2) + math.Pow(t[2], 2) + math.Pow(t[3], 2)
	return math.Sqrt(content)
}

func (t Tuple4) Normalize() Tuple4 {
	mag := t.Magnitude()
	return Tuple4{
		t[0] / mag,
		t[1] / mag,
		t[2] / mag,
		t[3] / mag,
	}
}

//vectorは第4成分が0
func NewVector(x, y, z float64) Tuple4 {
	return Tuple4{x, y, z, 0}
}

//pointは第4成分が1
func NewPoint(x, y, z float64) Tuple4 {
	return Tuple4{x, y, z, 1}
}

//tuple4自体にaddとか入れるのも一案だけど一旦関数で

//point+pointは第四成分が2になってしまうので未定義にした方がよい？
func AddTuple(a, b Tuple4) Tuple4 {
	//if a[3]+b[3] == 2 return emptyTuple4 or error
	return Tuple4{a[0] + b[0], a[1] + b[1], a[2] + b[2], a[3] + b[3]}
}

//vector - pointは第四成分が-1となるので未定義
func SubTuple(a, b Tuple4) Tuple4 {

	return Tuple4{a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3]}
}

func NegTuple(a Tuple4) Tuple4 {
	return Tuple4{
		-a[0],
		-a[1],
		-a[2],
		-a[3],
	}
}

func MulTupleByScalar(s float64, a Tuple4) Tuple4 {
	return Tuple4{
		a[0] * s,
		a[1] * s,
		a[2] * s,
		a[3] * s,
	}
}

func DivTupleByScalar(s float64, a Tuple4) Tuple4 {
	return Tuple4{
		a[0] / s,
		a[1] / s,
		a[2] / s,
		a[3] / s,
	}
}

func DotTuple(a, b Tuple4) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

//クロス積はa,bの順番が肝要
func CrossTuple(a, b Tuple4) Tuple4 {

	newX := a[1]*b[2] - a[2]*b[1]
	newY := a[2]*b[0] - a[0]*b[2]
	newZ := a[0]*b[1] - a[1]*b[0]
	return NewVector(
		newX, newY, newZ,
	)
}
