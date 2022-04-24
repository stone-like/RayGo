package scene

import (
	"testing"

	"rayGo/util"

	"github.com/stretchr/testify/require"
)

func colorCompare(c1, c2 Color) bool {
	return util.FloatEqual(c1.Red, c2.Red) && util.FloatEqual(c1.Green, c2.Green) && util.FloatEqual(c1.Blue, c2.Blue)
}

func TestColorAdd(t *testing.T) {
	c1 := NewColor(0.9, 0.6, 0.75)
	c2 := NewColor(0.7, 0.1, 0.25)
	require.True(t, colorCompare(NewColor(1.6, 0.7, 1.0), c1.Add(c2)))
}

func TestColorSub(t *testing.T) {
	c1 := NewColor(0.9, 0.6, 0.75)
	c2 := NewColor(0.7, 0.1, 0.25)
	require.True(t, colorCompare(NewColor(0.2, 0.5, 0.5), c1.Sub(c2)))
}

func TestColorMul(t *testing.T) {
	c1 := NewColor(1, 0.2, 0.4)
	c2 := NewColor(0.9, 1, 0.1)
	//goのfloatの特性で0.4*0.1=0.04000000000001みたいになるので比較はutilのfloatEqualを使う
	require.True(t, colorCompare(NewColor(0.9, 0.2, 0.04), c1.Mul(c2)))
}
