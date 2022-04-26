package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ray_Position(t *testing.T) {
	r := NewRay(calc.NewPoint(2, 3, 4), calc.NewVector(1, 0, 0))

	for _, target := range []struct {
		t   float64
		ans calc.Tuple4
	}{
		{
			0,
			calc.NewPoint(2, 3, 4),
		},
		{
			1,
			calc.NewPoint(3, 3, 4),
		},
		{
			-1,
			calc.NewPoint(1, 3, 4),
		},
		{
			2.5,
			calc.NewPoint(4.5, 3, 4),
		},
	} {
		require.Equal(t, target.ans, r.Position(target.t))
	}
}

func Test_Ray_Intersects_Sphere_At_TwoPoints(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, 4.0, xs.Intersections[0].Time)
	require.Equal(t, 6.0, xs.Intersections[1].Time)
	require.Equal(t, s, xs.Intersections[0].Object)
	require.Equal(t, s, xs.Intersections[1].Object)
}

func Test_Ray_Intersects_Sphere_At_Tangent(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 1, -5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, 5.0, xs.Intersections[0].Time)
	require.Equal(t, 5.0, xs.Intersections[1].Time)
}

func Test_Ray_Misses_Sphere(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 2, -5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	xs, err := s.Intersect(r)
	require.Nil(t, err)
	require.Equal(t, 0, xs.Count)
	require.Equal(t, 0, len(xs.Intersections))
}

func Test_Ray_Originates_Inside_Sphere(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, -1.0, xs.Intersections[0].Time)
	require.Equal(t, 1.0, xs.Intersections[1].Time)
}

func Test_Ray_Originates_Behind_Sphere(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, 5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, -6.0, xs.Intersections[0].Time)
	require.Equal(t, -4.0, xs.Intersections[1].Time)
}

func Test_Aggregating_Intersections(t *testing.T) {
	s := NewSphere(1)
	i1 := s.GetIntersection(1)
	i2 := s.GetIntersection(2)

	xs := AggregateIntersection(i1, i2)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, 1.0, xs.Intersections[0].Time)
	require.Equal(t, 2.0, xs.Intersections[1].Time)

}

func Test_Hit(t *testing.T) {
	s := NewSphere(1)

	for _, target := range []struct {
		sections Intersections
		ans      *Intersection
	}{
		{
			AggregateIntersection(s.GetIntersection(2), s.GetIntersection(1)),
			s.GetIntersection(1),
		},
		{
			AggregateIntersection(s.GetIntersection(-1), s.GetIntersection(1)),
			s.GetIntersection(1),
		},
		{
			AggregateIntersection(s.GetIntersection(-2), s.GetIntersection(-1)),
			nil,
		},
		{
			AggregateIntersection(s.GetIntersection(5), s.GetIntersection(7), s.GetIntersection(-3), s.GetIntersection(2)),
			s.GetIntersection(2),
		},
	} {
		require.Equal(t, target.ans, GenerateHit(target.sections))
	}

}

//実際はShapeのTransformのInverseをray.Transformに渡すことになる、
//ここではTransformがあっているかだけ見る
func Test_Translating_Ray(t *testing.T) {
	r := NewRay(calc.NewPoint(1, 2, 3), calc.NewVector(0, 1, 0))
	trans := calc.NewTranslation(3, 4, 5)
	r2 := r.Transform(trans)

	require.Equal(t, calc.NewPoint(4, 6, 8), r2.Origin)
	require.Equal(t, calc.NewVector(0, 1, 0), r2.Direction)

}

//実際のscailingではshapeが三倍になるとしたら、rayのVectorは1/3になる
func Test_Scailing_Ray(t *testing.T) {
	r := NewRay(calc.NewPoint(1, 2, 3), calc.NewVector(0, 1, 0))
	trans := calc.NewScale(2, 3, 4)
	r2 := r.Transform(trans)

	require.Equal(t, calc.NewPoint(2, 6, 12), r2.Origin)
	require.Equal(t, calc.NewVector(0, 3, 0), r2.Direction)

}

func Test_Ray_Intersects_Scaling_Sphere(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	s.SetTransform(calc.NewScale(2, 2, 2))

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, 3.0, xs.Intersections[0].Time)
	require.Equal(t, 7.0, xs.Intersections[1].Time)
}

func Test_Ray_Intersects_Translated_Sphere(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	s := NewSphere(1)

	s.SetTransform(calc.NewTranslation(5, 0, 0))

	xs, err := s.Intersect(r)
	require.Nil(t, err)

	require.Equal(t, 0, xs.Count)
	require.Equal(t, 0, len(xs.Intersections))
}
