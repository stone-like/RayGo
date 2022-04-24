package calc

import "math"

type Tuple4 struct {
	x float64
	y float64
	z float64
	w float64
}

func (t Tuple4) Magnitude() float64 {

	content := math.Pow(t.x, 2) + math.Pow(t.y, 2) + math.Pow(t.z, 2) + math.Pow(t.w, 2)
	return math.Sqrt(content)
}

func (t Tuple4) Normalize() Tuple4 {
	mag := t.Magnitude()
	return Tuple4{
		t.x / mag,
		t.y / mag,
		t.z / mag,
		t.w / mag,
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
	return Tuple4{a.x + b.x, a.y + b.y, a.z + b.z, a.w + b.w}
}

//vector - pointは第四成分が-1となるので未定義
func SubTuple(a, b Tuple4) Tuple4 {

	return Tuple4{a.x - b.x, a.y - b.y, a.z - b.z, a.w - b.w}
}

func NegTuple(a Tuple4) Tuple4 {
	return Tuple4{
		-a.x,
		-a.y,
		-a.z,
		-a.w,
	}
}

func MulTupleByScalar(s float64, a Tuple4) Tuple4 {
	return Tuple4{
		a.x * s,
		a.y * s,
		a.z * s,
		a.w * s,
	}
}

func DivTupleByScalar(s float64, a Tuple4) Tuple4 {
	return Tuple4{
		a.x / s,
		a.y / s,
		a.z / s,
		a.w / s,
	}
}

func DotTuple(a, b Tuple4) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z + a.w*b.w
}

//クロス積はa,bの順番が肝要
func CrossTuple(a, b Tuple4) Tuple4 {

	newX := a.y*b.z - a.z*b.y
	newY := a.z*b.x - a.x*b.z
	newZ := a.x*b.y - a.y*b.x
	return NewVector(
		newX, newY, newZ,
	)
}
