package scene

import "rayGo/calc"

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}

var Black = NewColor(0, 0, 0)
var White = NewColor(255, 255, 255)

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

func TupletoColor(t calc.Tuple4) Color {
	return NewColor(t[0], t[1], t[2])
}

//tupleとかもこっちの関数じゃなくてメソッドの方がよいかもね
func (c Color) Mul(c2 Color) Color {
	r := c.Red * c2.Red
	g := c.Green * c2.Green
	b := c.Blue * c2.Blue
	return NewColor(r, g, b)
}

func (c Color) Add(c2 Color) Color {
	r := c.Red + c2.Red
	g := c.Green + c2.Green
	b := c.Blue + c2.Blue
	return NewColor(r, g, b)
}

func (c Color) Sub(c2 Color) Color {
	r := c.Red - c2.Red
	g := c.Green - c2.Green
	b := c.Blue - c2.Blue
	return NewColor(r, g, b)
}

func (c Color) ToTuple4() calc.Tuple4 {
	return calc.Tuple4{
		c.Red,
		c.Green,
		c.Blue,
		1.0,
	}
}

func (c Color) GetByIndex(i int) float64 {
	switch i {
	case 0:
		return c.Red
	case 1:
		return c.Green
	case 2:
		return c.Blue
	default:
		return c.Red
	}
}
