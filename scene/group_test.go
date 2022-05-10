package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddChildToGroup(t *testing.T) {
	s := NewSphere(1)

	g := NewGroup()

	g.AddChildren(s)

	tt := s.GetParent()

	require.Equal(t, 1, len(g.Children))
	require.Equal(t, s, g.Children[0])
	require.Equal(t, g, tt)

}

func Test_Intersect_Ray_With_EmptyGroup(t *testing.T) {
	g := NewGroup()
	ray := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1))

	xs, err := g.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
}

func Test_Intersect_Ray_With_NonEmptyGroup(t *testing.T) {
	g := NewGroup()

	s1 := NewSphere(1)
	s2 := NewSphere(1)
	s2.SetTransform(calc.NewTranslation(0, 0, -3))
	s3 := NewSphere(1)
	s3.SetTransform(calc.NewTranslation(5, 0, 0))

	g.AddChildren(s1, s2, s3)

	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))

	xs, err := g.calcLocalIntersect(ray)
	require.Nil(t, err)
	require.Equal(t, 4, xs.Count)
	require.Equal(t, s2, xs.Intersections[0].Object)
	require.Equal(t, s2, xs.Intersections[1].Object)
	require.Equal(t, s1, xs.Intersections[2].Object)
	require.Equal(t, s1, xs.Intersections[3].Object)

}

func Test_Intersect_Transformed_Group(t *testing.T) {
	g := NewGroup()
	g.SetTransform(calc.NewScale(2, 2, 2))

	s := NewSphere(1)
	s.SetTransform(calc.NewTranslation(5, 0, 0))

	g.AddChildren(s)

	ray := NewRay(calc.NewPoint(10, 0, -10), calc.NewVector(0, 0, 1))

	xs, err := g.Intersect(ray)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
}

func Test_Converting_Point_From_World_to_Object_Space(t *testing.T) {

	// g1 -> g2 -> sの関係
	g1 := NewGroup()
	g1.SetTransform(calc.NewRotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(calc.NewScale(2, 2, 2))

	g1.AddChildren(g2)

	s := NewSphere(1)
	s.SetTransform(calc.NewTranslation(5, 0, 0))

	g2.AddChildren(s)

	p, err := s.WorldToObject(calc.NewPoint(-2, 0, -10))
	require.Nil(t, err)

	require.True(t, calc.TupleCompare(calc.NewPoint(0, 0, -1), p))

}

func Test_Converting_Normal_From_Object_To_WorldSpace(t *testing.T) {
	// g1 -> g2 -> sの関係
	g1 := NewGroup()
	g1.SetTransform(calc.NewRotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(calc.NewScale(1, 2, 3))

	g1.AddChildren(g2)

	s := NewSphere(1)
	s.SetTransform(calc.NewTranslation(5, 0, 0))

	g2.AddChildren(s)

	p, err := s.NormalToWorld(calc.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
	require.Nil(t, err)

	util.SetEpsilon(0.001)
	defer util.SetEpsilon(util.DefaultEpsilon)
	require.True(t, calc.TupleCompare(calc.NewVector(0.2857, 0.4286, -0.8571), p))

}
