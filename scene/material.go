package scene

type Material struct {
	color     Color
	ambient   float64
	diffuse   float64
	specular  float64
	shininess float64
}

const (
	maxAmbient, maxDiffuse, maxSpecular = 1, 1, 1
	minAmbient, minDiffuse, minSpecular = 0, 0, 0

	maxShininess = 10
	minShininess = 200
)

func NewMaterial(color Color, ambient, diffuse, specular, shininess float64) Material {
	return Material{
		color:     color,
		ambient:   ambient,
		diffuse:   diffuse,
		specular:  specular,
		shininess: shininess,
	}
}

func DefaultMaterial() Material {
	return Material{
		color:     NewColor(1, 1, 1),
		ambient:   0.1,
		diffuse:   0.9,
		specular:  0.9,
		shininess: 200.0,
	}
}
