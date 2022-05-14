package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

type PreComps struct {
	Time        float64
	Object      Shape
	RayPoint    calc.Tuple4
	OverPoint   calc.Tuple4
	UnderPoint  calc.Tuple4
	EyeVec      calc.Tuple4
	NormalVec   calc.Tuple4
	ReflectVec  calc.Tuple4
	IsRayInside bool
	N1          float64
	N2          float64
}

func isTotalInternalReflection(n1, n2, sin2_t float64) bool {

	if n1 <= n2 {
		return false
	}

	if sin2_t <= 1.0 {
		return false
	}

	return true

}

func (c PreComps) ComputeSchlick() float64 {
	cos := calc.DotTuple(c.EyeVec, c.NormalVec)
	n := c.N1 / c.N2
	sin2_t := math.Pow(n, 2) * (1.0 - math.Pow(cos, 2))

	if isTotalInternalReflection(c.N1, c.N2, sin2_t) {
		return 1.0
	}

	if c.N1 > c.N2 {
		cos = math.Sqrt(1.0 - sin2_t)
	}

	r0 := math.Pow(((c.N1 - c.N2) / (c.N1 + c.N2)), 2)

	return r0 + (1-r0)*math.Pow((1-cos), 5)

}

func findN1AndN2(xs Intersections, target Intersection) (N1, N2 float64) {

	if xs.Count == 0 {
		return
	}

	var container []*Intersection

	getN := func(section *Intersection) float64 {

		var N float64
		//TimeもObjectも完全一致
		if section.Time == target.Time && section.Object == target.Object {
			if len(container) == 0 {
				N = 1.0
			} else {
				N = container[len(container)-1].Object.GetMaterial().RefractiveIndex
			}
		}

		return N
	}

	getContainerIndex := func(targetSection *Intersection) int {
		for i, section := range container {
			//Objectのみ一致
			if section.Object == targetSection.Object {
				return i
			}
		}

		return -1
	}

	for _, section := range xs.Intersections {
		N1 = getN(section)

		index := getContainerIndex(section)
		if index == -1 {
			container = append(container, section)
		} else {
			container = append(container[:index], container[index+1:]...)
		}
		N2 = getN(section)

		//targetが終わったらbreak
		if N1 != 0 || N2 != 0 {
			break
		}
	}

	return
}

func PrepareComputations(intersection Intersection, ray Ray, xs Intersections) (PreComps, error) {

	t := intersection.Time
	obj := intersection.Object
	ray_point := ray.Position(t)
	eye_vec := calc.NegTuple(ray.Direction)
	normal_vec, err := obj.NormalAt(ray_point, intersection)

	if err != nil {
		return PreComps{}, err
	}

	var IsRayInside bool

	//rayのOriginがObjectのInsideにあるとき
	if calc.DotTuple(normal_vec, eye_vec) < 0 {
		IsRayInside = true
		normal_vec = calc.NegTuple(normal_vec)
	}

	reflect_vec := calc.Reflect(ray.Direction, normal_vec)

	//shadow用にOverPointを作る,normal方向に微小に↓にずらしたものがoverpoint
	over_point := calc.AddTuple(ray_point, calc.MulTupleByScalar(util.DefaultEpsilon, normal_vec))
	under_point := calc.SubTuple(ray_point, calc.MulTupleByScalar(util.DefaultEpsilon, normal_vec))

	n1, n2 := findN1AndN2(xs, intersection)

	return PreComps{
		Time:        t,
		Object:      obj,
		RayPoint:    ray_point,
		OverPoint:   over_point,
		UnderPoint:  under_point,
		EyeVec:      eye_vec,
		NormalVec:   normal_vec,
		ReflectVec:  reflect_vec,
		IsRayInside: IsRayInside,
		N1:          n1,
		N2:          n2,
	}, nil
}
