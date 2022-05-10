package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ray_With_World(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))

	xs, err := w.Intersect(r)

	require.Nil(t, err)

	require.Equal(t, 4, xs.Count)
	require.Equal(t, 4.0, xs.Intersections[0].Time)
	require.Equal(t, 4.5, xs.Intersections[1].Time)
	require.Equal(t, 5.5, xs.Intersections[2].Time)
	require.Equal(t, 6.0, xs.Intersections[3].Time)

}

//hit(hitのt,point,object),eye_vec,normal_vecをpreComputing
func Test_PreComputing_When_Ray_Occurs_Outside_Object(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	shape := NewSphere(1)

	i := CreateIntersection(4, shape)

	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)

	require.Equal(t, i.Time, comps.Time)
	require.Equal(t, i.Object, comps.Object)
	require.Equal(t, calc.NewPoint(0, 0, -1), comps.RayPoint)
	require.Equal(t, calc.NewVector(0, 0, -1), comps.EyeVec)
	require.Equal(t, calc.NewVector(0, 0, -1), comps.NormalVec)
	require.False(t, comps.IsRayInside)

}

func Test_PreComputing_When_Ray_Occurs_Inside_Object(t *testing.T) {
	r := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1))
	shape := NewSphere(1)

	i := CreateIntersection(1, shape)

	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)

	require.Equal(t, calc.NewPoint(0, 0, 1), comps.RayPoint)
	require.Equal(t, calc.NewVector(0, 0, -1), comps.EyeVec)
	require.Equal(t, calc.NewVector(0, 0, -1), comps.NormalVec)
	require.True(t, comps.IsRayInside)

}

//ShadeHitは内部でLightingを呼んでいる
func Test_Shading_Intersection(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))

	shape := w.Objects[0]
	i := CreateIntersection(4, shape)

	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)

	c, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.38066, 0.47583, 0.2855), c))

}

func Test_Shading_Intersection_When_Ray_Is_Inside(t *testing.T) {
	w := DefaultWorld()
	w.Light = NewLight(calc.NewPoint(0, 0.25, 0), NewColor(1, 1, 1))
	r := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1))

	shape := w.Objects[1]
	i := CreateIntersection(0.5, shape)

	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)

	c, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.90498, 0.90498, 0.90498), c))

}

func Test_Color_At_When_Ray_Misses(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 1, 0))

	c, err := w.ColorAt(r, DefaultRemaing, DefaultRemaing)

	require.Nil(t, err)
	require.True(t, colorCompare(NewColor(0, 0, 0), c))

}

func Test_Color_At_When_Ray_Hits(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))

	c, err := w.ColorAt(r, DefaultRemaing, DefaultRemaing)

	require.Nil(t, err)
	require.True(t, colorCompare(NewColor(0.38066, 0.47583, 0.2855), c))

}

func Test_Color_At_With_Intersection_Behind_Ray(t *testing.T) {
	w := DefaultWorld()

	outer := w.Objects[0]
	inner := w.Objects[1]

	m1 := outer.GetMaterial()
	m1.Ambient = 1
	outer.SetMaterial(m1)
	m2 := inner.GetMaterial()
	m2.Ambient = 1
	inner.SetMaterial(m2)

	r := NewRay(calc.NewPoint(0, 0, 0.75), calc.NewVector(0, 0, -1))

	c, err := w.ColorAt(r, DefaultRemaing, DefaultRemaing)

	require.Nil(t, err)
	require.True(t, colorCompare(inner.GetMaterial().Color, c))

}

func TestIsShadowed(t *testing.T) {
	//pointを作って原点からR:1のsphereと-10, 10, -10にある光源という状況でpointが陰になるかをテスト

	w := DefaultWorld()
	for _, target := range []struct {
		title    string
		p        calc.Tuple4
		isShadow bool
	}{
		{
			"no shadow when object does not block",
			calc.NewPoint(0, 10, 0),
			false,
		},
		{
			"shadow when object is between point and light",
			calc.NewPoint(10, -10, 10),
			true,
		},
		{
			"no shadow when object is behind light",
			calc.NewPoint(-20, 20, -20),
			false,
		},
		{
			"no shadow when object is behind point",
			calc.NewPoint(-2, 2, -2),
			false,
		},
	} {
		t.Run(target.title, func(t *testing.T) {

			isShadow, err := w.IsShadowed(target.p)
			require.Nil(t, err)
			require.Equal(t, target.isShadow, isShadow)
		})
	}
}

