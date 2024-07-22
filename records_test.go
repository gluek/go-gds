package gds

import (
	"bufio"
	"bytes"
	"testing"
)

var (
	ENDEL        = Record{Size: 4, Datatype: "ENDEL", Data: []byte{}}
	ENDSTR       = Record{Size: 4, Datatype: "ENDSTR", Data: []byte{}}
	ENDLIB       = Record{Size: 4, Datatype: "ENDLIB", Data: []byte{}}
	TestElements = []Element{
		Boundary{
			ElFlags:  1,
			Plex:     2,
			Layer:    2,
			Datatype: 1,
			XY:       []int32{0, 1, 2, 3},
		},
		Path{
			ElFlags:  1,
			Plex:     2,
			Layer:    3,
			Datatype: 4,
			Pathtype: 5,
			Width:    6,
			XY:       []int32{7, 8, 9, 10},
		},
		SRef{
			ElFlags: 1,
			Plex:    2,
			Sname:   "Test",
			Strans:  3,
			Mag:     4.0,
			Angle:   5.0,
			XY:      []int32{6, 7, 8, 9},
		},
	}
)

func mockFilehandler(data []byte) *bufio.Reader {
	return bufio.NewReader(bytes.NewReader(data))
}

func TestRecords(t *testing.T) {
	recordsTest := Record{
		Size:     4 + 2,
		Datatype: "HEADER",
		Data:     []byte{byte(0x00), byte(0x01)},
	}
	recordReader := mockFilehandler(recordsToBytes([]Record{recordsTest}))
	recordNew, err := decodeRecord(recordReader)
	if err != nil {
		t.Fatalf("could not decode record: %v", err)
	}
	if recordsTest.String() != recordNew.String() {
		t.Fatalf("%v not equal to %v", recordsTest, recordNew)
	}
}

