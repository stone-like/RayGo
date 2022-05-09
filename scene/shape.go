package scene

import (
	"rayGo/calc"
	"sort"
)

type Shape interface {
	Intersect(r Ray) (Intersections, error)
	NormalAt(worldPoint calc.Tuple4) (calc.Tuple4, error)
	GetMaterial() *Material
	SetMaterial(m *Material)
	GetTransform() calc.Mat4x4
	SetTransform(mat calc.Mat4x4)
	GetParent() Shape //parentはそのshapeが属するGroupを表す
	SetParent(s Shape)
	WorldToObject(point calc.Tuple4) (calc.Tuple4, error)
	NormalToWorld(normal_vec calc.Tuple4) (calc.Tuple4, error)
}

type BaseShape struct {
	Transform calc.Mat4x4
	Material  *Material
	Parent    Shape
}

func NewBaseShape() *BaseShape {
	return &BaseShape{
		Transform: calc.Ident4x4,
		Material:  DefaultMaterial(),
		Parent:    nil,
	}
}

func (b *BaseShape) GetTransform() calc.Mat4x4 {
	return b.Transform
}

func (b *BaseShape) SetTransform(mat calc.Mat4x4) {
	b.Transform = mat
}

func (b *BaseShape) GetParent() Shape {
	return b.Parent
}

func (b *BaseShape) SetParent(s Shape) {
	b.Parent = s
}

type CalcLocalNormal func(localPoint calc.Tuple4) calc.Tuple4

//各shapeごとにnormalだったりintersectを求める方法が違うのでそこはfuncで引数経由で渡せばいい
func (base *BaseShape) ShapeNormalAt(worldPoint calc.Tuple4, calcLocalNormal CalcLocalNormal) (calc.Tuple4, error) {
	invTrans, err := base.GetTransform().Inverse()
	if err != nil {
		return calc.Tuple4{}, err
	}

	localPoint := invTrans.MulByTuple(worldPoint)
	localNormal := calcLocalNormal(localPoint)

	//objectNormal -> worldNormalでなぜinverse->TransposeがいるかはPDFに記載
	worldNormal := invTrans.Transpose().MulByTuple(localNormal)
	worldNormal[3] = 0

	return worldNormal.Normalize(), nil

	// local_point, err := base.WorldToObject(worldPoint)
	// if err != nil {
	// 	return calc.Tuple4{}, err
	// }
	// local_normal := calcLocalNormal(local_point)
	// return base.NormalToWorld(local_normal)
}

func (base *BaseShape) NormalToWorld(normal_vec calc.Tuple4) (calc.Tuple4, error) {
	invTrans, err := base.GetTransform().Inverse()
	if err != nil {
		return calc.Tuple4{}, err
	}

	worldNormal := invTrans.Transpose().MulByTuple(normal_vec)
	worldNormal[3] = 0
	worldNormal = worldNormal.Normalize()

	if base.GetParent() != nil {
		worldNormal, err = base.GetParent().NormalToWorld(worldNormal)
		if err != nil {
			return calc.Tuple4{}, err
		}
	}

	return worldNormal, nil

}

type CalcLocalIntersect func(localRay Ray) (Intersections, error)

func (base *BaseShape) ShapeIntersect(r Ray, localIntersect CalcLocalIntersect) (Intersections, error) {
	invTrans, err := base.GetTransform().Inverse()
	if err != nil {
		return Intersections{}, err
	}

	r = r.Transform(invTrans)

	return localIntersect(r)
}

//子から親に向かってWorldToObjectを計算していく
func (base *BaseShape) WorldToObject(p calc.Tuple4) (calc.Tuple4, error) {

	point := p

	if base.GetParent() != nil {
		inversedPoint, err := base.GetParent().WorldToObject(point)
		if err != nil {
			return calc.Tuple4{}, err
		}

		point = inversedPoint
	}

	transInv, err := base.GetTransform().Inverse()

	if err != nil {
		return calc.Tuple4{}, err
	}

	return transInv.MulByTuple(point), nil
}