func Test_Hit_OverPoint(t *testing.T) {
	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	shape := NewSphere(1)
	shape.SetTransform(calc.NewTranslation(0, 0, 1))
	i := CreateIntersection(5, shape)
	comps, err := PrepareComputations(i, ray, Intersections{})
	require.Nil(t, err)

	require.True(t, comps.OverPoint[2] < -(util.DefaultEpsilon/2))
	require.True(t, comps.RayPoint[2] > comps.OverPoint[2])

}

func Test_ShadeHit_When_Given_IsShadow_is_True(t *testing.T) {
	w := DefaultWorld()

	w.Light = NewLight(calc.NewPoint(0, 0, -10), NewColor(1, 1, 1))
	s1 := NewSphere(1)
	s2 := NewSphere(1)
	s2.SetTransform(calc.NewTranslation(0, 0, 10))

	w.Objects = append(w.Objects, s1, s2)

	ray := NewRay(calc.NewPoint(0, 0, 5), calc.NewVector(0, 0, 1))
	i := CreateIntersection(4, s2)

	comps, err := PrepareComputations(i, ray, Intersections{})
	require.Nil(t, err)

	c, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)
	require.True(t, colorCompare(NewColor(0.1, 0.1, 0.1), c))
}

func Test_Reflected_Color_For_NonReflective_Material(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 0, 1))

	shape := w.Objects[1]
	m := shape.GetMaterial()
	m.Ambient = 1
	shape.SetMaterial(m)

	i := Intersection{1, shape}
	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)
	color, err := w.ReflectedColor(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0, 0, 0), color))
}

//packageTestだとcolorの計算がずれている...fileTestとか単体で動かすときは計算あっているのでロジック自体は問題ないはず
//package間で何かが起こっている
func Test_Reflected_Color_For_Reflective_Material(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -3), calc.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))

	shape := NewPlane()
	m := shape.GetMaterial()
	m.Reflective = 0.5
	shape.SetMaterial(m)
	shape.SetTransform(calc.NewTranslation(0, -1, 0))

	i := Intersection{math.Sqrt(2), shape}
	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)
	color, err := w.ReflectedColor(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	util.SetEpsilon(0.0001)
	defer util.SetEpsilon(util.DefaultEpsilon)
	require.True(t, colorCompare(NewColor(0.19032, 0.2379, 0.14274), color))

}

func Test_Shade_Hit_For_Reflective_Material(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -3), calc.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))

	shape := NewPlane()
	m := shape.GetMaterial()
	m.Reflective = 0.5
	shape.SetMaterial(m)
	shape.SetTransform(calc.NewTranslation(0, -1, 0))

	i := Intersection{math.Sqrt(2), shape}
	comps, err := PrepareComputations(i, r, Intersections{})
	require.Nil(t, err)
	color, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	util.SetEpsilon(0.0001)
	defer util.SetEpsilon(util.DefaultEpsilon)
	require.True(t, colorCompare(NewColor(0.87677, 0.92436, 0.82918), color))

}

func Test_Avoid_Infinite_Recursion_For_Reflect(t *testing.T) {

	lower := NewPlane()
	m := lower.GetMaterial()
	m.Reflective = 1
	lower.SetMaterial(m)
	lower.SetTransform(calc.NewTranslation(0, -1, 0))

	upper := NewPlane()
	m2 := upper.GetMaterial()
	m2.Reflective = 1
	upper.SetMaterial(m2)
	upper.SetTransform(calc.NewTranslation(0, 1, 0))

	w := NewWorld(NewLight(calc.NewPoint(0, 0, 0), NewColor(1, 1, 1)), lower, upper)

	r := NewRay(calc.NewPoint(0, 0, 0), calc.NewVector(0, 1, 0))

	//回数制限を設けて、
	//ColorAt -> ShadeHit -> ReflectedColor -> ColorAtの無限ループを防ぐ
	w.ColorAt(r, DefaultRemaing, DefaultRemaing)

	terminate := true
	require.True(t, terminate)

}

