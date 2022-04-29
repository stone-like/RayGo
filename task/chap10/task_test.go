package chap10

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func TestChap10(t *testing.T) {

	floor := scene.NewPlane()
	//ポインタじゃないとsetが反映されないっぽい
	floorMaterial := scene.DefaultMaterial()
	floorMaterial.SetPattern(scene.NewCheckersPattern(scene.White, scene.Black))

	floor.SetMaterial(floorMaterial)

	// wall := scene.NewPlane()
	// wall.SetTransform(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewTranslation(0, 0, 3)))
	// wallMaterial := scene.DefaultMaterial()
	// wallMaterial.Color = scene.NewColor(0, 0, 0)
	// wall.SetMaterial(wallMaterial)

	middle := scene.NewSphere(1)
	middle.SetTransform(calc.NewTranslation(-0.5, 1, 0.5))
	middleMaterial := scene.DefaultMaterial()
	// middleMaterial.Color = scene.NewColor(0.1, 1, 0.5)
	middleMaterial.Diffuse = 0.7
	middleMaterial.Specular = 0.3

	middlePattern := scene.NewRingPattern(scene.Orange, scene.Red)
	middlePattern.SetTransform(calc.NewRotateZ(math.Pi / 2).MulByMat4x4(calc.NewScale(0.1, 0.1, 0.5)))

	middleMaterial.SetPattern(middlePattern)
	middle.SetMaterial(middleMaterial)

	right := scene.NewSphere(1)
	right.SetTransform(calc.NewTranslation(1.5, 0.5, -0.5).MulByMat4x4(calc.NewScale(0.5, 0.5, 0.5)))
	rightMaterial := scene.DefaultMaterial()
	// rightMaterial.Color = scene.NewColor(0.5, 1, 0.1)
	rightMaterial.Diffuse = 0.7
	rightMaterial.Specular = 0.3
	rightPattern := scene.NewGradientPattern(scene.Blue, scene.MizuIro)
	rightPattern.SetTransform(calc.NewScale(0.5, 0.5, 0.5))

	rightMaterial.SetPattern(rightPattern)
	right.SetMaterial(rightMaterial)

	//チェッカーの貼り付けの不完全さを回避するためにはUVマッピングをしなくてはいけない
	left := scene.NewSphere(1)
	left.SetTransform(calc.NewTranslation(-1.5, 0.33, -0.75).MulByMat4x4(calc.NewScale(0.33, 0.33, 0.33)))
	leftMaterial := scene.DefaultMaterial()
	// leftMaterial.Color = scene.NewColor(1, 0.8, 0.1)
	leftMaterial.Diffuse = 0.7
	leftMaterial.Specular = 0.3

	leftPattern := scene.NewCheckersPattern(scene.Green, scene.FukaMidori)
	leftPattern.SetTransform(calc.NewScale(0.5, 0.5, 0.5))
	leftMaterial.SetPattern(leftPattern)
	left.SetMaterial(leftMaterial)

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), floor, middle, left, right)

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
