package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaneNormalAt(t *testing.T) {
	p := NewPlane()

	require.True(t, calc.TupleCompare(calc.NewVector(0, 1, 0), p.calcLocalNormal(calc.NewPoint(0, 0, 0), Intersection{})))
	require.True(t, calc.TupleCompare(calc.NewVector(0, 1, 0), p.calcLocalNormal(calc.NewPoint(10, 0, -10), Intersection{})))
	require.True(t, calc.TupleCompare(calc.NewVector(0, 1, 0), p.calcLocalNormal(calc.NewPoint(-5, 0, 150), Intersection{})))

}

func TestPlaneIntersect(t *testing.T) {
	p := NewPlane()

	for _, target := range []struct {
		title string
		ray   Ray
		xs    Intersections
	}{
		{
			"intersect with ray parallel to the plane",
			NewRay(calc.NewPoint(0, 10, 0), calc.NewVector(0, 0, 1)),
			Intersections{},
		},
		{
			"intersect with ray coplanar to the plane",
			NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1)),
			Intersections{},
		},
		{
			"intersect with ray above plane",
			NewRay(calc.NewPoint(0, 1, 0), calc.NewVector(0, -1, 0)),
			Intersections{
				Count: 1,
				Intersections: []*Intersection{
					{
						Time:   1,
						Object: p,
					},
				},
			},
		},
		{
			"intersect with ray below plane",
			NewRay(calc.NewPoint(0, -1, 0), calc.NewVector(0, 1, 0)),
			Intersections{
				Count: 1,
				Intersections: []*Intersection{
					{
						Time:   1,
						Object: p,
					},
				},
			},
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			xs, err := p.calcLocalIntersect(target.ray)
			require.Nil(t, err)
			require.Equal(t, target.xs, xs)
		})
	}
}
