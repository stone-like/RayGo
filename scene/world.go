package scene

import (
	"math"
	"rayGo/calc"
	"sort"
)

var DefaultRemaing = 5

type World struct {
	Light   Light
	Objects []Shape
}

func NewWorld(light Light, objects ...Shape) *World {
	return &World{
		Light:   light,
		Objects: objects,
	}
}

func (w *World) AddObjects(objects ...Shape) {
	w.Objects = append(w.Objects, objects...)
}

func (w *World) Intersect(r Ray) (Intersections, error) {

	count := 0
	var sections []*Intersection
	for _, shape := range w.Objects {
		shapeXS, err := shape.Intersect(r)
		if err != nil {
			return Intersections{}, err
		}

		sections = append(sections, shapeXS.Intersections...)
		count += shapeXS.Count

	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Time < sections[j].Time
	})

	return Intersections{
		Intersections: sections,
		Count:         count,
	}, nil
}

func (w *World) RefractedColor(comps PreComps, remainingReflection, remainingRefraction int) (Color, error) {
	if remainingRefraction <= 0 {
		return Black, nil
	}

	if comps.Object.GetMaterial().Transparency == 0 {
		return Black, nil
	}

	n_ratio := comps.N1 / comps.N2
	cos_i := calc.DotTuple(comps.EyeVec, comps.NormalVec)
	sin2_t := math.Pow(n_ratio, 2) * (1 - math.Pow(cos_i, 2))
	cos_t := math.Sqrt(1 - sin2_t)

	if sin2_t >= 1 {
		//Total Internal Reflection
		return Black, nil
	}

	direction := calc.SubTuple(
		calc.MulTupleByScalar(n_ratio*cos_i-cos_t, comps.NormalVec),
		calc.MulTupleByScalar(n_ratio, comps.EyeVec),
	)

	refract_ray := NewRay(comps.UnderPoint, direction)
	refract_color, err := w.ColorAt(refract_ray, remainingReflection, remainingRefraction-1)

	if err != nil {
		return Black, err
	}

	colorTuple := calc.MulTupleByScalar(
		comps.Object.GetMaterial().Transparency,
		refract_color.ToTuple4(),
	)

	return TupletoColor(colorTuple), nil
}

func (w *World) ReflectedColor(comps PreComps, remainingReflection, remainingRefraction int) (Color, error) {

	if remainingReflection <= 0 {
		return Black, nil
	}

	reflective := comps.Object.GetMaterial().Reflective
	if reflective == 0 {
		return Black, nil
	}

	reflect_ray := NewRay(comps.OverPoint, comps.ReflectVec)
	color, err := w.ColorAt(reflect_ray, remainingReflection-1, remainingRefraction)
	if err != nil {
		return Color{}, err
	}

	colorTuple := color.ToTuple4()
	return TupletoColor(calc.MulTupleByScalar(reflective, colorTuple)), nil
}

func isFresnelAppliable(material *Material) bool {
	if material.Reflective > 0 && material.Transparency > 0 {
		return true
	}

	return false

}

func (w *World) ApplyFresnel(comps PreComps, surface, reflected, refracted Color) Color {
	if !isFresnelAppliable(comps.Object.GetMaterial()) {
		return surface.Add(reflected).Add(refracted)
	}

	reflectance := comps.ComputeSchlick()

	appliedReflected := TupletoColor(calc.MulTupleByScalar(reflectance, reflected.ToTuple4()))
	appliedRefracted := TupletoColor(calc.MulTupleByScalar(1-reflectance, refracted.ToTuple4()))

	return surface.Add(appliedReflected).Add(appliedRefracted)

}

//rayとobjectの交点とずらしたOverPointを使わないと自分自身が自分と重なっている点として判定されてしまう
func (w *World) ShadeHit(comps PreComps, remainingReflection, remainingRefraction int) (Color, error) {

	in_shadow, err := w.IsShadowed(comps.OverPoint)
	if err != nil {
		return Color{}, err
	}

	sufaceColor, err := w.Light.Lighting(
		comps.Object.GetMaterial(),
		comps.RayPoint,
		comps.EyeVec,
		comps.NormalVec,
		in_shadow,
		comps.Object,
	)
	if err != nil {
		return Color{}, err
	}

	reflected, err := w.ReflectedColor(comps, remainingReflection, remainingRefraction)
	if err != nil {
		return Color{}, err
	}

	refracted, err := w.RefractedColor(comps, remainingReflection, remainingRefraction)
	if err != nil {
		return Color{}, err
	}

	return w.ApplyFresnel(comps, sufaceColor, reflected, refracted), nil

}

//光源とpointを結んでRayをつくってRayとWorldのIntersectionを求める
//hitがあり、tがdistanceより小さければpointはShadow
//それ以外はShadowでない
func (w *World) IsShadowed(point calc.Tuple4) (bool, error) {

	v := calc.SubTuple(w.Light.Position, point)
	distance := v.Magnitude()
	direction := v.Normalize()

	ray := NewRay(point, direction)
	xs, err := w.Intersect(ray)
	if err != nil {
		return false, err
	}

	hit := GenerateHit(xs)

	if hit != nil && hit.Time < distance {
		return true, nil
	}

	return false, nil
}

func (w *World) ColorAt(ray Ray, remainingReflection, remainingRefraction int) (Color, error) {
	xs, err := w.Intersect(ray)
	if err != nil {
		return Color{}, err
	}

	hit := GenerateHit(xs)

	if hit == nil {
		return Black, nil
	}

	comps, err := PrepareComputations(*hit, ray, xs)

	if err != nil {
		return Color{}, err
	}

	return w.ShadeHit(comps, remainingReflection, remainingRefraction)

}

func (w *World) Render(camera Camera) (*Canvas, error) {
	canvas := NewCanvas(int(camera.VSize), int(camera.HSize))

	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			ray, err := camera.RayForPixel(float64(x), float64(y))
			if err != nil {
				return nil, err
			}

			color, err := w.ColorAt(ray, DefaultRemaing, DefaultRemaing)
			if err != nil {
				return nil, err
			}

			canvas.WritePixel(x, y, color)
		}
	}

	return canvas, nil
}

func DefaultWorld() *World {
	light := (NewLight(calc.NewPoint(-10, 10, -10), NewColor(1, 1, 1)))

	s1 := NewSphere(1)

	m1 := s1.GetMaterial()
	m1.Color = NewColor(0.8, 1.0, 0.6)
	m1.Diffuse = 0.7
	m1.Specular = 0.2
	s1.SetMaterial(m1)

	s2 := NewSphere(1)
	s2.SetTransform(calc.NewScale(0.5, 0.5, 0.5))

	return NewWorld(light, s1, s2)
}
