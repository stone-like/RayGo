package calc

import (
	"math"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranslation(t *testing.T) {
	trans := NewTranslation(5, -3, 2)
	p := NewPoint(-3, 4, 5)

	target := NewPoint(2, 1, 7)
	require.Equal(t, target, trans.MulByTuple(p))
}

func TestInverseTranslation(t *testing.T) {
	trans := NewTranslation(5, -3, 2)
	inv, _ := trans.Inverse()
	p := NewPoint(-3, 4, 5)

	target := NewPoint(-8, 7, 3)
	require.Equal(t, target, inv.MulByTuple(p))
}

//Vectorは第四要素が0なのでtraslationの影響がない
func Test_Translation_Does_Not_Affect_Vector(t *testing.T) {
	trans := NewTranslation(5, -3, 2)
	p := NewVector(-3, 4, 5)
	require.Equal(t, p, trans.MulByTuple(p))

}

func Test_Scailing(t *testing.T) {
	trans := NewScale(2, 3, 4)
	p := NewPoint(-4, 6, 8)
	target := NewPoint(-8, 18, 32)
	require.Equal(t, target, trans.MulByTuple(p))
}

func Test_Scailing_Inverse(t *testing.T) {
	trans := NewScale(2, 3, 4)
	inv, _ := trans.Inverse()
	p := NewPoint(-4, 6, 8)
	target := NewPoint(-2, 2, 2)
	require.Equal(t, target, inv.MulByTuple(p))
}

//scalingはvectorにも作用
func Test_Scailing_On_Vector(t *testing.T) {
	trans := NewScale(2, 3, 4)
	p := NewVector(-4, 6, 8)
	target := NewVector(-8, 18, 32)
	require.Equal(t, target, trans.MulByTuple(p))
}

func Test_X_dir_Reflection_By_Scailing(t *testing.T) {
	trans := NewScale(-1, 1, 1)
	p := NewPoint(2, 3, 4)
	target := NewPoint(-2, 3, 4)
	require.Equal(t, target, trans.MulByTuple(p))
}

func TupleCompare(a, b Tuple4) bool {
	for i := 0; i < 4; i++ {
		ret := util.FloatEqual(a[i], b[i])
		if ret == false {
			return false
		}
	}

	return true
}

func Test_Rotate_X(t *testing.T) {
	p := NewPoint(0, 1, 0)
	half_quarter := NewRotateX(math.Pi / 4)
	full_quarter := NewRotateX(math.Pi / 2)

	halfTarget := NewPoint(0, math.Sqrt(2)/2, math.Sqrt(2)/2)
	fullTarget := NewPoint(0, 0, 1)

	require.True(t, TupleCompare(halfTarget, half_quarter.MulByTuple(p)))
	require.True(t, TupleCompare(fullTarget, full_quarter.MulByTuple(p)))

}

func Test_Rotate_X_Inverse(t *testing.T) {
	p := NewPoint(0, 1, 0)
	half_quarter := NewRotateX(math.Pi / 4)
	inv, _ := half_quarter.Inverse()

	halfInvTarget := NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2)

	require.True(t, TupleCompare(halfInvTarget, inv.MulByTuple(p)))
}

func Test_Rotate_Y(t *testing.T) {
	p := NewPoint(0, 0, 1)
	half_quarter := NewRotateY(math.Pi / 4)
	full_quarter := NewRotateY(math.Pi / 2)

	halfTarget := NewPoint(math.Sqrt(2)/2, 0, math.Sqrt(2)/2)
	fullTarget := NewPoint(1, 0, 0)

	require.True(t, TupleCompare(halfTarget, half_quarter.MulByTuple(p)))
	require.True(t, TupleCompare(fullTarget, full_quarter.MulByTuple(p)))

}

func Test_Rotate_Z(t *testing.T) {
	p := NewPoint(0, 1, 0)
	half_quarter := NewRotateZ(math.Pi / 4)
	full_quarter := NewRotateZ(math.Pi / 2)

	halfTarget := NewPoint(-math.Sqrt(2)/2, math.Sqrt(2)/2, 0)
	fullTarget := NewPoint(-1, 0, 0)

	require.True(t, TupleCompare(halfTarget, half_quarter.MulByTuple(p)))
	require.True(t, TupleCompare(fullTarget, full_quarter.MulByTuple(p)))

}

func TestShearing(t *testing.T) {
	for _, test := range []struct {
		trans Mat4x4
		p     Tuple4
		ans   Tuple4
	}{
		{
			NewShearing(0, 1, 0, 0, 0, 0),
			NewPoint(2, 3, 4),
			NewPoint(6, 3, 4),
		},
		{
			NewShearing(0, 0, 1, 0, 0, 0),
			NewPoint(2, 3, 4),
			NewPoint(2, 5, 4),
		},
		{
			NewShearing(0, 0, 0, 1, 0, 0),
			NewPoint(2, 3, 4),
			NewPoint(2, 7, 4),
		},
		{
			NewShearing(0, 0, 0, 0, 1, 0),
			NewPoint(2, 3, 4),
			NewPoint(2, 3, 6),
		},
		{
			NewShearing(0, 0, 0, 0, 0, 1),
			NewPoint(2, 3, 4),
			NewPoint(2, 3, 7),
		},
	} {
		require.True(t, TupleCompare(test.ans, test.trans.MulByTuple(test.p)))
	}
}

//matMulはassociativeは保つ A*(B*C) = (A*B)*C
//だけどcomutitveは保たない A*B != B*A

func TestChainTransform(t *testing.T) {
	p := NewPoint(1, 0, 1)
	a := NewRotateX(math.Pi / 2)
	b := NewScale(5, 5, 5)
	c := NewTranslation(10, 5, 7)

	//a->b->cの順に作用させたいなら,
	//C*B*AとMULする
	mat := c.MulByMat4x4(b).MulByMat4x4(a)
	mat2 := b.MulByMat4x4(c).MulByMat4x4(a)
	mat3 := a.MulByMat4x4(b).MulByMat4x4(c)

	target := NewPoint(15, 0, 7)

	require.True(t, TupleCompare(target, mat.MulByTuple(p)))
	require.False(t, TupleCompare(target, mat2.MulByTuple(p)))
	require.False(t, TupleCompare(target, mat3.MulByTuple(p)))

}
