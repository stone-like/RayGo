package scene

import (
	"bufio"
	"os"
	"rayGo/calc"
	"rayGo/files"
	"strconv"
	"strings"
)

type ParserError struct {
	message string
}

func (p ParserError) Error() string {
	return p.message
}

func NewParserError(message string) ParserError {
	return ParserError{
		message: message,
	}
}

type ParserGroup struct {
	Name  string
	Group *Group
}

func NewParserGroup(name string) *ParserGroup {
	return &ParserGroup{
		Name:  name,
		Group: NewGroup(),
	}
}

func (pg *ParserGroup) GetChildren() []Shape {
	return pg.Group.Children
}

func (pg *ParserGroup) GetChild(num int) Shape {
	return pg.Group.Children[num]
}

func (pg *ParserGroup) AddChildren(ss ...Shape) {
	pg.Group.AddChildren(ss...)
}

type Parser struct {
	ParserGroups []*ParserGroup
	Vertices     []calc.Tuple4
	Normals      []calc.Tuple4
}

func NewParser() *Parser {
	return &Parser{}
}

func getLines(fileName string) ([]string, error) {
	path := files.GetFilePath(fileName)
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil

}

type FaceComponent struct {
	VertexNum       int
	VertexNormalNum int
}

func NewFaceComponent() FaceComponent {
	return FaceComponent{
		-1,
		-1,
	}
}

func retrieveComponentFromData(data []string) ([]float64, error) {

	floatData := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		fp, err := strconv.ParseFloat(data[i], 64)
		if err != nil {
			return []float64{}, err
		}

		floatData[i] = fp
	}

	return floatData, nil

}

// num or
// num//num or
// num/num/num
func parseFaceComponent(data string) (FaceComponent, error) {

	faceComponent := NewFaceComponent()

	if strings.Contains(data, "//") {
		return parseDoubleSlash(data, faceComponent)
	}

	if strings.Contains(data, "/") {
		return parseSingleSlash(data, faceComponent)
	}

	return parseNumOnly(data, faceComponent)
}

func parseNumOnly(data string, faceComponent FaceComponent) (FaceComponent, error) {
	num, err := strconv.Atoi(data)
	if err != nil {
		return faceComponent, err
	}

	faceComponent.VertexNum = change0indexed(num)

	return faceComponent, nil
}

func parseDoubleSlash(data string, faceComponent FaceComponent) (FaceComponent, error) {
	// //区切り
	doubleSlashComponent := strings.Split(data, "//")
	if len(doubleSlashComponent) != 2 {
		return faceComponent, NewParserError("invalid face double slash columns")
	}

	vertexNum, err := strconv.Atoi(doubleSlashComponent[0])
	if err != nil {
		return faceComponent, err
	}
	vertexNormalNum, err := strconv.Atoi(doubleSlashComponent[1])
	if err != nil {
		return faceComponent, err
	}

	faceComponent.VertexNum = change0indexed(vertexNum)
	faceComponent.VertexNormalNum = change0indexed(vertexNormalNum)

	return faceComponent, nil

}

func parseSingleSlash(data string, faceComponent FaceComponent) (FaceComponent, error) {
	// //区切り
	singleSlashComponent := strings.Split(data, "/")
	if len(singleSlashComponent) != 3 {
		return faceComponent, NewParserError("invalid face single slash columns")
	}

	vertexNum, err := strconv.Atoi(singleSlashComponent[0])
	if err != nil {
		return faceComponent, err
	}
	vertexNormalNum, err := strconv.Atoi(singleSlashComponent[2])
	if err != nil {
		return faceComponent, err
	}

	faceComponent.VertexNum = change0indexed(vertexNum)
	faceComponent.VertexNormalNum = change0indexed(vertexNormalNum)

	return faceComponent, nil

}

func retrieveFaceComponentFromData(data []string) ([]FaceComponent, error) {

	faceData := make([]FaceComponent, len(data))

	for i := 0; i < len(data); i++ {

		if data[i] == "" {
			continue
		}
		faceComponent, err := parseFaceComponent(data[i])
		if err != nil {
			return nil, err
		}
		faceData[i] = faceComponent
	}

	return faceData, nil

}

