package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ring_Pattern(t *testing.T) {
	p := NewRingPattern(White, Black)
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(1, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(0, 0, 1))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(0.708, 0, 0.708))))
}
