package calc

import "math"

//ひょっとしたらMat4x4を作っていくよりもPointとVectorにTransformを実装していくのがいいかも
func NewTranslation(x, y, z float64) Mat4x4 {
	return Mat4x4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func NewScale(x, y, z float64) Mat4x4 {
	return Mat4x4{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func NewRotateX(angle float64) Mat4x4 {
	return Mat4x4{
		{1, 0, 0, 0},
		{0, math.Cos(angle), -math.Sin(angle), 0},
		{0, math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}
func NewRotateY(angle float64) Mat4x4 {
	return Mat4x4{
		{math.Cos(angle), 0, math.Sin(angle), 0},
		{0, 1, 0, 0},
		{-math.Sin(angle), 0, math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}
func NewRotateZ(angle float64) Mat4x4 {
	return Mat4x4{
		{math.Cos(angle), -math.Sin(angle), 0, 0},
		{math.Sin(angle), math.Cos(angle), 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

//shearingはx,y,zがそれぞれ残りの二つに対して変化する
//x in proportion to y
//x in proportion to z
//...
//といった感じ、なのでパラメータが６つ
func NewShearing(xtoY, xToZ, yToX, yToZ, zToX, zToY float64) Mat4x4 {
	return Mat4x4{
		{1, xtoY, xToZ, 0},
		{yToX, 1, yToZ, 0},
		{zToX, zToY, 1, 0},
		{0, 0, 0, 1},
	}
}
