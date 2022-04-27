package scene

func CreateIntersection(t float64, object Shape) Intersection {
	return Intersection{
		Time:   t,
		Object: object,
	}
}
