package scene

import (
	"math"
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewTransFormation(t *testing.T) {
	for _, target := range []struct {
		title string
		from  calc.Tuple4
		to    calc.Tuple4
		up    calc.Tuple4
		ans   calc.Mat4x4
	}{
		{
			title: "transform matrix for default orientation",
			from:  calc.NewPoint(0, 0, 0),
			to:    calc.NewPoint(0, 0, -1),
			up:    calc.NewVector(0, 1, 0),
			ans:   calc.Ident4x4,
		},
		{
			title: "transform matrix looking int positive z axis",
			from:  calc.NewPoint(0, 0, 0),
			to:    calc.NewPoint(0, 0, 1),
			up:    calc.NewVector(0, 1, 0),
			ans:   calc.NewScale(-1, 1, -1),
		},
		{
			title: "transform matrix moves the world",
			from:  calc.NewPoint(0, 0, 8),
			to:    calc.NewPoint(0, 0, 0),
			up:    calc.NewVector(0, 1, 0),
			ans:   calc.NewTranslation(0, 0, -8),
		},
		{
			title: "arbitary view transformation",
			from:  calc.NewPoint(1, 3, 2),
			to:    calc.NewPoint(4, -2, 8),
			up:    calc.NewVector(1, 1, 0),
			ans: calc.Mat4x4{
				{-0.50709, 0.50709, 0.67612, -2.36643},
				{0.76772, 0.60609, 0.12122, -2.82843},
				{-0.35857, 0.59761, -0.71714, 0.00000},
				{0.00000, 0.00000, 0.00000, 1.00000},
			},
		},
	} {
		t.Run(target.title, func(t *testing.T) {
			require.True(t, calc.Mat4x4Compare(target.ans, ViewTransform(target.from, target.to, target.up)))
		})
	}
}

func Test_PixelSize(t *testing.T) {
	camera := NewCamera(200, 125, math.Pi/2)

	require.Equal(t, 0.01, camera.PixelSize)

	camera = NewCamera(125, 200, math.Pi/2)

	require.Equal(t, 0.01, camera.PixelSize)
}

func Test_Ray_For_Pixel(t *testing.T) {
	camera := NewCamera(201, 101, math.Pi/2)

	for _, target := range []struct {
		title     string
		px        float64
		py        float64
		transform calc.Mat4x4
		point     calc.Tuple4
		ans       calc.Tuple4
	}{
		{
			title:     "constructing Ray through center of canvas",
			px:        100,
			py:        50,
			transform: calc.Ident4x4,
			point:     calc.NewPoint(0, 0, 0),
			ans:       calc.NewVector(0, 0, -1),
		},
		{
			title:     "constructing Ray through corner of canvas",
			px:        0,
			py:        0,
			transform: calc.Ident4x4,
			point:     calc.NewPoint(0, 0, 0),
			ans:       calc.NewVector(0.66519, 0.33259, -0.66851),
		},
		{
			title:     "constructing Ray when camera is transformed",
			px:        100,
			py:        50,
			transform: calc.NewRotateY(math.Pi / 4).MulByMat4x4(calc.NewTranslation(0, -2, 5)),
			point:     calc.NewPoint(0, 2, -5),
			ans:       calc.NewVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2),
		},
	} {

		t.Run(target.title, func(t *testing.T) {

			camera.Transform = target.transform

			ray, err := camera.RayForPixel(target.px, target.py)
			if err != nil {
				t.Error(err)
				return
			}

			require.True(t, calc.TupleCompare(target.point, ray.Origin))
			require.True(t, calc.TupleCompare(target.ans, ray.Direction))

		})

	}

}

func Test_Render_World_With_Camera(t *testing.T) {
	w := DefaultWorld()
	camera := NewCamera(11, 11, math.Pi/2)
	from := calc.NewPoint(0, 0, -5)
	to := calc.NewPoint(0, 0, 0)
	up := calc.NewVector(0, 1, 0)

	camera.Transform = ViewTransform(from, to, up)

	canvas, err := w.Render(camera)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.38066, 0.47583, 0.2855), canvas.Pixels[5][5]))

}