func ParseObj(fileName string) (*Parser, error) {

	parser := NewParser()

	lines, err := getLines(fileName)
	if err != nil {
		return &Parser{}, err
	}

	for _, line := range lines {
		err := parser.ParseLine(line)
		if err != nil {
			return &Parser{}, err
		}
	}

	return parser, nil
}

func createVertice(data []float64) (calc.Tuple4, error) {
	if len(data) != 3 {
		return calc.Tuple4{}, NewParserError("invalid vertice columns")
	}

	return calc.NewPoint(data[0], data[1], data[2]), nil
}

func (p *Parser) parseVertice(line string) error {
	//空白区切り

	verticeComponent := strings.Split(line, " ")

	if len(verticeComponent) != 4 {
		return NewParserError("invalid vertice columns")
	}

	verticeData, err := retrieveComponentFromData(verticeComponent[1:])
	if err != nil {
		return err
	}

	vertice, err := createVertice(verticeData)
	if err != nil {
		return err
	}

	p.Vertices = append(p.Vertices, vertice)

	return nil
}

func (p *Parser) isValidVertex(num int) bool {
	return 0 <= num && num <= len(p.Vertices)-1
}

func (p *Parser) isValidVertexNormal(num int) bool {
	return 0 <= num && num <= len(p.Normals)-1
}

//入力は1-indexedで受け付けているので
func change0indexed(num int) int {
	return num - 1
}

//ParseFaceで事前に得られたVeriticesから実際にtriangleを作っていく
func (p *Parser) parseFace(line string) error {
	//空白区切り
	faceComponent := strings.Split(line, " ")

	faceData, err := retrieveFaceComponentFromData(faceComponent[1:])
	if err != nil {
		return err
	}

	return p.createObject(faceData)

}

// func (p *Parser) checkFaceIsValid(faces []int) bool {
// 	for _, vertice := range faces {
// 		if !p.isValidFace(int(vertice)) {
// 			return false
// 		}
// 	}

// 	return true
// }

func convertToPoint(vertexData VertexData) []calc.Tuple4 {
	points := make([]calc.Tuple4, len(vertexData.data))

	for i, d := range vertexData.data {
		points[i] = d.vertex
	}

	return points
}

func fanTriangulation(vertexData VertexData) []Shape {
	var triangles []Shape

	points := convertToPoint(vertexData)

	for i := 1; i < len(points)-1; i++ {
		triangles = append(triangles, NewTriangle(points[0], points[i], points[i+1]))
	}

	return triangles
}

//

//こっちの場合はfanTriangulationみたいに五角形とかを想定していない...はず
//縦に[]{
// verttex
// vertexNormal
//}と持つのではなく
//横に
//{
// []vertex
// []vertexNormal
//}とした方が取り回し良さそうなのでtest通ったらリファクタ
func createSmoothTriangles(vertexData VertexData) []Shape {
	one, two, three := vertexData.data[0], vertexData.data[1], vertexData.data[2]

	return []Shape{
		NewSmoothTriangle(
			one.vertex, two.vertex, three.vertex,
			one.vertexNormal, two.vertexNormal, three.vertexNormal,
		),
	}
}

func (p *Parser) createTriangles(vertexData VertexData) []Shape {
	if vertexData.isUseVertexNormal {
		return createSmoothTriangles(vertexData)
	}

	return fanTriangulation(vertexData)

}

func (p *Parser) groupAlreadyCreated() bool {
	return len(p.ParserGroups) != 0
}

func (p *Parser) createDefaultGroup() {
	p.ParserGroups = append(p.ParserGroups, NewParserGroup("defaultGroup"))
}

func (p *Parser) latestGroup() *ParserGroup {

	if len(p.ParserGroups) == 0 {
		return nil
	}
	return p.ParserGroups[len(p.ParserGroups)-1]
}

type VertexDatum struct {
	vertex       calc.Tuple4
	vertexNormal calc.Tuple4
}

type VertexData struct {
	data              []VertexDatum
	isUseVertexNormal bool
}

