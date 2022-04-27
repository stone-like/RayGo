package scene

type Material struct {
	Color     Color
	Ambient   float64
	Diffuse   float64
	Specular  float64
	Shininess float64
}

const (
	maxAmbient, maxDiffuse, maxSpecular = 1, 1, 1
	minAmbient, minDiffuse, minSpecular = 0, 0, 0

	maxShininess = 10
	minShininess = 200
)

func NewMaterial(color Color, ambient, diffuse, specular, shininess float64) Material {
	return Material{
		Color:     color,
		Ambient:   ambient,
		Diffuse:   diffuse,
		Specular:  specular,
		Shininess: shininess,
	}
}

func DefaultMaterial() Material {
	return Material{
		Color:     NewColor(1, 1, 1),
		Ambient:   0.1,
		Diffuse:   0.9,
		Specular:  0.9,
		Shininess: 200.0,
	}
}