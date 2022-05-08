package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

type Cube struct {
	*BaseShape
}

var _ Shape = Cube{}

func NewCube() Cube {
	return Cube{
		NewBaseShape(),
	}
}

func (c Cube) calcLocalNormal(localPoint calc.Tuple4) calc.Tuple4 {
	x := math.Abs(localPoint[0])
	y := math.Abs(localPoint[1])
	z := math.Abs(localPoint[2])

	maxc := math.Max(x, math.Max(y, z))

	if maxc == x {
		return calc.NewVector(localPoint[0], 0, 0)
	}

	if maxc == y {
		return calc.NewVector(0, localPoint[1], 0)
	}

	return calc.NewVector(0, 0, localPoint[2])
}

func (c Cube) NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error) {
	return c.ShapeNormalAt(worldPoint, c.calcLocalNormal)
}

func calcTminAndTmax(tminNumerator, tmaxNumerator, directionComponent float64) (float64, float64) {

	//divisionByZero対策
	if util.IsNearlyEqualZero(directionComponent) {
		return tminNumerator * math.Inf(0), tmaxNumerator * math.Inf(0)
	}

	return tminNumerator / directionComponent, tmaxNumerator / directionComponent
}

func checkAxis(originComponent, directionComponent float64) (float64, float64) {

	tmin_numerator := (-1 - originComponent)
	tmax_numerator := (1 - originComponent)

	tmin, tmax := calcTminAndTmax(tmin_numerator, tmax_numerator, directionComponent)

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	return tmin, tmax
}

func (c Cube) calcLocalIntersect(r Ray) (Intersections, error) {
	xtmin, xtmax := checkAxis(r.Origin[0], r.Direction[0])
	ytmin, ytmax := checkAxis(r.Origin[1], r.Direction[1])
	ztmin, ztmax := checkAxis(r.Origin[2], r.Direction[2])

	tmin := math.Max(xtmin, math.Max(ytmin, ztmin))
	tmax := math.Min(xtmax, math.Min(ytmax, ztmax))

	//ray Missing
	if tmin > tmax {
		return Intersections{}, nil
	}

	i1 := Intersection{
		Time:   tmin,
		Object: c,
	}
	i2 := Intersection{
		Time:   tmax,
		Object: c,
	}

	return AggregateIntersection(&i1, &i2), nil

}

func (c Cube) Intersect(r Ray) (Intersections, error) {
	return c.ShapeIntersect(r, c.calcLocalIntersect)
}

func (c Cube) GetMaterial() *Material {
	return c.Material
}

func (c Cube) SetMaterial(m *Material) {
	c.Material = m
}
