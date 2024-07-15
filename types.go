package gds

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"reflect"
)

var Datatypes map[string]string = map[string]string{
	"00": "No data",
	"01": "bitarray",
	"02": "2-byte signed integers",
	"03": "4-byte signed integer",
	"04": "4-byte floats",
	"05": "8-byte floats",
	"06": "ASCII string",
}

var RecordTypes map[string]string = map[string]string{
	"0002": "HEADER",       // version number
	"0102": "BGNLIB",       // begin of library, last modification data
	"0206": "LIBNAME",      // name of library
	"0305": "UNITS",        // user and database units, contains two eight-byte real numbers. The first number is the size of a database unit in user units. The second number is the size of a database unit in meters.
	"0400": "ENDLIB",       // end of library
	"0502": "BGNSTR",       // begin of structure + creation and modification time
	"0606": "STRNAME",      // name of structure
	"0700": "ENDSTR",       // end of structure
	"0800": "BOUNDARY",     // begin of boundary element
	"0900": "PATH",         // begin of path element
	"0a00": "SREF",         // begin of structure reference element
	"0b00": "AREF",         // begin of array reference element
	"0c00": "TEXT",         // begin of text element
	"0d02": "LAYER",        // layer number of element
	"0e02": "DATATYPE",     // Datatype number of element
	"0f03": "WIDTH",        // width of element in db units
	"1003": "XY",           // list of xy coordinates in db units
	"1100": "ENDEL",        // end of element
	"1206": "SNAME",        // name of structure reference
	"1302": "COLROW",       // number of colrow[0]=columns and colrow[1]=rows in array reference
	"1500": "NODE",         // begin of node element
	"1602": "TEXTTYPE",     // texttype number
	"1701": "PRESENTATION", // text presentation, font
	"1906": "STRING",       // character string for text element
	"1a01": "STRANS",       // array reference, structure reference and text transform flags
	"1b05": "MAG",          // magnification factor for text and references
	"1c05": "ANGLE",        // rotation angle for text and references
	"1f06": "REFLIBS",      // name of referenced libraries
	"2006": "FONTS",        // name of text fonts definition files
	"2102": "PATHTYPE",     // type of PATH element end ( rounded, square)
	"2202": "GENERATIONS",  // number of deleted structure ?????
	"2306": "ATTRTABLE",    // attribute table, used in combination with element properties
	"2601": "ELFLAGS",      // template data
	"2a02": "NODETYPE",     // node type number for NODE element
	"2b02": "PROPATTR",     // attribute number
	"2c06": "PROPVALUE",    // attribute name
	"2d00": "BOX",          // begin of box element
	"2e02": "BOXTYPE",      // boxtype for box element
	"2f03": "PLEX",         // plex number
	"3202": "TAPENUM",      // tape number
	"3302": "TAPECODE",     // tape code
	"3602": "FORMAT",       // format type
	"3706": "MASK",         // list of layers
	"3800": "ENDMASKS",     // end of MASK
}

const HEADERSIZE = 4

type Record struct {
	Size     uint16
	Datatype string
	Data     []byte
}

type Element interface {
	String() string
	GetData() any
}

type Library struct {
	Header     int16
	BgnLib     []int16
	LibName    string
	Units      []float64
	Structures []*Structure
}

func (l Library) String() string {
	structureInfo := "\n"
	for _, structure := range l.Structures {
		structureInfo += "      " + structure.String() + "\n"
		structureElements := structure.ListElements()
		structureInfo += structureElements
	}
	return fmt.Sprintf(`Library:
   Version: %d
   Name: %s
   Units: %v
   Structures:%s`, l.Header, l.LibName, l.Units, structureInfo)
}

type Structure struct {
	BgnStr   []int16
	StrName  string
	Elements []Element
}

func (s Structure) String() string {
	return s.StrName
}

func (s Structure) ListElements() string {
	result := ""
	for _, v := range s.Elements {
		result += "         " + v.String() + "\n"
	}
	return result
}

