package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCSGCreatError(t *testing.T) {
	s := NewSphere(1)
	c := NewCube()

	_, err := NewCSG("unionError", s, c)
	require.Equal(t, "operation is invalid", err.Error())
}

type CSGRule struct {
	title     string
	operation string
	lhit      bool
	inl       bool
	inr       bool
	result    bool
}

func constructCSGRules() []CSGRule {
	return []CSGRule{
		{
			"union1",
			"union",
			true,
			true,
			true,
			false,
		},
		{
			"union2",
			"union",
			true,
			true,
			false,
			true,
		},
		{
			"union3",
			"union",
			true,
			false,
			true,
			false,
		},
		{
			"union4",
			"union",
			true,
			false,
			false,
			true,
		},
		{
			"union5",
			"union",
			false,
			true,
			true,
			false,
		},
		{
			"union6",
			"union",
			false,
			true,
			false,
			false,
		},
		{
			"union7",
			"union",
			false,
			false,
			true,
			true,
		},
		{
			"union8",
			"union",
			false,
			false,
			false,
			true,
		},
		{
			"intersection1",
			"intersection",
			true,
			true,
			true,
			true,
		},
		{
			"intersection2",
			"intersection",
			true,
			true,
			false,
			false,
		},
		{
			"intersection3",
			"intersection",
			true,
			false,
			true,
			true,
		},
		{
			"intersection4",
			"intersection",
			true,
			false,
			false,
			false,
		},
		{
			"intersection5",
			"intersection",
			false,
			true,
			true,
			true,
		},
		{
			"intersection6",
			"intersection",
			false,
			true,
			false,
			true,
		},
		{
			"intersection7",
			"intersection",
			false,
			false,
			true,
			false,
		},
		{
			"intersection8",
			"intersection",
			false,
			false,
			false,
			false,
		},
		{
			"difference1",
			"difference",
			true,
			true,
			true,
			false,
		},
		{
			"difference2",
			"difference",
			true,
			true,
			false,
			true,
		},
		{
			"difference3",
			"difference",
			true,
			false,
			true,
			false,
		},
		{
			"difference4",
			"difference",
			true,
			false,
			false,
			true,
		},
		{
			"difference5",
			"difference",
			false,
			true,
			true,
			true,
		},
		{
			"difference6",
			"difference",
			false,
			true,
			false,
			true,
		},
		{
			"difference7",
			"difference",
			false,
			false,
			true,
			false,
		},
		{
			"difference8",
			"difference",
			false,
			false,
			false,
			false,
		},
	}
}

func Test_Evaluating_Rule_Foe_CSG_Operation(t *testing.T) {

	s := NewSphere(1)
	c := NewCube()

	for _, rule := range constructCSGRules() {

		t.Run(rule.title, func(t *testing.T) {
			csg, err := NewCSG(rule.operation, s, c)
			require.Nil(t, err)

			require.Equal(t, rule.result, csg.checkIntersectionAllowed(rule.lhit, rule.inl, rule.inr))
		})
	}
}

type CSGFilter struct {
	title     string
	operation string
	x0Index   int
	x1Index   int
}

func constructCSGFilters() []CSGFilter {
	return []CSGFilter{
		{
			"filter1",
			"union",
			0,
			3,
		},
		{
			"filter2",
			"intersection",
			1,
			2,
		},
		{
			"filter3",
			"difference",
			0,
			1,
		},
	}
}

func Test_IsInclude_Work(t *testing.T) {
	s1 := NewSphere(1)
	s2 := NewSphere(1)

	require.True(t, s1.IsInclude(s1))
	require.False(t, s1.IsInclude(s2))

}

func Test_GroupIsInclude_Work(t *testing.T) {
	//group1
	//  group2
	//    s1
	//  s2
	s1 := NewSphere(1)
	s2 := NewSphere(1)
	s3 := NewSphere(1)

	group1 := NewGroup()
	group2 := NewGroup()

	group2.AddChildren(s1)

	group1.AddChildren(s2, group2)

	require.True(t, group1.IsInclude(s1))
	require.True(t, group1.IsInclude(s2))
	require.False(t, group1.IsInclude(s3))

}

func Test_Filtering_List_Of_Intersections(t *testing.T) {

	s1 := NewSphere(1)
	s2 := NewCube()

	xs := AggregateIntersection(
		&Intersection{1, s1, 0, 0},
		&Intersection{2, s2, 0, 0},
		&Intersection{3, s1, 0, 0},
		&Intersection{4, s2, 0, 0},
	)

	for _, filter := range constructCSGFilters() {

		t.Run(filter.title, func(t *testing.T) {
			csg, err := NewCSG(filter.operation, s1, s2)
			require.Nil(t, err)

			result := csg.filterIntersections(xs)

			require.Equal(t, 2, result.Count)
			require.Equal(t, xs.Intersections[filter.x0Index], result.Intersections[0])
			require.Equal(t, xs.Intersections[filter.x1Index], result.Intersections[1])

		})
	}
}

func Test_Ray_Misses_CSG_Object(t *testing.T) {

	s1 := NewSphere(1)
	s2 := NewCube()
	c, err := NewCSG(CSGUnion, s1, s2)
	require.Nil(t, err)

	ray := NewRay(calc.NewPoint(0, 2, -5), calc.NewVector(0, 0, 1))
	xs, err := c.calcLocalIntersect(ray)
	require.Nil(t, err)

	require.Equal(t, 0, xs.Count)
}

func Test_Ray_Hits_CSG_Object(t *testing.T) {

	s1 := NewSphere(1)
	s2 := NewSphere(1)

	s2.SetTransform(calc.NewTranslation(0, 0, 0.5))

	c, err := NewCSG(CSGUnion, s1, s2)
	require.Nil(t, err)

	ray := NewRay(calc.NewPoint(0, 0, -5), calc.NewVector(0, 0, 1))
	xs, err := c.calcLocalIntersect(ray)
	require.Nil(t, err)

	require.Equal(t, 2, xs.Count)
	require.Equal(t, float64(4), xs.Intersections[0].Time)
	require.Equal(t, s1, xs.Intersections[0].Object)
	require.Equal(t, float64(6.5), xs.Intersections[1].Time)
	require.Equal(t, s2, xs.Intersections[1].Object)

}
