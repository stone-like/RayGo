package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Gradient_Pattern(t *testing.T) {
	p := NewGradientPattern(White, Black)

	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(NewColor(0.75, 0.75, 0.75), p.PatternAt(calc.NewPoint(0.25, 0, 0))))
	require.True(t, colorCompare(NewColor(0.5, 0.5, 0.5), p.PatternAt(calc.NewPoint(0.5, 0, 0))))
	require.True(t, colorCompare(NewColor(0.25, 0.25, 0.25), p.PatternAt(calc.NewPoint(0.75, 0, 0))))

}
