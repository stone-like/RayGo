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
		Time:   time,
		Object: s,
	}
}

func (s Sphere) calcLocalNormal(localPoint calc.Tuple4, hit Intersection) calc.Tuple4 {
	return calc.SubTuple(localPoint, calc.NewPoint(0, 0, 0))
}

func (s Sphere) NormalAt(worldPoint calc.Tuple4, hit Intersection) (calc.Tuple4, error) {
	return s.ShapeNormalAt(worldPoint, hit, s.calcLocalNormal)
}

func (s Sphere) calcLocalIntersect(r Ray) (Intersections, error) {
	//shapeは原点にいるのでcalc.NewPoint(0, 0, 0)
	shapeToRay := calc.SubTuple(r.Origin, calc.NewPoint(0, 0, 0))

	a := calc.DotTuple(r.Direction, r.Direction)
	b := 2 * calc.DotTuple(r.Direction, shapeToRay)
	c := calc.DotTuple(shapeToRay, shapeToRay) - 1

	discriminant := math.Pow(b, 2) - 4*a*c

	if discriminant < 0 {
		return Intersections{}, nil
	}

	t1 := ((-b - math.Sqrt(discriminant)) / (2 * a))
	t2 := ((-b + math.Sqrt(discriminant)) / (2 * a))

	intersection1 := Intersection{
		Time:   t1,
		Object: s,
	}

	intersection2 := Intersection{
		Time:   t2,
		Object: s,
	}

	return Intersections{
		Count:         2,
		Intersections: []*Intersection{&intersection1, &intersection2},
	}, nil
}

func (s Sphere) Intersect(r Ray) (Intersections, error) {
	return s.ShapeIntersect(r, s.calcLocalIntersect)
}

func (s Sphere) GetMaterial() *Material {
	return s.Material
}

func (s Sphere) SetMaterial(m *Material) {
	s.Material = m
}

// //sphere自体を動かす代わりにNormalを動かして計算
// func (s Sphere) NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error) {
// 	invTrans, err := s.GetTransform().Inverse()
// 	if err != nil {
// 		return calc.Tuple4{}, err
// 	}

// 	objectPoint := invTrans.MulByTuple(worldPoint)
// 	objectNormal := calc.SubTuple(objectPoint, calc.NewPoint(0, 0, 0))

// 	//objectNormal -> worldNormalでなぜinverse->TransposeがいるかはPDFに記載
// 	worldNormal := invTrans.Transpose().MulByTuple(objectNormal)
// 	worldNormal[3] = 0

// 	return worldNormal.Normalize(), nil
// }

// func (s Sphere) Intersect(r Ray) (Intersections, error) {
// 	//spehereの代わりにrayをTransform
// 	invTrans, err := s.GetTransform().Inverse()
// 	if err != nil {
// 		return Intersections{}, err
// 	}

// 	r = r.Transform(invTrans)

// 	//sphereは原点にいるのでcalc.NewPoint(0, 0, 0)
// 	sphereToRay := calc.SubTuple(r.Origin, calc.NewPoint(0, 0, 0))

// 	a := calc.DotTuple(r.Direction, r.Direction)
// 	b := 2 * calc.DotTuple(r.Direction, sphereToRay)
// 	c := calc.DotTuple(sphereToRay, sphereToRay) - 1

// 	discriminant := math.Pow(b, 2) - 4*a*c

// 	if discriminant < 0 {
// 		return Intersections{}, nil
// 	}

// 	t1 := ((-b - math.Sqrt(discriminant)) / (2 * a))
// 	t2 := ((-b + math.Sqrt(discriminant)) / (2 * a))

// 	intersection1 := Intersection{
// 		Time:   t1,
// 		Object: s,
// 	}

// 	intersection2 := Intersection{
// 		Time:   t2,
// 		Object: s,
// 	}

// 	return Intersections{
// 		Count:         2,
// 		Intersections: []*Intersection{&intersection1, &intersection2},
// 	}, nil
// }
