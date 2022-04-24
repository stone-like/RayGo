package calc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	a1 := NewPoint(3, -2, 5)
	a2 := NewVector(-2, 3, 1)

	add := AddTuple(a1, a2)

	require.Equal(t, add, Tuple4{1, 1, 6, 1})
}

func TestSub(t *testing.T) {
	for _, target := range []struct {
		a   Tuple4
		b   Tuple4
		ans Tuple4
	}{
		{
			NewPoint(3, 2, 1),
			NewPoint(5, 6, 7),
			NewVector(-2, -4, -6),
		},
		{
			NewPoint(3, 2, 1),
			NewVector(5, 6, 7),
			NewPoint(-2, -4, -6),
		},
		{
			NewVector(3, 2, 1),
			NewVector(5, 6, 7),
			NewVector(-2, -4, -6),
		},
	} {
		require.Equal(t, SubTuple(target.a, target.b), target.ans)
	}
}

func TestNegate(t *testing.T) {
	tu := Tuple4{1, -2, 3, -4}

	require.Equal(t, NegTuple(tu), Tuple4{-1, 2, -3, 4})
}

func TestMul(t *testing.T) {
	for _, target := range []struct {
		tup Tuple4
		mul float64
		ans Tuple4
	}{
		{
			Tuple4{1, -2, 3, -4},
			3.5,
			Tuple4{3.5, -7, 10.5, -14},
		},
		{
			Tuple4{1, -2, 3, -4},
			0.5,
			Tuple4{0.5, -1, 1.5, -2},
		},
	} {
		require.Equal(t, MulTupleByScalar(target.mul, target.tup), target.ans)
	}
}

func TestDiv(t *testing.T) {
	tup := DivTupleByScalar(2, Tuple4{1, -2, 3, -4})

	require.Equal(t, tup, Tuple4{0.5, -1, 1.5, -2})
}

func TestMagnitude(t *testing.T) {
	for _, target := range []struct {
		target Tuple4
		ans    float64
	}{
		{
			NewVector(1, 0, 0),
			1,
		},
		{
			NewVector(1, 2, 3),
			math.Sqrt(14),
		},
		{
			NewVector(-1, -2, -3),
			math.Sqrt(14),
		},
	} {
		require.Equal(t, target.target.Magnitude(), target.ans)
	}
}

func TestNormalize(t *testing.T) {
	v := NewVector(1, 2, 3)
	m := math.Sqrt(14)
	target := NewVector(1/m, 2/m, 3/m)

	n := v.Normalize()
	require.Equal(t, n, target)
	require.Equal(t, n.Magnitude(), float64(1))
}

func TestDot(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)
	require.Equal(t, DotTuple(a, b), float64(20))
}

func TestCross(t *testing.T) {
	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)
	require.Equal(t, NewVector(-1, 2, -1), CrossTuple(a, b))
	require.Equal(t, NewVector(1, -2, 1), CrossTuple(b, a))

}
