package scene

import (
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Intersection_Can_Contain_U_And_V(t *testing.T) {
	tri := NewSmoothTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0), calc.NewVector(0, 1, 0), calc.NewVector(-1, 0, 0), calc.NewVector(1, 0, 0))

	ray := NewRay(calc.NewPoint(-0.2, 0.3, -2), calc.NewVector(0, 0, 1))

	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 1, xs.Count)
	require.True(t, util.FloatEqual(0.45, xs.Intersections[0].U))
	require.True(t, util.FloatEqual(0.25, xs.Intersections[0].V))

}

func Test_Use_U_And_V_To_Interpolate_The_Normal(t *testing.T) {
	tri := NewSmoothTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0), calc.NewVector(0, 1, 0), calc.NewVector(-1, 0, 0), calc.NewVector(1, 0, 0))

	ray := NewRay(calc.NewPoint(-0.2, 0.3, -2), calc.NewVector(0, 0, 1))

	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)

	n, err := tri.NormalAt(calc.NewPoint(0, 0, 0), *xs.Intersections[0])
	require.Nil(t, err)

	require.True(t, calc.TupleCompare(calc.NewVector(-0.5547, 0.83205, 0), n))
}

func Test_Passing_The_Normal_on_Smooth_Triangle(t *testing.T) {
	tri := NewSmoothTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0), calc.NewVector(0, 1, 0), calc.NewVector(-1, 0, 0), calc.NewVector(1, 0, 0))

	ray := NewRay(calc.NewPoint(-0.2, 0.3, -2), calc.NewVector(0, 0, 1))

	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)

	comp, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)

	require.True(t, calc.TupleCompare(calc.NewVector(-0.5547, 0.83205, 0), comp.NormalVec))
}