func TestText(t *testing.T) {
	textTest := Text{
		ElFlags:      0,
		Plex:         0,
		Layer:        1,
		Texttype:     0,
		Presentation: 0,
		Strans:       0,
		Mag:          1.0,
		Angle:        0.0,
		XY:           []int32{0, 0, 1, 1},
		StringBody:   "Test",
	}
	recordsText, err := fieldsToRecords(textTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsText = append(recordsText, ENDEL)
	textNew, err := decodeText(mockFilehandler(recordsToBytes(recordsText)))
	if err != nil {
		t.Fatalf("error decoding aref: %v", err)
	}
	if textTest.String() != textNew.String() {
		t.Fatalf("%v not equal to %v", textNew, textNew)
	}
}

func TestBoundary(t *testing.T) {
	boundaryTest := Boundary{
		ElFlags:  1,
		Plex:     2,
		Layer:    2,
		Datatype: 1,
		XY:       []int32{0, 1, 2, 3},
	}
	recordsBoundary, err := fieldsToRecords(boundaryTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsBoundary = append(recordsBoundary, ENDEL)
	boundaryNew, err := decodeBoundary(mockFilehandler(recordsToBytes(recordsBoundary)))
	if err != nil {
		t.Fatalf("could not decode boundary: %v", err)
	}
	if boundaryTest.String() != boundaryNew.String() {
		t.Fatalf("%v not equal to %v", boundaryTest, boundaryNew)
	}
}

func TestPath(t *testing.T) {
	pathTest := Path{
		ElFlags:  1,
		Plex:     2,
		Layer:    3,
		Datatype: 4,
		Pathtype: 5,
		Width:    6,
		XY:       []int32{7, 8, 9, 10},
	}
	recordsPath, err := fieldsToRecords(pathTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsPath = append(recordsPath, ENDEL)
	pathNew, err := decodePath(mockFilehandler(recordsToBytes(recordsPath)))
	if err != nil {
		t.Fatalf("could not decode path: %v", err)
	}
	if pathTest.String() != pathNew.String() {
		t.Fatalf("%v not equal to %v", pathTest, pathNew)
	}
}

func TestSref(t *testing.T) {
	srefTest := SRef{
		ElFlags: 1,
		Plex:    2,
		Sname:   "Test",
		Strans:  3,
		Mag:     4.0,
		Angle:   5.0,
		XY:      []int32{6, 7, 8, 9},
	}
	recordsSref, err := fieldsToRecords(srefTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsSref = append(recordsSref, ENDEL)
	srefNew, err := decodeSREF(mockFilehandler(recordsToBytes(recordsSref)))
	if err != nil {
		t.Fatalf("could not decode sref: %v", err)
	}
	if srefTest.String() != srefNew.String() {
		t.Fatalf("%v not equal to %v", srefTest, srefNew)
	}
}

func TestAref(t *testing.T) {
	arefTest := ARef{
		ElFlags: 0,
		Plex:    0,
		Sname:   "Test",
		Strans:  0,
		Mag:     1.0,
		Angle:   0.0,
		Colrow:  []int16{0, 0},
		XY:      []int32{0, 0, 1, 1},
	}
	recordsAref, err := fieldsToRecords(arefTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsAref = append(recordsAref, ENDEL)
	arefNew, err := decodeAREF(mockFilehandler(recordsToBytes(recordsAref)))
	if err != nil {
		t.Fatalf("error decoding aref: %v", err)
	}
	if arefTest.String() != arefNew.String() {
		t.Fatalf("%v not equal to %v", arefTest, arefNew)
	}
}

func TestNode(t *testing.T) {
	nodeTest := Node{
		ElFlags:  1,
		Plex:     2,
		Layer:    3,
		Nodetype: 4,
		XY:       []int32{5, 6, 7, 8},
	}
	recordsNode, err := fieldsToRecords(nodeTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsNode = append(recordsNode, ENDEL)
	nodeNew, err := decodeNode(mockFilehandler(recordsToBytes(recordsNode)))
	if err != nil {
		t.Fatalf("error decoding node: %v", err)
	}
	if nodeTest.String() != nodeNew.String() {
		t.Fatalf("%v not equal to %v", nodeTest, nodeNew)
	}
}

func TextBox(t *testing.T) {
	boxTest := Box{
		ElFlags: 1,
		Plex:    2,
		Layer:   3,
		Boxtype: 4,
		XY:      []int32{5, 6, 7, 8},
	}
	recordsBox, err := fieldsToRecords(boxTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsBox = append(recordsBox, ENDEL)
	boxNew, err := decodeBox(mockFilehandler(recordsToBytes(recordsBox)))
	if err != nil {
		t.Fatalf("error decoding box: %v", err)
	}
	if boxTest.String() != boxNew.String() {
		t.Fatalf("%v not equal to %v", boxTest, boxNew)
	}
}

func TestStructure(t *testing.T) {
	structureTest := Structure{
		BgnStr:   []int16{1, 2},
		StrName:  "TestStructure",
		Elements: TestElements,
	}
	recordsStructure, err := fieldsToRecords(structureTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsStructure = append(recordsStructure, ENDSTR)
	structureNew, err := decodeStructure(mockFilehandler(recordsToBytes(recordsStructure)),
		&Record{Size: 6, Datatype: "BGNSTR", Data: []byte{byte(0x00), byte(0x01), byte(0x00), byte(0x02)}})
	if err != nil {
		t.Fatalf("error decoding structure: %v", err)
	}
	if structureTest.String() != structureNew.String() {
		t.Fatalf("%v not equal to %v", structureTest, structureNew)
	}
}

func TestLibrary(t *testing.T) {
	libraryTest := Library{
		Header:  1,
		BgnLib:  []int16{2, 3},
		LibName: "TestLibrary",
		Units:   []float64{4.0, 5.0},
		Structures: []Structure{
			{
				BgnStr:   []int16{1, 2},
				StrName:  "TestStructure",
				Elements: TestElements,
			},
		},
	}
	recordsLibrary, err := fieldsToRecords(libraryTest)
	if err != nil {
		t.Fatalf("error fields to records: %v", err)
	}
	recordsLibrary = append(recordsLibrary, ENDLIB)
	libraryNew, err := decodeLibrary(mockFilehandler(recordsToBytes(recordsLibrary)))
	if err != nil {
		t.Fatalf("could not decode library: %v", err)
	}
	if libraryTest.String() != libraryNew.String() {
		t.Fatalf("%v not equal to %v", libraryTest, libraryNew)
	}
}