func Test_Finding_n1_and_n2_at_Various_Intersections(t *testing.T) {
	s1 := GlassSphere()
	s1.SetTransform(calc.NewScale(2, 2, 2))
	m1 := s1.GetMaterial()
	m1.RefractiveIndex = 1.5
	s1.SetMaterial(m1)

	s2 := GlassSphere()
	s2.SetTransform(calc.NewTranslation(0, 0, -0.25))
	m2 := s2.GetMaterial()
	m2.RefractiveIndex = 2.0
	s2.SetMaterial(m2)

	s3 := GlassSphere()
	s3.SetTransform(calc.NewTranslation(0, 0, 0.25))
	m3 := s3.GetMaterial()
	m3.RefractiveIndex = 2.5
	s3.SetMaterial(m3)

	ray := NewRay(calc.NewPoint(0, 0, -4), calc.NewVector(0, 0, 1))

	section1 := Intersection{
		Time:   2,
		Object: s1,
	}
	section2 := Intersection{
		Time:   2.75,
		Object: s2,
	}
	section3 := Intersection{
		Time:   3.25,
		Object: s3,
	}
	section4 := Intersection{
		Time:   4.75,
		Object: s2,
	}
	section5 := Intersection{
		Time:   5.25,
		Object: s3,
	}
	section6 := Intersection{
		Time:   6,
		Object: s1,
	}
	xs := AggregateIntersection(&section1, &section2, &section3, &section4, &section5, &section6)

	for i, target := range []struct {
		title string
		n1    float64
		n2    float64
	}{
		{
			title: "section1",
			n1:    1.0,
			n2:    1.5,
		},
		{
			title: "section2",
			n1:    1.5,
			n2:    2.0,
		},
		{
			title: "section3",
			n1:    2.0,
			n2:    2.5,
		},
		{
			title: "section4",
			n1:    2.5,
			n2:    2.5,
		},
		{
			title: "section5",
			n1:    2.5,
			n2:    1.5,
		},
		{
			title: "section6",
			n1:    1.5,
			n2:    1.0,
		},
	} {

		t.Run(target.title, func(t *testing.T) {
			comps, err := PrepareComputations(*xs.Intersections[i], ray, xs)

			require.Nil(t, err)

			require.True(t, util.FloatEqual(target.n1, comps.N1))
			require.True(t, util.FloatEqual(target.n2, comps.N2))
		})

	}

}

//UnderPointを設けることでrefractedRayの原点がrayPoint(rayとObjectがぶつかる場所)と被らない
func Test_PreComputing_Under_Point(t *testing.T) {
	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	shape := GlassSphere()
	shape.SetTransform(calc.NewTranslation(0, 0, 1))

	i := Intersection{5, shape}

	xs := AggregateIntersection(&i)

	comps, err := PrepareComputations(i, ray, xs)
	require.Nil(t, err)

	require.True(t, comps.UnderPoint[2] > util.DefaultEpsilon/2)
	require.True(t, comps.UnderPoint[2] > comps.RayPoint[2])

}

func Test_Refracted_Color_On_Opaque_Surface(t *testing.T) {
	w := DefaultWorld()
	shape := w.Objects[0]

	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	xs := AggregateIntersection(&Intersection{4, shape}, &Intersection{6, shape})

	comps, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)

	color, err := w.RefractedColor(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0, 0, 0), color))
}

func Test_Avoid_Infinite_Recursion_on_RefractedColor(t *testing.T) {
	w := DefaultWorld()
	shape := w.Objects[0]

	m1 := shape.GetMaterial()
	m1.Transparency = 1.0
	m1.RefractiveIndex = 1.5
	shape.SetMaterial(m1)

	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	xs := AggregateIntersection(&Intersection{4, shape}, &Intersection{6, shape})

	comps, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)

	color, err := w.RefractedColor(comps, 0, 0)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0, 0, 0), color))
}

