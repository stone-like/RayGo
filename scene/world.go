package scene

import (
	"rayGo/calc"
	"sort"
)

type World struct {
	Light   Light
	Objects []Shape
}

func NewWorld(light Light, objects ...Shape) World {
	return World{
		Light:   light,
		Objects: objects,
	}
}

func (w World) Intersect(r Ray) (Intersections, error) {

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

func (w World) ShadeHit(comps PreComps) Color {
	return w.Light.Lighting(
		comps.Object.GetMaterial(),
		comps.RayPoint,
		comps.EyeVec,
		comps.NormalVec,
	)
}

func (w World) ColorAt(ray Ray) (Color, error) {
	xs, err := w.Intersect(ray)
	if err != nil {
		return Color{}, err
	}

	hit := GenerateHit(xs)

	if hit == nil {
		return Black, nil
	}

	comps, err := PrepareComputations(*hit, ray)

	if err != nil {
		return Color{}, err
	}

	return w.ShadeHit(comps), nil

}

func DefaultWorld() World {
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