type Boundary struct {
	ElFlags  uint16
	Plex     int16
	Layer    int16
	Datatype int16
	XY       []int32
}

func (b Boundary) GetData() any {
	return b.XY
}
func (b Boundary) String() string {
	return fmt.Sprintf("Boundary - ElFlags: %v, Plex: %v, Layer: %v, Datatype: %v, XY: %v", b.ElFlags, b.Plex, b.Layer, b.Datatype, b.XY)
}

type Path struct {
	ElFlags  uint16
	Plex     int16
	Layer    int16
	Datatype int16
	Pathtype int16
	Width    int32
	XY       []int32
}

func (p Path) GetData() any {
	return p.XY
}
func (p Path) String() string {
	return fmt.Sprintf("Path - ElFlags: %v, Plex: %v, Layer: %v, Datatype: %v, Pathtype: %v, Width: %v, XY: %v",
		p.ElFlags, p.Plex, p.Layer, p.Datatype, p.Pathtype, p.Width, p.XY)
}

type Text struct {
	ElFlags uint16
	Plex    int16
	Layer   int16
	Textbody
}

func (t Text) GetData() any {
	return t.Textbody.String()
}
func (t Text) String() string {
	return fmt.Sprintf("Text - ElFlags: %v, Plex: %v, Layer: %v, Textbody: %s, XY: %v", t.ElFlags, t.Plex, t.Layer, t.Textbody, t.Textbody.XY)
}

type Node struct {
	ElFlags  uint16
	Plex     int16
	Layer    int16
	Nodetype int16
	XY       []int32
}

func (n Node) GetData() any {
	return n.XY
}
func (n Node) String() string {
	return fmt.Sprintf("Node - ElFlags: %v, Plex: %v, Layer: %v, Nodetype: %v, XY: %v", n.ElFlags, n.Plex, n.Layer, n.Nodetype, n.XY)
}

type Box struct {
	ElFlags uint16
	Plex    int16
	Layer   int16
	Boxtype int16
	XY      []int32
}

func (b Box) GetData() any {
	return b.XY
}
func (b Box) String() string {
	return fmt.Sprintf("Box - ElFlags: %v, Plex: %v, Layer: %v, Boxtype: %v, XY: %v", b.ElFlags, b.Plex, b.Layer, b.Boxtype, b.XY)
}

type SRef struct {
	ElFlags uint16
	Plex    int16
	Sname   string
	Strans  uint16
	Mag     float64
	Angle   float64
	XY      []int32
}

func (s SRef) GetData() any {
	return s.XY
}
func (s SRef) String() string {
	return fmt.Sprintf("SRef - ElFlags: %v, Plex: %v, Sname: %v, Strans: %v, Mag: %v, Angle: %v, XY: %v",
		s.ElFlags, s.Plex, s.Sname, s.Strans, s.Mag, s.Angle, s.XY)
}

type ARef struct {
	ElFlags uint16
	Plex    int16
	Sname   string
	Strans  uint16
	Mag     float64
	Angle   float64
	Colrow  []int16
	XY      []int32
}

func (a ARef) GetData() any {
	return a.XY
}
func (a ARef) String() string {
	return fmt.Sprintf("ARef - ElFlags: %v, Plex: %v, Sname: %v, Strans: %v, Mag: %v, Angle: %v, Colrow: %v, XY: %v",
		a.ElFlags, a.Plex, a.Sname, a.Strans, a.Mag, a.Angle, a.Colrow, a.XY)
}

type Textbody struct {
	Texttype     int16
	Presentation uint16
	Strans       uint16
	Mag          float64
	Angle        float64
	XY           []int32
	StringBody   string
}

func (t Textbody) String() string {
	return t.StringBody
}

func binaryToFloat(f string) float64 {
	result := 0.0
	for i := range len(f) {
		if f[i] == '1' {
			result += math.Pow(2, float64(-1*(i+1)))
		}
	}
	return result
}

