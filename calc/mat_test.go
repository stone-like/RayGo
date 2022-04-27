package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMulMat(t *testing.T) {
	a := Mat4x4{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	b := Mat4x4{{-2, 1, 2, 3}, {3, 2, 1, -1}, {4, 3, 6, 5}, {1, 2, 7, 8}}

	target := Mat4x4{{20, 22, 50, 48}, {44, 54, 114, 108}, {40, 58, 110, 102}, {16, 26, 46, 42}}

	require.Equal(t, target, a.MulByMat4x4(b))
}

func TestMulMatByTuple(t *testing.T) {
	a := Mat4x4{{1, 2, 3, 4}, {2, 4, 4, 2}, {8, 6, 4, 1}, {0, 0, 0, 1}}
	b := Tuple4{1, 2, 3, 1}

	target := Tuple4{18, 24, 33, 1}

	require.Equal(t, target, a.MulByTuple(b))
}

func TestMulMatIdent(t *testing.T) {
	m := Mat4x4{{0, 1, 2, 4}, {1, 2, 4, 8}, {2, 4, 8, 16}, {4, 8, 16, 32}}
	tu := Tuple4{1, 2, 3, 4}

	ident := m.GetIdentityMat()
	require.Equal(t, m, ident.MulByMat4x4(m))
	require.Equal(t, tu, ident.MulByTuple(tu))
}

func TestTarnsposeMat(t *testing.T) {
	m := Mat4x4{{0, 9, 3, 0}, {9, 8, 0, 8}, {1, 8, 5, 3}, {0, 0, 5, 8}}

	target := Mat4x4{{0, 9, 1, 0}, {9, 8, 8, 0}, {3, 0, 5, 5}, {0, 8, 3, 8}}
	require.Equal(t, target, m.Transpose())

}

func TestSubMat(t *testing.T) {
	m3 := Mat3x3{{1, 5, 0}, {-3, 2, 7}, {0, 6, -3}}
	m4 := Mat4x4{{-6, 1, 1, 6}, {-8, 5, 8, 6}, {-1, 0, 8, 2}, {-7, 1, -1, 1}}

	target3 := Mat2x2{{-3, 2}, {0, 6}}
	target4 := Mat3x3{{-6, 1, 6}, {-8, 8, 6}, {-7, -1, 1}}

	require.Equal(t, target3, m3.SubMatrix(0, 2))
	require.Equal(t, target4, m4.SubMatrix(2, 1))

}

func TestMinor3x3(t *testing.T) {
	m3 := Mat3x3{{3, 5, 0}, {2, -1, -7}, {6, -1, 5}}

	require.Equal(t, float64(25), m3.Minor(1, 0))

}

func TestCofactor3x3(t *testing.T) {
	m3 := Mat3x3{{3, 5, 0}, {2, -1, -7}, {6, -1, 5}}

	minor1 := m3.Minor(0, 0)
	co1 := m3.Cofactor(0, 0)

	minor2 := m3.Minor(1, 0)
	co2 := m3.Cofactor(1, 0)

	require.Equal(t, float64(-12), minor1)
	require.Equal(t, float64(-12), co1)
	require.Equal(t, float64(25), minor2)
	require.Equal(t, float64(-25), co2)

}

func TestDet3x3(t *testing.T) {
	m3 := Mat3x3{{1, 2, 6}, {-5, 8, -4}, {2, 6, 4}}

	co1 := m3.Cofactor(0, 0)
	co2 := m3.Cofactor(0, 1)
	co3 := m3.Cofactor(0, 2)
	det := m3.Det()

	require.Equal(t, float64(56), co1)
	require.Equal(t, float64(12), co2)
	require.Equal(t, float64(-46), co3)
	require.Equal(t, float64(-196), det)

}

func TestDet4x4(t *testing.T) {
	m4 := Mat4x4{{-2, -8, 3, 5}, {-3, 1, 7, 3}, {1, 2, -9, 6}, {-6, 7, 7, -9}}
	co1 := m4.Cofactor(0, 0)
	co2 := m4.Cofactor(0, 1)
	co3 := m4.Cofactor(0, 2)
	co4 := m4.Cofactor(0, 3)

	det := m4.Det()

	require.Equal(t, float64(690), co1)
	require.Equal(t, float64(447), co2)
	require.Equal(t, float64(210), co3)
	require.Equal(t, float64(51), co4)

	require.Equal(t, float64(-4071), det)
}

func TestInvertMat(t *testing.T) {
	m := Mat4x4{{-5, 2, 6, -8}, {1, -5, 1, 8}, {7, 7, -6, -7}, {1, -3, 7, 4}}

	target := Mat4x4{{0.21805, 0.45113, 0.24060, -0.04511}, {-0.80827, -1.45677, -0.44361, 0.52068}, {-0.07895, -0.22368, -0.05263, 0.19737}, {-0.52256, -0.81391, -0.30075, 0.30639}}

	invert, err := m.Inverse()
	require.Nil(t, err)
	require.True(t, Mat4x4Compare(target, invert))

	m2 := Mat4x4{{8, -5, 9, 2}, {7, 5, 6, 1}, {-6, 0, 9, 6}, {-3, 0, -9, -4}}

	target2 := Mat4x4{{-0.15385, -0.15385, -0.28205, -0.53846}, {-0.07692, 0.12308, 0.02564, 0.03077}, {0.35897, 0.35897, 0.43590, 0.92308}, {-0.69231, -0.69231, -0.76923, -1.92308}}

	invert2, err := m2.Inverse()
	require.Nil(t, err)
	require.True(t, Mat4x4Compare(target2, invert2))

	m3 := Mat4x4{{9, 3, 0, 9}, {-5, -2, -6, -3}, {-4, 9, 6, 4}, {-7, 6, 6, 2}}

	target3 := Mat4x4{{-0.04074, -0.07778, 0.14444, -0.22222}, {-0.07778, 0.03333, 0.36667, -0.33333}, {-0.02901, -0.14630, -0.10926, 0.12963}, {0.17778, 0.06667, -0.26667, 0.33333}}

	invert3, err := m3.Inverse()
	require.Nil(t, err)
	require.True(t, Mat4x4Compare(target3, invert3))

}

//float64で計算しているため完全に一致する逆行列ではないがepsilon程度の誤差の逆行列ならばOK
func TestInvertCalc(t *testing.T) {
	a := Mat4x4{{3, -9, 7, 3}, {3, -8, 2, -9}, {-4, 4, 4, 1}, {-6, 5, -1, 1}}
	b := Mat4x4{{8, 2, 2, 2}, {3, -1, 7, 0}, {7, 0, 5, 4}, {6, -2, 0, 5}}

	c := a.MulByMat4x4(b)
	inbertB, _ := b.Inverse()
	ret := c.MulByMat4x4(inbertB)

	require.True(t, Mat4x4Compare(ret, a))
}
