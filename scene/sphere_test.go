package scene

import (
	"math"
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Sphere_Transformation(t *testing.T) {
	s := NewSphere(1)
	require.Equal(t, calc.Mat4x4(calc.Ident4x4), s.GetTransform())

	trans := calc.NewTranslation(2, 3, 4)
	s.SetTransform(trans)
	require.Equal(t, trans, s.GetTransform())

}

func Test_Normal(t *testing.T) {
	s := NewSphere(1)

	for _, target := range []struct {
		point calc.Tuple4
		ans   calc.Tuple4
	}{
		{calc.NewPoint(1, 0, 0),
			calc.NewVector(1, 0, 0)},
		{calc.NewPoint(0, 1, 0),
			calc.NewVector(0, 1, 0)},
		{calc.NewPoint(0, 0, 1),
			calc.NewVector(0, 0, 1)},
		{calc.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
			calc.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3)},
	} {

		n, err := s.NormalAt(target.point)

		require.Nil(t, err)
		require.Equal(t, target.ans, n)
	}
}

func Test_Normal_On_Translated_Sphere(t *testing.T) {
	s := NewSphere(1)
	s.SetTransform(calc.NewTranslation(0, 1, 0))

	n, err := s.NormalAt(calc.NewPoint(0, 1.70711, -0.70711))
	require.Nil(t, err)

	require.True(t, calc.TupleCompare(calc.NewVector(0, 0.70711, -0.70711), n))
}

func Test_Normal_On_Transformed_Sphere(t *testing.T) {
	s := NewSphere(1)

	trans := calc.NewScale(1, 0.5, 1).MulByMat4x4(calc.NewRotateZ(math.Pi / 5))
	s.SetTransform(trans)

	n, err := s.NormalAt(calc.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
	require.Nil(t, err)

	require.True(t, calc.TupleCompare(calc.NewVector(0, 0.97014, -0.24254), n))
}

func Test_Sphere_Has_Default_Material(t *testing.T) {
	s := NewSphere(1)

	m := DefaultMaterial()

	require.Equal(t, m, s.GetMaterial())
}

func Test_Sphere_Can_Set_Material(t *testing.T) {
	s := NewSphere(1)

	m := DefaultMaterial()
	m.ambient = 1

	s.SetMaterial(m)

	require.Equal(t, m, s.GetMaterial())
}
