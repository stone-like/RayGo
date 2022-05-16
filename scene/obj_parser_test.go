package scene

import (
	"rayGo/calc"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ignore_Unrecognized_Lines(t *testing.T) {
	parser, err := ParseObj("test/test1.txt")
	require.Nil(t, err)
	require.Equal(t, 0, len(parser.Vertices))
}

func Test_Parse_Verice(t *testing.T) {
	parser, err := ParseObj("test/test2.txt")
	require.Nil(t, err)
	require.Equal(t, 4, len(parser.Vertices))

	for i, point := range []calc.Tuple4{
		calc.NewPoint(-1, 1, 0),
		calc.NewPoint(-1, 0.5, -0),
		calc.NewPoint(1, 0, 0),
		calc.NewPoint(1, 1, 0),
	} {
		require.True(t, calc.TupleCompare(point, parser.Vertices[i]))
	}
}

func Test_Parse_Invalid_Verice_Columns_Error(t *testing.T) {
	_, err := ParseObj("test/test3.txt")

	require.Equal(t, "invalid vertice columns", err.Error())
}

func Test_Parse_Not_Float_Verice_Error(t *testing.T) {
	_, err := ParseObj("test/test4.txt")

	require.Equal(t, "strconv.ParseFloat: parsing \"aaaaa\": invalid syntax", err.Error())
}

func Test_Parse_Triangle_Face(t *testing.T) {
	parser, err := ParseObj("test/test5.txt")
	require.Nil(t, err)

	require.Equal(t, 1, len(parser.ParserGroups))
	parserGroup := parser.ParserGroups[0]

	require.Equal(t, 2, len(parserGroup.GetChildren()))

	t1, ok := parserGroup.GetChild(0).(Triangle)
	require.True(t, ok)
	t2, ok := parserGroup.GetChild(1).(Triangle)
	require.True(t, ok)

	require.Equal(t, parser.Vertices[0], t1.P1)
	require.Equal(t, parser.Vertices[1], t1.P2)
	require.Equal(t, parser.Vertices[2], t1.P3)

	require.Equal(t, parser.Vertices[0], t2.P1)
	require.Equal(t, parser.Vertices[2], t2.P2)
	require.Equal(t, parser.Vertices[3], t2.P3)

}

func Test_Parse_Invalid_Face_Error(t *testing.T) {
	_, err := ParseObj("test/test6.txt")
	require.Equal(t, "invalid face num,please make sure It is valid vertice", err.Error())
}

func Test_Parse_Triangle_Polygons(t *testing.T) {
	parser, err := ParseObj("test/test7.txt")
	require.Nil(t, err)

	require.Equal(t, 1, len(parser.ParserGroups))
	parserGroup := parser.ParserGroups[0]

	require.Equal(t, 3, len(parserGroup.GetChildren()))

	t1, ok := parserGroup.GetChild(0).(Triangle)
	require.True(t, ok)
	t2, ok := parserGroup.GetChild(1).(Triangle)
	require.True(t, ok)
	t3, ok := parserGroup.GetChild(2).(Triangle)
	require.True(t, ok)

	require.Equal(t, parser.Vertices[0], t1.P1)
	require.Equal(t, parser.Vertices[1], t1.P2)
	require.Equal(t, parser.Vertices[2], t1.P3)

	require.Equal(t, parser.Vertices[0], t2.P1)
	require.Equal(t, parser.Vertices[2], t2.P2)
	require.Equal(t, parser.Vertices[3], t2.P3)

	require.Equal(t, parser.Vertices[0], t3.P1)
	require.Equal(t, parser.Vertices[3], t3.P2)
	require.Equal(t, parser.Vertices[4], t3.P3)

}

func Test_Parse_Triangle_Groups(t *testing.T) {
	parser, err := ParseObj("test/test8.txt")
	require.Nil(t, err)

	require.Equal(t, 2, len(parser.ParserGroups))
	group1 := parser.ParserGroups[0]
	group2 := parser.ParserGroups[1]

	require.Equal(t, 1, len(group1.GetChildren()))
	require.Equal(t, "FirstGroup", group1.Name)
	require.Equal(t, 1, len(group2.GetChildren()))
	require.Equal(t, "SecondGroup", group2.Name)

	t1, ok := group1.GetChild(0).(Triangle)
	require.True(t, ok)

	t2, ok := group2.GetChild(0).(Triangle)
	require.True(t, ok)

	require.Equal(t, parser.Vertices[0], t1.P1)
	require.Equal(t, parser.Vertices[1], t1.P2)
	require.Equal(t, parser.Vertices[2], t1.P3)

	require.Equal(t, parser.Vertices[0], t2.P1)
	require.Equal(t, parser.Vertices[2], t2.P2)
	require.Equal(t, parser.Vertices[3], t2.P3)

}

func Test_Parse_Triangle_Groups_Error(t *testing.T) {
	_, err := ParseObj("test/test9.txt")
	require.Equal(t, "invalid group columns", err.Error())
}

func Test_Parse_Triangle_DefaultGroup(t *testing.T) {
	parser, err := ParseObj("test/test5.txt")
	require.Nil(t, err)

	require.Equal(t, 1, len(parser.ParserGroups))
	parserGroup := parser.ParserGroups[0]
	require.Equal(t, "defaultGroup", parserGroup.Name)
}

func Test_Convert_File_To_Group(t *testing.T) {
	parser, err := ParseObj("test/test8.txt")
	require.Nil(t, err)
	g := parser.ToGroup()

	require.Equal(t, 2, len(g.Children))

	t1, ok := g.Children[0].(*Group).Children[0].(Triangle)
	require.True(t, ok)

	t2, ok := g.Children[1].(*Group).Children[0].(Triangle)
	require.True(t, ok)

	require.Equal(t, parser.Vertices[0], t1.P1)
	require.Equal(t, parser.Vertices[1], t1.P2)
	require.Equal(t, parser.Vertices[2], t1.P3)

	require.Equal(t, parser.Vertices[0], t2.P1)
	require.Equal(t, parser.Vertices[2], t2.P2)
	require.Equal(t, parser.Vertices[3], t2.P3)
}

func Test_Parser_Vertex_Normal(t *testing.T) {
	parser, err := ParseObj("test/vertexNormal.txt")
	require.Nil(t, err)
	require.Equal(t, calc.NewVector(0, 0, 1), parser.Normals[0])
	require.Equal(t, calc.NewVector(0.707, 0, -0.707), parser.Normals[1])
	require.Equal(t, calc.NewVector(1, 2, 3), parser.Normals[2])

}

func Test_Parser_Vertex_Normal_Error(t *testing.T) {
	_, err := ParseObj("test/vertexNormalError.txt")
	require.Equal(t, "invalid vn columns", err.Error())
}

func Test_Parser_Face_With_Normal(t *testing.T) {
	parser, err := ParseObj("test/faceWithNormal.txt")
	require.Nil(t, err)

	g := parser.ParserGroups[0]

	children := g.GetChildren()
	require.Equal(t, 2, len(children))

	t1, ok := children[0].(SmoothTriangle)
	require.True(t, ok)

	t2, ok := children[1].(SmoothTriangle)
	require.True(t, ok)

	require.Equal(t, parser.Vertices[0], t1.P1)
	require.Equal(t, parser.Vertices[1], t1.P2)
	require.Equal(t, parser.Vertices[2], t1.P3)

	require.Equal(t, parser.Normals[2], t1.N1)
	require.Equal(t, parser.Normals[0], t1.N2)
	require.Equal(t, parser.Normals[1], t1.N3)

	require.Equal(t, t1, t2)

}
