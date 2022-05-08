package scene

import (
	"math"
	"rayGo/calc"
	"rayGo/util"
)

type Options struct {
	R      float64
	Min    float64
	Max    float64
	Closed bool
}

type Option func(*Options)

func CynR(r float64) Option {
	return func(o *Options) {
		o.R = r
	}
}

func CynMin(m float64) Option {
	return func(o *Options) {
		o.Min = m
	}
}

func CynMax(m float64) Option {
	return func(o *Options) {
		o.Max = m
	}
}

func CynClosed(isClosed bool) Option {
	return func(o *Options) {
		o.Closed = isClosed
	}
}

//Minは-y、Maxはyをそれぞれ指定
type Cyliner struct {
	*BaseShape
	r      float64
	Min    float64
	Max    float64
	Closed bool
}

var _ Shape = Cyliner{}

func NewCyliner(options ...Option) Cyliner {

	defaultOptions := &Options{
		1,
		-util.Inf,
		util.Inf,
		false,
	}

	for _, fn := range options {
		fn(defaultOptions)
	}

	return Cyliner{
		NewBaseShape(),
		defaultOptions.R,
		defaultOptions.Min,
		defaultOptions.Max,
		defaultOptions.Closed,
	}

}

func (c Cyliner) calcLocalNormal(localPoint calc.Tuple4) calc.Tuple4 {

	dist := math.Pow(localPoint[0], 2) + math.Pow(localPoint[2], 2)

	//上面と下面のNormal
	y := localPoint[1]
	//上面
	if dist < 1 && y >= c.Max-util.EPSILON {
		return calc.NewVector(0, 1, 0)
	}
	//下面
	if dist < 1 && y <= c.Min+util.EPSILON {
		return calc.NewVector(0, -1, 0)
	}

	//側面
	return calc.NewVector(localPoint[0], 0, localPoint[2])

}

func (c Cyliner) NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error) {
	return c.ShapeNormalAt(worldPoint, c.calcLocalNormal)
}

func calcT0AndT1(a, b, ans float64) (float64, float64) {
	t0 := (-b - math.Sqrt(ans)) / (2 * a)
	t1 := (-b + math.Sqrt(ans)) / (2 * a)

	if t0 > t1 {
		t0, t1 = t1, t0
	}

	return t0, t1

}

type YTPair struct {
	y float64
	t float64
}

func (c Cyliner) calcXs(t0, t1 float64, ray Ray) Intersections {
	var xs []*Intersection

	y0 := ray.Origin[1] + t0*ray.Direction[1]
	y1 := ray.Origin[1] + t1*ray.Direction[1]

	for _, pair := range []YTPair{{y0, t0}, {y1, t1}} {
		if c.Min < pair.y && pair.y < c.Max {
			xs = append(xs, &Intersection{
				Time:   pair.t,
				Object: c,
			})
		}
	}

	return c.intersectCaps(ray, xs)

}

func (c Cyliner) checkCap(ray Ray, t float64) bool {
	x := ray.Origin[0] + t*ray.Direction[0]
	z := ray.Origin[2] + t*ray.Direction[2]

	return (math.Pow(x, 2) + math.Pow(z, 2)) <= 1

}

func (c Cyliner) intersectCaps(ray Ray, sections []*Intersection) Intersections {

	//cyn does not have caps or rayDirection.y is cloase to zero
	if !c.Closed || util.IsNearlyEqualZero(ray.Direction[1]) {
		return AggregateIntersection(sections...)
	}

	mint, maxt := (c.Min-ray.Origin[1])/ray.Direction[1], (c.Max-ray.Origin[1])/ray.Direction[1]

	for _, t := range []float64{mint, maxt} {
		if c.checkCap(ray, t) {
			sections = append(sections, &Intersection{Time: t, Object: c})
		}
	}

	return AggregateIntersection(sections...)
}

func (c Cyliner) calcLocalIntersect(r Ray) (Intersections, error) {
	_a := math.Pow(r.Direction[0], 2) + math.Pow(r.Direction[2], 2)

	//ray is parallel to y axis, because ray.Direction.X and ray.Direction.Z are nearlyEqualZero, r.Direction is composed almost Y component
	//上記の状況でほぼY要素しかなかったとしてもcylinderの上面と下面のCapに接しているときがるのでintersectCapを確かめる
	if util.IsNearlyEqualZero(_a) {
		return c.intersectCaps(r, []*Intersection{}), nil
	}

	_b := 2*(r.Origin[0]*r.Direction[0]) + 2*(r.Origin[2]*r.Direction[2])
	_c := math.Pow(r.Origin[0], 2) + math.Pow(r.Origin[2], 2) - 1

	ans := math.Pow(_b, 2) - 4*_a*_c

	if ans < 0 {
		return Intersections{}, nil
	}

	t0, t1 := calcT0AndT1(_a, _b, ans)

	xs := c.calcXs(t0, t1, r)

	return xs, nil
}

func (c Cyliner) Intersect(r Ray) (Intersections, error) {
	return c.ShapeIntersect(r, c.calcLocalIntersect)
}

func (c Cyliner) GetMaterial() *Material {
	return c.Material
}

func (c Cyliner) SetMaterial(m *Material) {
	c.Material = m
}
