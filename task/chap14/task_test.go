package chap14

import (
	"fmt"
	"math"
	"os"
	"rayGo/calc"
	"rayGo/scene"
	"testing"
)

func hexagon_corner() scene.Shape {
	corner := scene.NewSphere(1)
	corner.SetTransform(calc.NewTranslation(0, 0, -1).MulByMat4x4(calc.NewScale(0.25, 0.25, 0.25)))

	return corner
}

func hexagon_edge() scene.Shape {
	edge := scene.NewCyliner(scene.CynMin(0), scene.CynMax(1))
	edge.SetTransform(calc.NewTranslation(0, 0, -1).MulByMat4x4(calc.NewRotateY(-math.Pi / 6).MulByMat4x4(calc.NewRotateZ(-math.Pi / 2).MulByMat4x4(calc.NewScale(0.25, 1, 0.25)))))

	return edge
}

func hexagon_side() scene.Shape {
	side := scene.NewGroup()

	side.AddChildren(hexagon_corner(), hexagon_edge())
	return side
}

func hexagon() scene.Shape {
	hex := scene.NewGroup()

	for i := 0; i <= 5; i++ {
		side := hexagon_side()
		side.SetTransform(calc.NewRotateY(float64(i) * math.Pi / 3))
		hex.AddChildren(side)
	}

	return hex
}
func TestChap14(t *testing.T) {

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

	hex := hexagon()
	// hex.SetTransform(calc.NewRotateZ(math.Pi / 6).MulByMat4x4(calc.NewRotateX(math.Pi / 6).MulByMat4x4(calc.NewTranslation(0, 1, 0))))
	hex.SetTransform(calc.MulMatMulti(calc.NewRotateZ(math.Pi/6), calc.NewRotateX(math.Pi/6), calc.NewTranslation(0, 1, 0)))

	world := scene.NewWorld(scene.NewLight(calc.NewPoint(-10, 10, -10), scene.NewColor(1, 1, 1)), glassfloor, hex)

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
