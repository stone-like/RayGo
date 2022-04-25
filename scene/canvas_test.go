package scene

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanvas(t *testing.T) {
	c := NewCanvas(10, 20)
	require.Equal(t, 10, c.Width)
	require.Equal(t, 20, c.Height)

	target := NewColor(0, 0, 0)

	for i := 0; i < c.Height; i++ {
		for j := 0; j < c.Width; j++ {
			require.True(t, colorCompare(c.Pixels[i][j], target))
		}
	}
}

func TestWrite(t *testing.T) {
	c := NewCanvas(10, 20)

	red := NewColor(1, 0, 0)
	c.WritePixel(3, 2, red)

	require.Equal(t, red, c.Pixels[2][3])

}

func TestPPM(t *testing.T) {
	c := NewCanvas(5, 3)
	c1 := NewColor(1.5, 0, 0)
	c2 := NewColor(0, 0.5, 0)
	c3 := NewColor(-0.5, 0, 1)

	c.WritePixel(0, 0, c1)
	c.WritePixel(2, 1, c2)
	c.WritePixel(4, 2, c3)

	target := "P3\n5 3\n255\n255 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n0 0 0 0 0 0 0 128 0 0 0 0 0 0 0\n0 0 0 0 0 0 0 0 0 0 0 0 0 0 255\n"

	require.Equal(t, target, c.ToPPM())
}

func TestPPMLongLine(t *testing.T) {
	c := NewCanvas(10, 2)
	nc := NewColor(1, 0.8, 0.6)

	for i := 0; i < c.Height; i++ {
		for j := 0; j < c.Width; j++ {
			c.WritePixel(j, i, nc)
		}
	}

	target := "P3\n10 2\n255\n255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204\n153 255 204 153 255 204 153 255 204 153 255 204 153\n255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204\n153 255 204 153 255 204 153 255 204 153 255 204 153\n"

	require.Equal(t, target, c.ToPPM())
}