type Intersection struct {
	Time   float64
	Object Shape
}

type Intersections struct {
	Intersections []*Intersection
	Count         int
}

func AggregateIntersection(sections ...*Intersection) Intersections {

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Time < sections[j].Time
	})

	return Intersections{
		Count:         len(sections),
		Intersections: sections,
	}

}

//hitはintersectionsの中で最初の正のtimeをもつintersection
func GenerateHit(intersections Intersections) *Intersection {
	for _, intersection := range intersections.Intersections {
		if intersection.Time >= 0 {
			return intersection
		}
	}

	//正のtimeのintersectionがなかったらnil
	return nil
}

// type PreComps struct {
// 	Time        float64
// 	Object      Shape
// 	RayPoint    calc.Tuple4
// 	OverPoint   calc.Tuple4
// 	UnderPoint  calc.Tuple4
// 	EyeVec      calc.Tuple4
// 	NormalVec   calc.Tuple4
// 	ReflectVec  calc.Tuple4
// 	IsRayInside bool
// 	N1          float64
// 	N2          float64
// }

// func FindN1AndN2(xs Intersections, target Intersection) (N1, N2 float64) {

// 	if xs.Count == 0 {
// 		return
// 	}

// 	var container []*Intersection

// 	getN := func(section *Intersection) float64 {

// 		var N float64
// 		//TimeもObjectも完全一致
// 		if section.Time == target.Time && section.Object == target.Object {
// 			if len(container) == 0 {
// 				N = 1.0
// 			} else {
// 				N = container[len(container)-1].Object.GetMaterial().RefractiveIndex
// 			}
// 		}

// 		return N
// 	}

// 	getContainerIndex := func(targetSection *Intersection) int {
// 		for i, section := range container {
// 			//Objectのみ一致
// 			if section.Object == targetSection.Object {
// 				return i
// 			}
// 		}

// 		return -1
// 	}

// 	for _, section := range xs.Intersections {
// 		N1 = getN(section)

// 		index := getContainerIndex(section)
// 		if index == -1 {
// 			container = append(container, section)
// 		} else {
// 			container = append(container[:index], container[index+1:]...)
// 		}
// 		N2 = getN(section)

// 		//targetが終わったらbreak
// 		if N1 != 0 || N2 != 0 {
// 			break
// 		}
// 	}

// 	return
// }

// //normalの計算がおかしいっぽい？
// func PrepareComputations(intersection Intersection, ray Ray, xs Intersections) (PreComps, error) {

// 	t := intersection.Time
// 	obj := intersection.Object
// 	ray_point := ray.Position(t)
// 	eye_vec := calc.NegTuple(ray.Direction)
// 	normal_vec, err := obj.NormalAt(ray_point)

// 	if err != nil {
// 		return PreComps{}, err
// 	}

// 	var IsRayInside bool

// 	//rayのOriginがObjectのInsideにあるとき
// 	if calc.DotTuple(normal_vec, eye_vec) < 0 {
// 		IsRayInside = true
// 		normal_vec = calc.NegTuple(normal_vec)
// 	}

// 	reflect_vec := calc.Reflect(ray.Direction, normal_vec)

// 	//shadow用にOverPointを作る,normal方向に微小に↓にずらしたものがoverpoint
// 	over_point := calc.AddTuple(ray_point, calc.MulTupleByScalar(util.EPSILON, normal_vec))
// 	under_point := calc.SubTuple(ray_point, calc.MulTupleByScalar(util.EPSILON, normal_vec))

// 	n1, n2 := FindN1AndN2(xs, intersection)

// 	return PreComps{
// 		Time:        t,
// 		Object:      obj,
// 		RayPoint:    ray_point,
// 		OverPoint:   over_point,
// 		UnderPoint:  under_point,
// 		EyeVec:      eye_vec,
// 		NormalVec:   normal_vec,
// 		ReflectVec:  reflect_vec,
// 		IsRayInside: IsRayInside,
// 		N1:          n1,
// 		N2:          n2,
// 	}, nil
// }
