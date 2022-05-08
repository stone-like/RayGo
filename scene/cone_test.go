package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Intersect_Cone(t *testing.T) {
	c := NewCone()

	for _, target := range []struct {
		title string
		ray   Ray
		t0    float64
		t1    float64
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1).Normalize()),
			t0:    5,
			t1:    5,
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(1, 1, 1).Normalize()),
			t0:    8.66025,
			t1:    8.66025,
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(1, 1, -5), calc.NewVector(-0.5, -1, 1).Normalize()),
			t0:    4.55006,
			t1:    49.44994,
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, 2, xs.Count)
			require.True(t, util.FloatEqual(target.t0, xs.Intersections[0].Time))
			require.True(t, util.FloatEqual(target.t1, xs.Intersections[1].Time))

		})
	}
}

func Test_Intersect_Cone_When_Ray_Parallel_To_One_of_its_Halves(t *testing.T) {
	c := NewCone()
	ray := NewRay(calc.NewPoint(0, 0, -1), calc.NewVector(0, 1, 1).Normalize())

	xs, err := c.calcLocalIntersect(ray)
	require.Nil(t, err)

	require.Equal(t, 1, xs.Count)
	require.True(t, util.FloatEqual(0.35355, xs.Intersections[0].Time))
}

func Test_Intersect_Cone_End_Cap(t *testing.T) {

	c := NewCone(ConeMin(-0.5), ConeMax(0.5), ConeClosed(true))
	for _, target := range []struct {
		title string
		ray   Ray
		count int
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 1, 0).Normalize()),
			count: 0,
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, 0, -0.25), calc.NewVector(0, 1, 1).Normalize()),
			count: 2,
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(0, 0, -0.25), calc.NewVector(0, 1, 0).Normalize()),
			count: 4,
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, target.count, xs.Count)

		})
	}
}

func Test_Normal_On_Cone(t *testing.T) {
	c := NewCone()

	for _, target := range []struct {
		title  string
		point  calc.Tuple4
		normal calc.Tuple4
	}{
		{
			title:  "1",
			point:  calc.NewPoint(0, 0, 0),
			normal: calc.NewVector(0, 0, 0),
		},
		{
			title:  "2",
			point:  calc.NewPoint(1, 1, 1),
			normal: calc.NewVector(1, -math.Sqrt(2), 1),
		},
		{
			title:  "3",
			point:  calc.NewPoint(-1, -1, 0),
			normal: calc.NewVector(-1, 1, 0),
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			normal := c.calcLocalNormal(target.point)
			require.True(t, calc.TupleCompare(target.normal, normal))
		})
	}
}
