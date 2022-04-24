package scene

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
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