func getRealSlice(data Record) []float64 {
	initSlice := make([]uint64, int((data.Size-HEADERSIZE)/8))
	finalSlice := make([]float64, len(initSlice))

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &initSlice)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	for i, number := range initSlice {
		sign := float64(1)
		if fmt.Sprintf("%064b", number)[0] == '1' {
			sign = float64(-1)
		}
		exponent := int8((number >> 56))
		mantisse := binaryToFloat(fmt.Sprintf("%064b", number<<8))
		value := sign * mantisse * math.Pow(16, math.Abs(float64(exponent))-64)
		finalSlice[i] = value
	}
	return finalSlice
}

func getDataSlice[T any](data Record) []T {
	var typeInit T
	typeSize := reflect.TypeOf(typeInit).Size()
	result := make([]T, int((data.Size-HEADERSIZE)/uint16(typeSize)))
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	return result
}

func getRealPoint(data Record) float64 {
	var number uint64

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &number)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	sign := float64(1)
	if fmt.Sprintf("%064b", number)[0] == '1' {
		sign = float64(-1)
	}
	exponent := int8((number >> 56))
	mantisse := binaryToFloat(fmt.Sprintf("%064b", number<<8))
	floatValue := sign * mantisse * math.Pow(16, math.Abs(float64(exponent))-64)
	return floatValue
}

func getDataPoint[T any](data Record) T {
	var result T
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	return result
}

func getDataString(data Record) string {
	return string(data.Data)
}

func (r Record) GetData() any {
	switch r.Datatype {
	case "HEADER":
		return getDataPoint[int16](r)
	case "BGNLIB":
		return getDataSlice[int16](r)
	case "LIBNAME":
		return getDataString(r)
	case "UNITS":
		return getRealSlice(r)
	case "ENDLIB":
		return "No data"
	case "BGNSTR":
		return getDataSlice[int16](r)
	case "STRNAME":
		return getDataString(r)
	case "ENDSTR":
		return "No data"
	case "BOUNDARY":
		return "No data"
	case "PATH":
		return "No data"
	case "SREF":
		return "No data"
	case "AREF":
		return "No data"
	case "TEXT":
		return "No data"
	case "LAYER":
		return getDataPoint[int16](r)
	case "DATATYPE":
		return getDataPoint[int16](r)
	case "WIDTH":
		return getDataPoint[int32](r)
	case "XY":
		return getDataSlice[int32](r)
	case "ENDEL":
		return "No data"
	case "SNAME":
		return getDataString(r)
	case "COLROW":
		return getDataSlice[int16](r)
	case "NODE":
		return "No data"
	case "TEXTTYPE":
		return getDataPoint[int16](r)
	case "PRESENTATION":
		return getDataPoint[uint16](r)
	case "STRING":
		return getDataString(r)
	case "STRANS":
		return getDataPoint[uint16](r)
	case "MAG":
		return getRealPoint(r)
	case "ANGLE":
		return getRealPoint(r)
	case "REFLIBS":
		return getDataString(r)
	case "FONTS":
		return getDataString(r)
	case "PATHTYPE":
		return getDataPoint[int16](r)
	case "GENERATIONS":
		return getDataPoint[int16](r)
	case "ATTRTABLE":
		return getDataString(r)
	case "ELFLAGS":
		return getDataPoint[uint16](r)
	case "NODETYPE":
		return getDataPoint[int16](r)
	case "PROPATTR":
		return getDataPoint[int16](r)
	case "PROPVALUE":
		return getDataString(r)
	case "BOX":
		return "No data"
	case "BOXTYPE":
		return getDataPoint[int16](r)
	case "PLEX":
		return getDataPoint[int32](r)
	case "TAPENUM":
		return getDataPoint[int16](r)
	case "TAPECODE":
		return getDataPoint[int16](r)
	case "FORMAT":
		return getDataPoint[int16](r)
	case "MASK":
		return getDataString(r)
	case "ENDMASKS":
		return "No data"
	default:
		panic("unexpected datatype")
	}
}

func (r Record) String() string {
	return fmt.Sprintf("%s, Bytes: %d Data: %v", r.Datatype, r.Size, r.GetData())
}
