package calc

import (
	"errors"
	"math"
	"rayGo/util"
)

type Mat4x4 [4][4]float64

var Ident4x4 = [4][4]float64{
	{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1},
}

func (m Mat4x4) MulByMat4x4(m2 Mat4x4) Mat4x4 {

	var temp Mat4x4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			temp[i][j] = m[i][0]*m2[0][j] + m[i][1]*m2[1][j] + m[i][2]*m2[2][j] + m[i][3]*m2[3][j]
		}
	}
	return temp
}

//4x4 * 4x1 = 4x1
func (m Mat4x4) MulByTuple(t Tuple4) Tuple4 {

	var temp Tuple4
	for i := 0; i < 4; i++ {
		temp[i] = m[i][0]*t[0] + m[i][1]*t[1] + m[i][2]*t[2] + m[i][3]*t[3]
	}
	return temp
}

func (m Mat4x4) GetIdentityMat() Mat4x4 {
	return Ident4x4
}

func (m Mat4x4) Transpose() Mat4x4 {

	var temp Mat4x4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				temp[i][j] = m[i][j]
			}

			temp[i][j] = m[j][i]
		}
	}

	return temp

}

func (m Mat4x4) SubMatrix(y, x int) Mat3x3 {
	var temp Mat3x3

	countY := 0
	for i := 0; i < 4; i++ {
		if i == y {
			continue
		}
		countX := 0

		for j := 0; j < 4; j++ {

			if j == x {
				continue
			}

			temp[countY][countX] = m[i][j]

			countX++
		}

		countY++
	}

	return temp
}

func (m Mat4x4) Minor(y, x int) float64 {
	m3x3 := m.SubMatrix(y, x)
	return m3x3.Det()
}
func (m Mat4x4) Cofactor(y, x int) float64 {
	return getSign(y, x) * m.Minor(y, x)
}

func (m Mat4x4) Det() float64 {
	det := 0.0
	for i := 0; i < 4; i++ {
		det += m[0][i] * m.Cofactor(0, i)
	}
	return det
}

func (m Mat4x4) IsInvertible() bool {
	if m.Det() == 0 {
		return false
	}

	return true
}

func (m Mat4x4) Inverse() (Mat4x4, error) {
	if !m.IsInvertible() {
		return Ident4x4, errors.New("this matrix is not Invertible")
	}

	var temp Mat4x4

	det := m.Det()

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			co := m.Cofactor(i, j)
			temp[j][i] = co / det
		}
	}

	return temp, nil
}

type Mat3x3 [3][3]float64

func (m Mat3x3) SubMatrix(y, x int) Mat2x2 {

	var temp Mat2x2

	countY := 0
	for i := 0; i < 3; i++ {
		if i == y {
			continue
		}
		countX := 0

		for j := 0; j < 3; j++ {

			if j == x {
				continue
			}

			temp[countY][countX] = m[i][j]

			countX++
		}

		countY++
	}

	return temp
}

func (m Mat3x3) Minor(y, x int) float64 {
	m2x2 := m.SubMatrix(y, x)
	return m2x2.Det()
}

func (m Mat3x3) Cofactor(y, x int) float64 {
	return getSign(y, x) * m.Minor(y, x)
}

//もしもっとfloatを足しこんでかつ、ずれが許容されないならdecimalを使う方がよい
func (m Mat3x3) Det() float64 {
	det := 0.0
	for i := 0; i < 3; i++ {
		det += m[0][i] * m.Cofactor(0, i)
	}
	return det
}

func getSign(y, x int) float64 {
	return math.Pow(-1, float64(y+x))
}

type Mat2x2 [2][2]float64

func (m Mat2x2) Det() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

func Mat4x4Compare(a, b Mat4x4) bool {

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			ret := util.FloatEqual(a[i][j], b[i][j])

			if ret == false {
				return false
			}
		}
	}
	return true
}

func Mat3x3Compare(a, b Mat3x3) bool {

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			ret := util.FloatEqual(a[i][j], b[i][j])

			if ret == false {
				return false
			}
		}
	}
	return true
}

func Mat2x2Compare(a, b Mat2x2) bool {

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			ret := util.FloatEqual(a[i][j], b[i][j])

			if ret == false {
				return false
			}
		}
	}
	return true
}
