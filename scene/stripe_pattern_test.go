package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Create_Stripe(t *testing.T) {
	p := NewStripePattern(White, Black)

	require.True(t, colorCompare(White, p.Color1))
	require.True(t, colorCompare(Black, p.Color2))

}

func Test_StripePattern_is_Constant_in_Y(t *testing.T) {
	p := NewStripePattern(White, Black)

	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 1, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 2, 0))))

}

func Test_StripePattern_is_Constant_in_Z(t *testing.T) {
	p := NewStripePattern(White, Black)

	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 1))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 2))))

}

func Test_StripePattern_is_Alternates_in_X(t *testing.T) {
	p := NewStripePattern(White, Black)
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(0.9, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(1, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(-0.1, 0, 0))))
	require.True(t, colorCompare(Black, p.PatternAt(calc.NewPoint(-1, 0, 0))))
	require.True(t, colorCompare(White, p.PatternAt(calc.NewPoint(-1.1, 0, 0))))

}

func Test_Stripe_With_Object_Transformation(t *testing.T) {
	s := NewSphere(1)

	s.SetTransform(calc.NewScale(2, 2, 2))
	//sphereが2倍にscaleしているので、Stripeも幅が二倍にならなければいけない

	pattern := NewStripePattern(White, Black)
	//objectのスケールで0~1までがWhiteだったのが、0~2までWhiteになる
	c, err := pattern.PatternAtShape(calc.NewPoint(1.5, 0, 0), s)
	require.Nil(t, err)

	require.True(t, colorCompare(White, c))

}

func Test_Stripe_With_Pattern_Transformation(t *testing.T) {
	s := NewSphere(1)

	pattern := NewStripePattern(White, Black)
	pattern.SetTransform(calc.NewScale(2, 2, 2))
	//patternのスケールで0~1までがWhiteだったのが、0~2までWhiteになる
	c, err := pattern.PatternAtShape(calc.NewPoint(1.5, 0, 0), s)
	require.Nil(t, err)

	require.True(t, colorCompare(White, c))

}

func Test_Stripe_With_Both_Object_And_Pattern_Transformation(t *testing.T) {
	s := NewSphere(1)
	s.SetTransform(calc.NewScale(2, 2, 2))

	pattern := NewStripePattern(White, Black)
	pattern.SetTransform(calc.NewTranslation(0.5, 0, 0))
	//patternのスケールで0~1までがWhiteだったのが、0~2までWhiteになる
	//Transitionで1~2.5までwhite
	c, err := pattern.PatternAtShape(calc.NewPoint(2.5, 0, 0), s)
	require.Nil(t, err)

	require.True(t, colorCompare(White, c))

	c, err = pattern.PatternAtShape(calc.NewPoint(1, 0, 0), s)
	require.Nil(t, err)

	require.True(t, colorCompare(White, c))

	c, err = pattern.PatternAtShape(calc.NewPoint(0.9, 0, 0), s)
	require.Nil(t, err)

	require.True(t, colorCompare(Black, c))

}