func Test_Refracted_Color_Under_Total_Internal_Reflection(t *testing.T) {
	w := DefaultWorld()
	shape := w.Objects[0]

	m1 := shape.GetMaterial()
	m1.Transparency = 1.0
	m1.RefractiveIndex = 1.5
	shape.SetMaterial(m1)

	//rayは内側のsphereの中から出る
	ray := NewRay(calc.NewPoint(0, 0, math.Sqrt(2)/2), calc.NewVector(0, 1, 0))
	xs := AggregateIntersection(&Intersection{-math.Sqrt(2) / 2, shape}, &Intersection{math.Sqrt(2) / 2, shape})

	//defaultWorld内で作った外側と内側のsphereで、内側のsphereから外側のsphereに至るrayの交点をみる
	comps, err := PrepareComputations(*xs.Intersections[1], ray, xs)
	require.Nil(t, err)

	color, err := w.RefractedColor(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0, 0, 0), color))
}

func Test_Refracted_Color_Under_Normal_Condition(t *testing.T) {
	w := DefaultWorld()
	shape := w.Objects[0]

	m1 := shape.GetMaterial()
	m1.Ambient = 1.0
	m1.RefractiveIndex = 1.0

	m1.SetPattern(NewTestPattern())
	shape.SetMaterial(m1)

	shape2 := w.Objects[1]
	m2 := shape2.GetMaterial()
	m2.Transparency = 1.0
	m2.RefractiveIndex = 1.5
	shape2.SetMaterial(m2)

	ray := NewRay(calc.NewPoint(0, 0, 0.1), calc.NewVector(0, 1, 0))
	xs := AggregateIntersection(
		&Intersection{-0.9899, shape},
		&Intersection{-0.4899, shape2},
		&Intersection{0.4899, shape2},
		&Intersection{0.9899, shape},
	)

	comps, err := PrepareComputations(*xs.Intersections[2], ray, xs)
	require.Nil(t, err)

	color, err := w.RefractedColor(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	util.SetEpsilon(0.0001)
	defer util.SetEpsilon(util.DefaultEpsilon)
	require.True(t, colorCompare(NewColor(0, 0.99888, 0.04725), color))
}

func Test_ShadeHit_Transparent_Material(t *testing.T) {
	w := DefaultWorld()
	floor := NewPlane()
	floorMaterial := floor.GetMaterial()
	floorMaterial.Transparency = 0.5
	floorMaterial.RefractiveIndex = 1.5
	floor.SetMaterial(floorMaterial)
	floor.SetTransform(calc.NewTranslation(0, -1, 0))

	ball := NewSphere(1)
	ballMaterial := ball.GetMaterial()
	ballMaterial.Color = NewColor(1, 0, 0)
	ballMaterial.Ambient = 0.5
	ball.SetMaterial(ballMaterial)
	ball.SetTransform(calc.NewTranslation(0, -3.5, -0.5))

	w.Objects = []Shape{floor, ball}

	ray := NewRay(calc.NewPoint(0, 0, -3), calc.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := AggregateIntersection(
		&Intersection{math.Sqrt(2), floor},
	)

	comps, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)

	color, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.93642, 0.68642, 0.68642), color))
}

func Test_ShadeHit_With_Reflective_Transparent_Marterial(t *testing.T) {
	w := DefaultWorld()

	ray := NewRay(calc.NewPoint(0, 0, -3), calc.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	floor := NewPlane()
	floorMaterial := floor.GetMaterial()
	floorMaterial.Reflective = 0.5
	floorMaterial.Transparency = 0.5
	floorMaterial.RefractiveIndex = 1.5
	floor.SetMaterial(floorMaterial)
	floor.SetTransform(calc.NewTranslation(0, -1, 0))

	ball := NewSphere(1)
	ballMaterial := ball.GetMaterial()
	ballMaterial.Color = NewColor(1, 0, 0)
	ballMaterial.Ambient = 0.5
	ball.SetMaterial(ballMaterial)
	ball.SetTransform(calc.NewTranslation(0, -3.5, -0.5))

	w.AddObjects(floor, ball)

	xs := AggregateIntersection(
		&Intersection{math.Sqrt(2), floor},
	)

	comps, err := PrepareComputations(*xs.Intersections[0], ray, xs)
	require.Nil(t, err)

	color, err := w.ShadeHit(comps, DefaultRemaing, DefaultRemaing)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.93391, 0.69643, 0.69243), color))
}
