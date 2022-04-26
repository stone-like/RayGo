package scene

import (
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
