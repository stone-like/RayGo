package scene

import (
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

	c, err := w.ShadeHit(comps)
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

	c, err := w.ShadeHit(comps)
	require.Nil(t, err)

	require.True(t, colorCompare(NewColor(0.90498, 0.90498, 0.90498), c))

}

func Test_Color_At_When_Ray_Misses(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 1, 0))

	c, err := w.ColorAt(r)

	require.Nil(t, err)
	require.True(t, colorCompare(NewColor(0, 0, 0), c))

}

func Test_Color_At_When_Ray_Hits(t *testing.T) {
	w := DefaultWorld()
	r := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))

	c, err := w.ColorAt(r)

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

	c, err := w.ColorAt(r)

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

	require.True(t, comps.OverPoint[2] < -(util.EPSILON/2))
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

	c, err := w.ShadeHit(comps)
	require.Nil(t, err)
	require.True(t, colorCompare(NewColor(0.1, 0.1, 0.1), c))
}
