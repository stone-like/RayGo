package chap11

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func TestChap11(t *testing.T) {

	glassfloor := scene.NewPlane()
	//ポインタじゃないとsetが反映されないっぽい
	glassfloorMaterial := scene.DefaultMaterial()
	glassfloorMaterial.SetPattern(scene.NewCheckersPattern(scene.White, scene.Black))
	glassfloorMaterial.Transparency = 1
	glassfloorMaterial.Reflective = 1
	glassfloorMaterial.Shininess = 300
	glassfloorMaterial.Specular = 1
	glassfloorMaterial.Diffuse = 0.1

	glassfloor.SetMaterial(glassfloorMaterial)

	// wall := scene.NewPlane()
	// wall.SetTransform(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewTranslation(0, 0, 3)))
	// wallMaterial := scene.DefaultMaterial()
	// wallMaterial.Color = scene.NewColor(0, 0, 0)
	// wall.SetMaterial(wallMaterial)

	middle := scene.NewSphere(1)
	middle.SetTransform(calc.NewTranslation(-0.5, 1, 0.5))
	middleMaterial := scene.DefaultMaterial()
	middleMaterial.Color = scene.NewColor(1, 0, 0)
	middleMaterial.Diffuse = 0.1
	middleMaterial.Ambient = 0.1
	middleMaterial.Transparency = 1.0
	middleMaterial.RefractiveIndex = 1.5

	middle.SetMaterial(middleMaterial)

	right := scene.NewSphere(1)
	right.SetTransform(calc.NewTranslation(1.5, 0.5, -0.5).MulByMat4x4(calc.NewScale(0.5, 0.5, 0.5)))
	rightMaterial := scene.DefaultMaterial()
	rightMaterial.Color = scene.NewColor(0.5, 1, 0.1)
	rightMaterial.Diffuse = 0.1
	rightMaterial.Ambient = 0.1
	rightMaterial.Transparency = 1.0
	rightMaterial.RefractiveIndex = 1.5

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

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), glassfloor, middle, left, right)

	camera := scene.NewCamera(300, 450, math.Pi/3)
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
