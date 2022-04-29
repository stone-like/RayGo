package scene

import "rayGo/calc"

type Pattern interface {
	PatternAtShape(point calc.Tuple4, shape Shape) (Color, error)
	PatternAt(point calc.Tuple4) Color
	GetMaterial() *Material
	SetMaterial(m *Material)
	GetTransform() calc.Mat4x4
	SetTransform(mat calc.Mat4x4)
}

type BasePattern struct {
	Transform calc.Mat4x4
	Material  *Material
}

func NewBasePattern() *BasePattern {
	return &BasePattern{
		Transform: calc.Ident4x4,
		Material:  DefaultMaterial(),
	}
}

func (b *BasePattern) GetTransform() calc.Mat4x4 {
	return b.Transform
}

func (b *BasePattern) SetTransform(mat calc.Mat4x4) {
	b.Transform = mat
}

func (b *BasePattern) GetMaterial() *Material {
	return b.Material
}

func (b *BasePattern) SetMaterial(m *Material) {
	b.Material = m
}

type PatternAt func(point calc.Tuple4) Color

func (b *BasePattern) PatternAtShapeOnBase(world_point calc.Tuple4, shape Shape, fn PatternAt) (Color, error) {
	shapeTransInv, err := shape.GetTransform().Inverse()

	if err != nil {
		return Color{}, err
	}

	object_point := shapeTransInv.MulByTuple(world_point)

	patternTransInv, err := b.GetTransform().Inverse()
	if err != nil {
		return Color{}, err
	}

	pattern_point := patternTransInv.MulByTuple(object_point)

	return fn(pattern_point), nil
}
