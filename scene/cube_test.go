package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCubeIntersect(t *testing.T) {
	c := NewCube()

	for _, target := range []struct {
		title string
		ray   Ray
		t1    float64
		t2    float64
	}{
		{
			title: "+x",
			ray:   NewRay(calc.NewPoint(5, 0.5, 0), calc.NewVector(-1, 0, 0)),
			t1:    4,
			t2:    6,
		},
		{
			title: "-x",
			ray:   NewRay(calc.NewPoint(-5, 0.5, 0), calc.NewVector(1, 0, 0)),
			t1:    4,
			t2:    6,
		},
		{
			title: "+y",
			ray:   NewRay(calc.NewPoint(0.5, 5, 0), calc.NewVector(0, -1, 0)),
			t1:    4,
			t2:    6,
		},
		{
			title: "-y",
			ray:   NewRay(calc.NewPoint(0.5, -5, 0), calc.NewVector(0, 1, 0)),
			t1:    4,
			t2:    6,
		},
		{
			title: "+z",
			ray:   NewRay(calc.NewPoint(0.5, 0, 5), calc.NewVector(0, 0, -1)),
			t1:    4,
			t2:    6,
		},
		{
			title: "-z",
			ray:   NewRay(calc.NewPoint(0.5, 0, -5), calc.NewVector(0, 0, 1)),
			t1:    4,
			t2:    6,
		},
		{
			title: "inside",
			ray:   NewRay(calc.NewPoint(0, 0.5, 0), calc.NewVector(0, 0, 1)),
			t1:    -1,
			t2:    1,
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, 2, xs.Count)
			require.Equal(t, target.t1, xs.Intersections[0].Time)
			require.Equal(t, target.t2, xs.Intersections[1].Time)

		})
	}

}

func TestRayMissing(t *testing.T) {
	c := NewCube()
	for _, target := range []struct {
		title string
		ray   Ray
	}{
		{
			title: "1",
			ray:   NewRay(calc.NewPoint(-2, 0, 0), calc.NewVector(0.2673, 0.5345, 0.8018)),
		},
		{
			title: "2",
			ray:   NewRay(calc.NewPoint(0, -2, 0), calc.NewVector(0.8018, 0.2673, 0.5345)),
		},
		{
			title: "3",
			ray:   NewRay(calc.NewPoint(0, 0, -2), calc.NewVector(0.5345, 0.8018, 0.2673)),
		},
		{
			title: "4",
			ray:   NewRay(calc.NewPoint(2, 0, 2), calc.NewVector(0, 0, -1)),
		},
		{
			title: "5",
			ray:   NewRay(calc.NewPoint(0, 2, 2), calc.NewVector(0, -1, 0)),
		},
		{
			title: "6",
			ray:   NewRay(calc.NewPoint(2, 2, 0), calc.NewVector(-1, 0, 0)),
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			xs, err := c.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, 0, xs.Count)
		})
	}
}

func TestNormalOnCube(t *testing.T) {
	c := NewCube()

	for _, target := range []struct {
		title string
		point calc.Tuple4
		ans   calc.Tuple4
	}{
		{
			"1",
			calc.NewPoint(1, 0.5, -0.8),
			calc.NewVector(1, 0, 0),
		},
		{
			"2",
			calc.NewPoint(-1, -0.2, 0.9),
			calc.NewVector(-1, 0, 0),
		},
		{
			"3",
			calc.NewPoint(-0.4, 1, -0.1),
			calc.NewVector(0, 1, 0),
		},
		{
			"4",
			calc.NewPoint(0.3, -1, -0.7),
			calc.NewVector(0, -1, 0),
		},
		{
			"5",
			calc.NewPoint(-0.6, 0.3, 1),
			calc.NewVector(0, 0, 1),
		},
		{
			"6",
			calc.NewPoint(0.4, 0.4, -1),
			calc.NewVector(0, 0, -1),
		},
		{
			"7",
			calc.NewPoint(1, 1, 1),
			calc.NewVector(1, 0, 0),
		},
		{
			"8",
			calc.NewPoint(-1, -1, -1),
			calc.NewVector(-1, 0, 0),
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			normal := c.calcLocalNormal(target.point)
			require.True(t, calc.TupleCompare(target.ans, normal))
		})
	}
}
