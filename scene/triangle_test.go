package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTriangle(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	require.True(t, calc.TupleCompare(calc.NewVector(-1, -1, 0), tri.E1))
	require.True(t, calc.TupleCompare(calc.NewVector(1, -1, 0), tri.E2))
	require.True(t, calc.TupleCompare(calc.NewVector(0, 0, -1), tri.NormalVec))

}

func TestNormalOnTriangle(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	n1 := tri.calcLocalNormal(calc.NewPoint(0, 0.5, 0), Intersection{})
	n2 := tri.calcLocalNormal(calc.NewPoint(-0.5, 0.75, 0), Intersection{})
	n3 := tri.calcLocalNormal(calc.NewPoint(0.5, 0.25, 0), Intersection{})

	require.True(t, calc.TupleCompare(tri.NormalVec, n1))
	require.True(t, calc.TupleCompare(tri.NormalVec, n2))
	require.True(t, calc.TupleCompare(tri.NormalVec, n3))

}

func Test_Intersect_Ray_Parallel_To_The_Triangle(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	ray := NewRay(calc.NewPoint(0, -1, -2), calc.NewVector(0, 1, 0))
	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
}

func Test_Ray_Misses_The_P1ToP3Edge(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	ray := NewRay(calc.NewPoint(1, 1, -2), calc.NewVector(0, 0, 1))
	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
}

func Test_Ray_Misses_The_P1ToP2Edge(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	ray := NewRay(calc.NewPoint(-1, 1, -2), calc.NewVector(0, 0, 1))
	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
}

func Test_Ray_Misses_The_P2ToP3Edge(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	ray := NewRay(calc.NewPoint(0, -1, -2), calc.NewVector(0, 0, 1))
	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
}

func Test_Ray_Strikes_a_Triangle(t *testing.T) {
	tri := NewTriangle(calc.NewPoint(0, 1, 0), calc.NewPoint(-1, 0, 0), calc.NewPoint(1, 0, 0))

	ray := NewRay(calc.NewPoint(0, 0.5, -2), calc.NewVector(0, 0, 1))
	xs, err := tri.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 1, xs.Count)
	require.Equal(t, float64(2), xs.Intersections[0].Time)

}
