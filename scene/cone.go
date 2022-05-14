package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

func ConeR(r float64) Option {
	return func(o *Options) {
		o.R = r
	}
}

func ConeMin(m float64) Option {
	return func(o *Options) {
		o.Min = m
	}
}

func ConeMax(m float64) Option {
	return func(o *Options) {
		o.Max = m
	}
}

func ConeClosed(isClosed bool) Option {
	return func(o *Options) {
		o.Closed = isClosed
	}
}

//Minは-y、Maxはyをそれぞれ指定
type Cone struct {
	*BaseShape
	r      float64
	Min    float64
	Max    float64
	Closed bool
}

var _ Shape = Cone{}

func NewCone(options ...Option) Cone {

	defaultOptions := &Options{
		1,
		-util.Inf,
		util.Inf,
		false,
	}

	for _, fn := range options {
		fn(defaultOptions)
	}

	return Cone{
		NewBaseShape(),
		defaultOptions.R,
		defaultOptions.Min,
		defaultOptions.Max,
		defaultOptions.Closed,
	}

}

func (c Cone) calcLocalNormal(localPoint calc.Tuple4, hit Intersection) calc.Tuple4 {

	dist := math.Pow(localPoint[0], 2) + math.Pow(localPoint[2], 2)

	//上面と下面のNormal
	y := localPoint[1]
	//上面
	if dist < 1 && y >= c.Max-util.DefaultEpsilon {
		return calc.NewVector(0, 1, 0)
	}
	//下面
	if dist < 1 && y <= c.Min+util.DefaultEpsilon {
		return calc.NewVector(0, -1, 0)
	}

	//側面
	yComponent := math.Sqrt(dist)

	if y > 0 {
		yComponent = -yComponent
	}

	return calc.NewVector(localPoint[0], yComponent, localPoint[2])

}

func (c Cone) NormalAt(worldPoint calc.Tuple4, hit Intersection) (calc.Tuple4, error) {
	return c.ShapeNormalAt(worldPoint, hit, c.calcLocalNormal)
}

func (c Cone) calcXs(ray Ray, ts []float64) Intersections {
	var xs []*Intersection
	var pairs []YTPair

	for _, time := range ts {
		y := ray.Origin[1] + time*ray.Direction[1]
		pairs = append(pairs, YTPair{y, time})
	}

	for _, pair := range pairs {
		if c.Min < pair.y && pair.y < c.Max {
			xs = append(xs, &Intersection{
				Time:   pair.t,
				Object: c,
			})
		}
	}

	return c.intersectCaps(ray, xs)

}

func (c Cone) intersectCaps(ray Ray, sections []*Intersection) Intersections {

	if !c.Closed || util.IsNearlyEqualZero(ray.Direction[1]) {
		return AggregateIntersection(sections...)
	}

	mint, maxt := (c.Min-ray.Origin[1])/ray.Direction[1], (c.Max-ray.Origin[1])/ray.Direction[1]

	for _, pair := range []YTPair{{c.Min, mint}, {c.Max, maxt}} {
		if c.checkCap(ray, pair.t, pair.y) {
			sections = append(sections, &Intersection{Time: pair.t, Object: c})
		}
	}

	return AggregateIntersection(sections...)
}

func (c Cone) checkCap(ray Ray, t, y float64) bool {
	x := ray.Origin[0] + t*ray.Direction[0]
	z := ray.Origin[2] + t*ray.Direction[2]

	return (math.Pow(x, 2) + math.Pow(z, 2)) <= math.Abs(y)
}

func calcTimes(a, b, c, ans float64) []float64 {

	if util.IsNearlyEqualZero(a) && !util.IsNearlyEqualZero(b) {
		t := -c / (2.0 * b)
		return []float64{t}
	}

	t0 := (-b - math.Sqrt(ans)) / (2 * a)
	t1 := (-b + math.Sqrt(ans)) / (2 * a)

	return []float64{t0, t1}

}

func (c Cone) calcLocalIntersect(r Ray) (Intersections, error) {

	_a := math.Pow(r.Direction[0], 2) - math.Pow(r.Direction[1], 2) + math.Pow(r.Direction[2], 2)
	_b := 2*(r.Origin[0]*r.Direction[0]) - 2*(r.Origin[1]*r.Direction[1]) + 2*(r.Origin[2]*r.Direction[2])
	_c := math.Pow(r.Origin[0], 2) - math.Pow(r.Origin[1], 2) + math.Pow(r.Origin[2], 2)

	if util.IsNearlyEqualZero(_a) && util.IsNearlyEqualZero(_b) {
		return Intersections{}, nil
	}

	ans := math.Pow(_b, 2) - 4*_a*_c
	//解なしなので交差せず
	if ans < 0 {
		return Intersections{}, nil
	}

	ts := calcTimes(_a, _b, _c, ans)

	xs := c.calcXs(r, ts)

	return xs, nil
}

func (c Cone) Intersect(r Ray) (Intersections, error) {
	return c.ShapeIntersect(r, c.calcLocalIntersect)
}

func (c Cone) GetMaterial() *Material {
	return c.Material
}

func (c Cone) SetMaterial(m *Material) {
	c.Material = m
}