func (p *Parser) convertDataToVertexAndVertexNormal(faceData []FaceComponent) (VertexData, error) {

	setVertex := func(vertexs *VertexDatum, face FaceComponent) error {
		if !p.isValidVertex(face.VertexNum) {
			return NewParserError("invalid face num,please make sure It is valid vertice")
		}
		vertexs.vertex = p.Vertices[face.VertexNum]

		return nil
	}

	setVertexNormal := func(vertexs *VertexDatum, face FaceComponent) (bool, error) {

		if face.VertexNormalNum == -1 {
			return false, nil
		}

		if !p.isValidVertexNormal(face.VertexNum) {
			return true, NewParserError("invalid face num,please make sure It is valid vertex Normal")
		}

		vertexs.vertexNormal = p.Normals[face.VertexNormalNum]

		return true, nil
	}

	data := make([]VertexDatum, len(faceData))

	isUseVertexNormal := true

	for i, eachFace := range faceData {

		vertexs := &VertexDatum{}

		err := setVertex(vertexs, eachFace)
		if err != nil {
			return VertexData{}, err
		}

		useVertexNormal, err := setVertexNormal(vertexs, eachFace)
		if err != nil {
			return VertexData{}, err
		}

		isUseVertexNormal = useVertexNormal

		data[i] = *vertexs

	}

	return VertexData{
		data:              data,
		isUseVertexNormal: isUseVertexNormal,
	}, nil
}

func (p *Parser) createObject(faceData []FaceComponent) error {

	//g FirstGroup
	//f 1 2 3
	//g SecondGroup
	//f 1 2 3
	//みたいにあるとき、fでつくったObjは直前のgに追加するようにしたい
	//なのでp.Groups[:len(p.Groups)-1]に追加すればよい
	//gがないときはp.Groupsが空なのでdefaultGroupを作る
	if !p.groupAlreadyCreated() {
		p.createDefaultGroup()
	}

	vertexData, err := p.convertDataToVertexAndVertexNormal(faceData)

	if err != nil {
		return err
	}

	for _, tri := range p.createTriangles(vertexData) {
		p.latestGroup().AddChildren(tri)
	}

	return nil
}

func (p *Parser) parseGroup(line string) error {
	//空白区切り
	groupComponent := strings.Split(line, " ")

	if len(groupComponent) != 2 {
		return NewParserError("invalid group columns")
	}

	p.ParserGroups = append(p.ParserGroups, NewParserGroup(groupComponent[1]))

	return nil
}

func (p *Parser) ToGroup() *Group {
	retGroup := NewGroup()

	for _, parserGroup := range p.ParserGroups {
		retGroup.AddChildren(parserGroup.Group)
	}

	return retGroup
}

func createVerticeNormal(data []float64) (calc.Tuple4, error) {
	if len(data) != 3 {
		return calc.Tuple4{}, NewParserError("invalid vn columns")
	}

	return calc.NewVector(data[0], data[1], data[2]), nil
}

func (p *Parser) ParseVertexNormal(line string) error {
	//空白区切り
	vnComponent := strings.Split(line, " ")

	if len(vnComponent) != 4 {
		return NewParserError("invalid vn columns")
	}

	vnDara, err := retrieveComponentFromData(vnComponent[1:])
	if err != nil {
		return err
	}

	vn, err := createVerticeNormal(vnDara)
	if err != nil {
		return err
	}

	p.Normals = append(p.Normals, vn)

	return nil
}

//v~の奴を処理
func (p *Parser) ParseVxxx(line string) error {
	nextChar := line[1]

	switch nextChar {
	//vn
	case 'n':
		return p.ParseVertexNormal(line)
	default:
		return p.parseVertice(line)
	}
}

func (p *Parser) ParseLine(line string) error {

	//改行オンリーの行は飛ばす
	if len(line) == 0 {
		return nil
	}

	startChar := line[0]

	switch startChar {
	case 'v':
		return p.ParseVxxx(line)
	case 'f':
		return p.parseFace(line)
	case 'g':
		return p.parseGroup(line)
	default:
		return nil

	}

}
