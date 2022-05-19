package chap16

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func createGemini() (scene.Shape, error) {
	s := scene.NewSphere(1)
	sMaterial := s.GetMaterial()
	sMaterial.Color = scene.Red
	s.SetMaterial(sMaterial)

	s2 := scene.NewSphere(1)
	s2Material := s2.GetMaterial()
	s2Material.Color = scene.Blue
	s2.SetMaterial(s2Material)
	s2.SetTransform(calc.NewTranslation(1.5, 0, 0))

	return scene.NewCSG(scene.CSGUnion, s, s2)

}

func TestChap16(t *testing.T) {

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

	gemini, err := createGemini()

	if err != nil {
		fmt.Println(err)
		return
	}

	gemini.SetTransform(calc.MulMatMulti(calc.NewTranslation(-0.5, 1, 0)))

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), glassfloor, gemini)

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
