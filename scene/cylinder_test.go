package scene

import (
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ray_Missing_Cylinder(t *testing.T) {
	c := NewCyliner()

	for _, target := range []struct {
		title string
		point calc.Tuple4
		dir   calc.Tuple4
	}{
		{
			title: "1",
			point: calc.NewPoint(1, 0, 0),
			dir:   calc.NewVector(0, 1, 0),
		},
		{
			title: "2",
			point: calc.NewPoint(0, 0, 0),
			dir:   calc.NewVector(0, 1, 0),
		},
		{
			title: "3",
			point: calc.NewPoint(0, 0, -5),
			dir:   calc.NewVector(1, 1, 1),
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			ray := NewRay(target.point, target.dir.Normalize())

			xs, err := c.calcLocalIntersect(ray)
			require.Nil(t, err)
			require.Equal(t, 0, xs.Count)
		})
	}
}

func Test_Ray_Hits_Cylinder(t *testing.T) {
	c := NewCyliner()

	for _, target := range []struct {
		title string
		ray   Ray
		t0    float64
		t1    float64
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(1, 0, -5), calc.NewVector(0, 0, 1).Normalize()),
			t0:    5,
			t1:    5,
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1).Normalize()),
			t0:    4,
			t1:    6,
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(0.5, 0, -5), calc.NewVector(0.1, 1, 1).Normalize()),
			t0:    6.80798,
			t1:    7.08872,
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

func Test_Normal_On_Cylinder(t *testing.T) {
	c := NewCyliner()

	for _, target := range []struct {
		title  string
		point  calc.Tuple4
		normal calc.Tuple4
	}{
		{
			title:  "1",
			point:  calc.NewPoint(1, 0, 0),
			normal: calc.NewVector(1, 0, 0),
		},
		{
			title:  "2",
			point:  calc.NewPoint(0, 5, -1),
			normal: calc.NewVector(0, 0, -1),
		},
		{
			title:  "3",
			point:  calc.NewPoint(-1, 1, 0),
			normal: calc.NewVector(-1, 0, 0),
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			normal := c.calcLocalNormal(target.point)
			require.True(t, calc.TupleCompare(target.normal, normal))
		})
	}
}

func Test_Truncated_Cynlinder(t *testing.T) {
	c := NewCyliner(CynMin(1), CynMax(2))
	for _, target := range []struct {
		title string
		ray   Ray
		count int
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(0, 1.5, 0), calc.NewVector(0.1, 1, 0).Normalize()),
			count: 0,
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, 3, -5), calc.NewVector(0, 0, 1).Normalize()),
			count: 0,
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1).Normalize()),
			count: 0,
		},
		{
			title: "4",
			ray:   NewRay(calc.NewPoint(0, 2, -5), calc.NewVector(0, 0, 1).Normalize()),
			count: 0,
		},
		{
			title: "5",
			ray:   NewRay(calc.NewPoint(0, 1, -5), calc.NewVector(0, 0, 1).Normalize()),
			count: 0,
		},
		{
			title: "6",
			ray:   NewRay(calc.NewPoint(0, 1.5, -2), calc.NewVector(0, 0, 1).Normalize()),
			count: 2,
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, target.count, xs.Count)

		})
	}

}

//Cappedはシリンダーの上面と下面を閉じたもの

func Test_Intersect_capped_cylinder(t *testing.T) {
	c := NewCyliner(CynMin(1), CynMax(2), CynClosed(true))

	for _, target := range []struct {
		title string
		ray   Ray
		count int
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(0, 3, 0), calc.NewVector(0, -1, 0).Normalize()),
			count: 2,
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, 3, -2), calc.NewVector(0, -1, 2).Normalize()),
			count: 2,
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(0, 4, -2), calc.NewVector(0, -1, 1).Normalize()),
			count: 2,
		},
		{
			title: "4",
			ray:   NewRay(calc.NewPoint(0, 0, -2), calc.NewVector(0, 1, 2).Normalize()),
			count: 2,
		},
		{
			title: "5",
			ray:   NewRay(calc.NewPoint(0, -1, -2), calc.NewVector(0, 1, 1).Normalize()),
			count: 2,
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, target.count, xs.Count)

		})
	}

}

func Test_Normal_On_capped_cylinder(t *testing.T) {
	c := NewCyliner(CynMin(1), CynMax(2), CynClosed(true))

	for _, target := range []struct {
		title  string
		point  calc.Tuple4
		normal calc.Tuple4
	}{
		{
			title:  "1",
			point:  calc.NewPoint(0, 1, 0),
			normal: calc.NewVector(0, -1, 0),
		},
		{
			title:  "2",
			point:  calc.NewPoint(0.5, 1, 0),
			normal: calc.NewVector(0, -1, 0),
		},
		{
			title:  "3",
			point:  calc.NewPoint(0, 1, 0.5),
			normal: calc.NewVector(0, -1, 0),
		},
		{
			title:  "4",
			point:  calc.NewPoint(0, 2, 0),
			normal: calc.NewVector(0, 1, 0),
		},
		{
			title:  "5",
			point:  calc.NewPoint(0.5, 2, 0),
			normal: calc.NewVector(0, 1, 0),
		},
		{
			title:  "6",
			point:  calc.NewPoint(0, 2, 0.5),
			normal: calc.NewVector(0, 1, 0),
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			normal := c.calcLocalNormal(target.point)
			require.True(t, calc.TupleCompare(target.normal, normal))
		})
	}

}
