package main

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
)

func createPacman() (scene.Shape, error) {
	// s := scene.NewSphere(1)
	// sMaterial := s.GetMaterial()
	// sMaterial.Color = scene.Red
	// s.SetMaterial(sMaterial)

	cube := scene.NewCube()
	// cube.SetTransform(calc.MulMatMulti(calc.NewRotateY(math.Pi/4), calc.NewScale(0.3, 1, 0.3), calc.NewTranslation(-0.3, 0, -0.5)))
	cube.SetTransform(calc.MulMatMulti(calc.NewTranslation(-0.5, 0, -0.7), calc.NewRotateX(-math.Pi/6), calc.NewRotateY(-math.Pi/6), calc.NewScale(1, 1, 0.5)))

	cubeMaterial := cube.GetMaterial()
	cubeMaterial.Color = scene.NewColor(1, 1, 1)
	cubeMaterial.Transparency = 1
	cubeMaterial.Reflective = 1
	cubeMaterial.Shininess = 0
	cubeMaterial.Specular = 0
	cubeMaterial.Diffuse = 0

	cube.SetMaterial(cubeMaterial)

	return cube, nil

}

func main() {
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

	pacman, err := createPacman()

	pacman.SetTransform(calc.MulMatMulti(calc.NewTranslation(0, 1, 0)))

	if err != nil {
		fmt.Println(err)
		return
	}

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), glassfloor, pacman)

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
