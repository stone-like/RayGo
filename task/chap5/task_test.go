package task

import (
	"fmt"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func TestChap5(t *testing.T) {
	ray_origin := calc.NewPoint(0, 0, -5)
	wall_z := 10.0
	wall_size := 7.0

	canvas_pixels := 100
	pixel_size := wall_size / float64(canvas_pixels)
	half := wall_size / 2

	canvas := scene.NewCanvas(canvas_pixels, canvas_pixels)
	color := scene.NewColor(1, 0, 0) //red

	sphere := scene.NewSphere(1)

	for y := 0; y < canvas_pixels; y++ {

		world_y := half - pixel_size*float64(y)

		for x := 0; x < canvas_pixels; x++ {
			world_x := -half + pixel_size*float64(x)

			position := calc.NewPoint(world_x, world_y, wall_z)

			ray := scene.NewRay(ray_origin, calc.SubTuple(position, ray_origin).Normalize())

			xs, err := sphere.Intersect(ray)

			if err != nil {
				panic("error occur!")
			}

			hit := scene.GenerateHit(xs)

			if hit != nil {
				canvas.WritePixel(x, y, color)
			}

		}
	}

	fp, err := os.Create("test.ppm")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fp.Close()

	fp.WriteString(canvas.ToPPM())
}
