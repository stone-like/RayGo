package scene

import (
	"math"
	"rayGo/calc"
)

//どんな図形でも中心がある場所は原点とする(ObjectCordination)
//図形を動かすのではなく、rayをコピーしてから動かすことで、各図形に合わせたrayを作る
type Sphere struct {
	*BaseShape
	r float64
}

var _ Shape = Sphere{}

func NewSphere(r float64) Sphere {
	return Sphere{
		NewBaseShape(),
		r,
	}
}

func (s Sphere) GetIntersection(time float64) *Intersection {
	return &Intersection{
		time:   time,
		object: s,
	}
}

func (s Sphere) Intersect(r Ray) (Intersections, error) {
	//spehereの代わりにrayをTransform
	invTrans, err := s.GetTransform().Inverse()
	if err != nil {
		return Intersections{}, err
	}

	r = r.Transform(invTrans)

	//sphereは原点にいるのでcalc.NewPoint(0, 0, 0)
	sphereToRay := calc.SubTuple(r.Origin, calc.NewPoint(0, 0, 0))

	a := calc.DotTuple(r.Direction, r.Direction)
	b := 2 * calc.DotTuple(r.Direction, sphereToRay)
	c := calc.DotTuple(sphereToRay, sphereToRay) - 1

	discriminant := math.Pow(b, 2) - 4*a*c

	if discriminant < 0 {
		return Intersections{}, nil
	}

	t1 := ((-b - math.Sqrt(discriminant)) / (2 * a))
	t2 := ((-b + math.Sqrt(discriminant)) / (2 * a))

	intersection1 := Intersection{
		time:   t1,
		object: s,
	}

	intersection2 := Intersection{
		time:   t2,
		object: s,
	}

	return Intersections{
		count:         2,
		Intersections: []*Intersection{&intersection1, &intersection2},
	}, nil
}