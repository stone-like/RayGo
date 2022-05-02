package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Schlick_Approximation_Under_Total_Internal_Reflection(t *testing.T) {
	shape := GlassSphere()
	ray := NewRay(calc.NewPoint(0, 0, math.Sqrt(2)), calc.NewVector(0, 1, 0))

	xs := AggregateIntersection(
		&Intersection{
			-math.Sqrt(2) / 2,
			shape,
		},
		&Intersection{
			math.Sqrt(2) / 2,
			shape,
		},
	)

	comps, err := PrepareComputations(*xs.Intersections[1], ray, xs)
	require.Nil(t, err)
	require.True(t, util.FloatEqual(1, comps.ComputeSchlick()))

}

func Test_Schlick_Approximation_Under_Perpendicular_Angle(t *testing.T) {
	shape := GlassSphere()
	ray := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 1, 0))

	xs := AggregateIntersection(
		&Intersection{
			-1,
			shape,
		},
		&Intersection{
			1,
			shape,
		},
	)

	comps, err := PrepareComputations(*xs.Intersections[1], ray, xs)
	require.Nil(t, err)
	require.True(t, util.FloatEqual(0.04, comps.ComputeSchlick()))

}

func Test_Schlick_Approximation_Under_Small_Angle_And_N2_Greater_Than_N1(t *testing.T) {
	shape := GlassSphere()
	ray := NewRay(calc.NewPoint(0, 0.99, -2), calc.NewVector(0, 0, 1))

	xs := AggregateIntersection(
		&Intersection{
			1.8589,
			shape,
		},
	)

	comps, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)
	require.True(t, util.FloatEqual(0.48873, comps.ComputeSchlick()))

}
