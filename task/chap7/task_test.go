package chap7

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func TestChap7(t *testing.T) {

	floor := scene.NewSphere(1)
	floor.SetTransform(calc.NewScale(10, 0.01, 10))
	floorMaterial := scene.DefaultMaterial()
	floorMaterial.Color = scene.NewColor(1, 0.9, 0.9)
	floorMaterial.Specular = 0
	floor.SetMaterial(floorMaterial)

	left_wall := scene.NewSphere(1)
	left_wall_Trans := calc.NewTranslation(0, 0, 5).MulByMat4x4(calc.NewRotateY(-math.Pi / 4).MulByMat4x4(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewScale(10, 0.01, 10))))
	left_wall.SetTransform(left_wall_Trans)
	left_wall.SetMaterial(floorMaterial)

	right_wall := scene.NewSphere(1)
	right_wall_Trans := calc.NewTranslation(0, 0, 5).MulByMat4x4(calc.NewRotateY(math.Pi / 4).MulByMat4x4(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewScale(10, 0.01, 10))))
	right_wall.SetTransform(right_wall_Trans)
	right_wall.SetMaterial(floorMaterial)

	middle := scene.NewSphere(1)
	middle.SetTransform(calc.NewTranslation(-0.5, 1, 0.5))
	middleMaterial := scene.DefaultMaterial()
	middleMaterial.Color = scene.NewColor(0.1, 1, 0.5)
	middleMaterial.Diffuse = 0.7
	middleMaterial.Specular = 0.3
	middle.SetMaterial(middleMaterial)

	right := scene.NewSphere(1)
	right.SetTransform(calc.NewTranslation(1.5, 0.5, -0.5).MulByMat4x4(calc.NewScale(0.5, 0.5, 0.5)))
	rightMaterial := scene.DefaultMaterial()
	rightMaterial.Color = scene.NewColor(0.5, 1, 0.1)
	rightMaterial.Diffuse = 0.7
	rightMaterial.Specular = 0.3
	right.SetMaterial(rightMaterial)

	left := scene.NewSphere(1)
	left.SetTransform(calc.NewTranslation(-1.5, 0.33, -0.75).MulByMat4x4(calc.NewScale(0.33, 0.33, 0.33)))
	leftMaterial := scene.DefaultMaterial()
	leftMaterial.Color = scene.NewColor(1, 0.8, 0.1)
	leftMaterial.Diffuse = 0.7
	leftMaterial.Specular = 0.3
	left.SetMaterial(leftMaterial)

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), floor, left_wall, right_wall, middle, left, right)

	camera := scene.NewCamera(600, 900, math.Pi/3)
	camera.Transform = scene.ViewTransform(
		calc.NewPoint(0, 1.5, -5),
		calc.NewPoint(-1.25, -0.7, 0),
		calc.NewVector(0, 1, 0),
	)

	canvas, err := world.Render(camera)
	if err != nil {
		fmt.Println(err)
		return
	}

	fp, err := os.Create("test.ppm")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fp.Close()

	fp.WriteString(canvas.ToPPM())
}
