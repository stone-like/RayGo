package scene

import (
	"math"
	"rayGo/calc"
)

type Camera struct {
	HSize       float64
	VSize       float64
	FieldOfView float64
	Transform   calc.Mat4x4
	HalfHeight  float64
	HalfWidth   float64
	PixelSize   float64
}

func calcHalfWidthAndHeight(hSize, vSize, fieldOfView float64) (half_width, half_height float64) {
	half_view := math.Tan(fieldOfView / 2)
	aspect := hSize / vSize

	if aspect >= 1 {
		half_width = half_view
		half_height = half_view / aspect
	} else {
		half_width = half_view * aspect
		half_height = half_view
	}

	return

}

func NewCamera(hSize, vSize, fieldOfView float64) Camera {

	camera := Camera{
		HSize:       hSize,
		VSize:       vSize,
		FieldOfView: fieldOfView,
		Transform:   calc.Ident4x4,
	}

	half_width, half_height := calcHalfWidthAndHeight(hSize, vSize, fieldOfView)

	camera.HalfWidth, camera.HalfHeight = half_width, half_height

	camera.PixelSize = (camera.HalfWidth * 2) / camera.HSize

	return camera
}

func (c Camera) RayForPixel(px, py float64) (Ray, error) {
	x_offset := (px + 0.5) * c.PixelSize
	y_offset := (py + 0.5) * c.PixelSize

	world_x := c.HalfWidth - x_offset
	world_y := c.HalfHeight - y_offset

	cameraTransInv, err := c.Transform.Inverse()
	if err != nil {
		return Ray{}, err
	}

	pixel := cameraTransInv.MulByTuple(calc.NewPoint(world_x, world_y, -1))
	origin := cameraTransInv.MulByTuple(calc.NewPoint(0, 0, 0))
	direction := calc.SubTuple(pixel, origin).Normalize()

	return NewRay(origin, direction), nil

}

//ViewTransformの導出はPDFに書いてある
func ViewTransform(from, to, up calc.Tuple4) calc.Mat4x4 {
	forward := calc.SubTuple(to, from).Normalize()
	upn := up.Normalize()
	left := calc.CrossTuple(forward, upn)
	true_up := calc.CrossTuple(left, forward)

	neg_forward := calc.NegTuple(forward)
	neg_from := calc.NegTuple(from)

	orientation := calc.Mat4x4{
		{left[0], left[1], left[2], 0},
		{true_up[0], true_up[1], true_up[2], 0},
		{neg_forward[0], neg_forward[1], neg_forward[2], 0},
		{0, 0, 0, 1},
	}

	return orientation.MulByMat4x4(calc.NewTranslation(neg_from[0], neg_from[1], neg_from[2]))
}
