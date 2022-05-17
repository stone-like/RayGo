package scene

import (
	"rayGo/calc"
)

type CSGError struct {
	msg string
}

func (c CSGError) Error() string {
	return c.msg
}

func NewCSGError(msg string) CSGError {
	return CSGError{
		msg: msg,
	}
}

type CSG struct {
	*BaseShape
	Operation string
	Left      Shape
	Right     Shape
}

var _ Shape = &CSG{}

const (
	CSGUnion        = "union"
	CSGDifference   = "difference"
	CSGIntersection = "intersection"
)

func checkValidOperation(operation string) bool {
	switch operation {
	case CSGUnion,
		CSGDifference,
		CSGIntersection:
		return true
	default:
		return false
	}
}

func NewCSG(operation string, left, right Shape) (*CSG, error) {
	if !checkValidOperation(operation) {
		return nil, NewCSGError("operation is invalid")
	}
	csg := &CSG{
		NewBaseShape(),
		operation,
		left,
		right,
	}

	csg.setParentToShapes(left, right)

	return csg, nil
}

func (c *CSG) setParentToShapes(left, right Shape) {
	left.SetParent(c)
	right.SetParent(c)
}

func (c *CSG) isAllowedOnUnion(lhit, inl, inr bool) bool {
	return (lhit && !inr) || (!lhit && !inl)
}
func (c *CSG) isAllowedOnIntersection(lhit, inl, inr bool) bool {
	return (lhit && inr) || (!lhit && inl)
}
func (c *CSG) isAllowedOnDifference(lhit, inl, inr bool) bool {
	return (lhit && !inr) || (!lhit && inl)
}

func (c *CSG) checkIntersectionAllowed(lhit, inl, inr bool) bool {
	switch c.Operation {
	case CSGUnion:
		return c.isAllowedOnUnion(lhit, inl, inr)
	case CSGIntersection:
		return c.isAllowedOnIntersection(lhit, inl, inr)
	case CSGDifference:
		return c.isAllowedOnDifference(lhit, inl, inr)
	default:
		return false
	}
}

func (c *CSG) filterIntersections(xs Intersections) Intersections {

	var inl bool
	var inr bool
	var newSection []*Intersection

	for _, section := range xs.Intersections {
		lhit := c.Left.IsInclude(section.Object)

		if c.checkIntersectionAllowed(lhit, inl, inr) {
			newSection = append(newSection, section)
		}

		if lhit {
			inl = !inl
		} else {
			inr = !inr
		}
	}

	return AggregateIntersection(newSection...)
}

func (c *CSG) calcLocalNormal(localPoint calc.Tuple4, hit Intersection) calc.Tuple4 {
	return calc.NewVector(0, 1, 0)
}

func (c *CSG) NormalAt(worldPoint calc.Tuple4, hit Intersection) (calc.Tuple4, error) {
	return c.ShapeNormalAt(worldPoint, hit, c.calcLocalNormal)
}

func (c *CSG) calcLocalIntersect(r Ray) (Intersections, error) {
	leftXs, err := c.Left.Intersect(r)
	if err != nil {
		return Intersections{}, err
	}

	rightXs, err := c.Right.Intersect(r)
	if err != nil {
		return Intersections{}, err
	}

	var sections []*Intersection

	sections = append(sections, leftXs.Intersections...)
	sections = append(sections, rightXs.Intersections...)

	xs := AggregateIntersection(sections...)

	return c.filterIntersections(xs), nil
}

func (c CSG) Intersect(r Ray) (Intersections, error) {
	return c.ShapeIntersect(r, c.calcLocalIntersect)
}

func (c CSG) GetMaterial() *Material {
	return c.Material
}

func (c CSG) SetMaterial(m *Material) {
	c.Material = m
}

func (c CSG) IsInclude(s Shape) bool {

	for _, child := range []Shape{
		c.Left,
		c.Right,
	} {
		if child.IsInclude(s) {
			return true
		}
	}

	return false
}
