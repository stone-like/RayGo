package box

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func TestBoxRoom(t *testing.T) {

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

	wall1 := scene.NewPlane()
	wall1.SetTransform(calc.MulMatMulti(calc.NewTranslation(0, 0, 10), calc.NewRotateY(-math.Pi/4), calc.NewRotateX(math.Pi/2)))
	wall1Material := scene.DefaultMaterial()
	wall1Material.Color = scene.White
	wall1.SetMaterial(wall1Material)

	wall2 := scene.NewPlane()
	wall2.SetTransform(calc.MulMatMulti(calc.NewTranslation(0, 0, 10), calc.NewRotateY(math.Pi/4), calc.NewRotateX(math.Pi/2)))
	wall2Material := scene.DefaultMaterial()
	wall2Material.Color = scene.White
	wall2.SetMaterial(wall2Material)

	// wall3 := scene.NewPlane()
	// wall3.SetTransform(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewTranslation(0, 0, 3)))
	// wall3Material := scene.DefaultMaterial()
	// wall3Material.Color = scene.NewColor(0, 0, 0)
	// wall3.SetMaterial(wall3Material)

	// wall4 := scene.NewPlane()
	// wall4.SetTransform(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewTranslation(0, 0, 3)))
	// wall4Material := scene.DefaultMaterial()
	// wall4Material.Color = scene.NewColor(0, 0, 0)
	// wall4.SetMaterial(wall4Material)

	// roof := scene.NewPlane()
	// roof.SetTransform(calc.NewRotateX(math.Pi / 2).MulByMat4x4(calc.NewTranslation(0, 0, 3)))
	// roofMaterial := scene.DefaultMaterial()
	// roofMaterial.Color = scene.NewColor(0, 0, 0)
	// roof.SetMaterial(roofMaterial)

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

	left := scene.NewSphere(1)
	left.SetTransform(calc.NewTranslation(-1.5, 0.33, -0.75).MulByMat4x4(calc.NewScale(0.33, 0.33, 0.33)))
	leftMaterial := scene.DefaultMaterial()
	leftMaterial.Color = scene.Blue
	leftMaterial.Diffuse = 0.1
	leftMaterial.Ambient = 0.1
	leftMaterial.Transparency = 1.0
	leftMaterial.RefractiveIndex = 1.5

	left.SetMaterial(leftMaterial)

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), glassfloor, middle, left, right, wall1, wall2)

	camera := scene.NewCamera(960, 1200, math.Pi/3)

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
