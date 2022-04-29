package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Checkers_Should_Repeat_X(t *testing.T) {
	p := NewCheckersPattern(White, Black)
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0.99, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(1.01, 0, 0))))
}

func Test_Checkers_Should_Repeat_Y(t *testing.T) {
	p := NewCheckersPattern(White, Black)
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0.99, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(0, 1.01, 0))))
}

func Test_Checkers_Should_Repeat_Z(t *testing.T) {
	p := NewCheckersPattern(White, Black)
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0.99))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(0, 0, 1.01))))
}
