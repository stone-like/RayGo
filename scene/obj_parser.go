package scene

import (
	"bufio"
	"os"
	"rayGo/calc"
	"rayGo/scene/files"
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

func (p *Parser) isValidFace(num int) bool {
	return 0 <= num && num <= len(p.Vertices)-1
}

//入力は1-indexedで受け付けているので
func change0indexed(num float64) int {
	return int(num) - 1
}

func retrieveFace(data []string) ([]int, error) {

	var faces []int
	floatData, err := retrieveComponentFromData(data)
	if err != nil {
		return faces, err
	}

	for i := 0; i < len(floatData); i++ {
		faces = append(faces, change0indexed(floatData[i]))
	}

	return faces, nil
}

//ParseFaceで事前に得られたVeriticesから実際にtriangleを作っていく
func (p *Parser) parseFace(line string) error {
	//空白区切り
	faceComponent := strings.Split(line, " ")

	faces, err := retrieveFace(faceComponent[1:])
	if err != nil {
		return err
	}

	return p.createObject(faces)

}

func (p *Parser) checkFaceIsValid(faces []int) bool {
	for _, vertice := range faces {
		if !p.isValidFace(int(vertice)) {
			return false
		}
	}

	return true
}

func (p *Parser) fanTriangulation(points []calc.Tuple4) []Triangle {
	var triangles []Triangle

	for i := 1; i < len(points)-1; i++ {
		triangles = append(triangles, NewTriangle(points[0], points[i], points[i+1]))
	}

	return triangles
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

func (p *Parser) createObject(faces []int) error {

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

	if !p.checkFaceIsValid(faces) {
		return NewParserError("invalid face num,please make sure It is valid vertice")
	}

	points := make([]calc.Tuple4, len(faces))

	for i, verticeNum := range faces {
		points[i] = p.Vertices[verticeNum]
	}

	for _, tri := range p.fanTriangulation(points) {
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

func (p *Parser) ParseLine(line string) error {

	//改行オンリーの行は飛ばす
	if len(line) == 0 {
		return nil
	}

	startChar := line[0]

	switch startChar {
	case 'v':
		return p.parseVertice(line)
	case 'f':
		return p.parseFace(line)
	case 'g':
		return p.parseGroup(line)
	default:
		return nil

	}

}
