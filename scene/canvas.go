package scene

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Canvas struct {
	Width  int
	Height int
	Pixels [][]Color
}

func NewCanvas(width, height int) *Canvas {

	pixels := make([][]Color, height)

	for i := 0; i < height; i++ {
		pixels[i] = make([]Color, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pixels[i][j] = NewColor(0, 0, 0)
		}
	}

	return &Canvas{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (c *Canvas) WritePixel(x, y int, color Color) {
	c.Pixels[y][x] = color
}

var maxColorValue int = 255
var minColorValue int = 0

func getColorValue(c float64) int {

	target := int(math.Ceil(float64(maxColorValue) * c))
	if target <= minColorValue {
		return minColorValue
	}

	if maxColorValue <= target {
		return maxColorValue
	}

	return target
}

func (c *Canvas) createWidthString(i int) string {
	var buf strings.Builder

	//red ,green,blueとプラスしていくときにプラスした時点で70charを超えたら改行
	//スペースなしで70越えなければOK
	written := 0
	threshold := 69

	for j := 0; j < c.Width; j++ {
		target := c.Pixels[i][j]
		for i := 0; i <= 2; i++ {
			value := getColorValue(target.GetByIndex(i))

			if written+3 > threshold {
				//70に収まらないとき
				buf.WriteString("\n")
				buf.WriteString(strconv.Itoa(value))
				buf.WriteString(" ")
				written = 4
				continue
			}

			//通常
			buf.WriteString(strconv.Itoa(value))
			if written+6 < threshold {
				buf.WriteString(" ")
			}

			written += 4
		}

	}

	return buf.String()
}

func (c *Canvas) ToPPM() string {
	var buf strings.Builder

	buf.WriteString("P3\n")
	buf.WriteString(fmt.Sprintf("%d %d\n", c.Width, c.Height))
	buf.WriteString(fmt.Sprintf("%d\n", maxColorValue))

	for i := 0; i < c.Height; i++ {
		widthString := strings.TrimSuffix(c.createWidthString(i), " ")
		buf.WriteString(widthString)
		buf.WriteString("\n")
	}

	return buf.String()

}
