package gds

import (
	"fmt"
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
	"1906": "STRINGBODY",   // character string for text element
	"1a01": "STRANS",       // array reference, structure reference and text transform flags
	"1b05": "MAG",          // magnification factor for text and references
	"1c05": "ANGLE",        // rotation angle for text and references
	"1f06": "REFLIBS",      // name of referenced libraries
	"2006": "FONTS",        // name of text fonts definition files
	"2102": "PATHTYPE",     // type of PATH element end ( rounded, square)
	"2202": "GENERATIONS",  // number of deleted structure ?????
	"2306": "ATTRTABLE",    // attribute table, used in combination with element properties
	"2601": "ELFLAGS",      // template data
	"2703": "ELKEY",        // from KLayout Source Code ???
	"2802": "LINKTYPE",     // unreleased Feature
	"2903": "LINKKEYS",     // unreleased Feature
	"2a02": "NODETYPE",     // node type number for NODE element
	"2b02": "PROPATTR",     // attribute number
	"2c06": "PROPVALUE",    // attribute name
	"2d00": "BOX",          // begin of box element
	"2e02": "BOXTYPE",      // boxtype for box element
	"2f03": "PLEX",         // plex number
	"3003": "BGNEXTN",      // path type 4 extension start
	"3103": "ENDEXTN",      // path type 4 extension end
	"3202": "TAPENUM",      // tape number
	"3302": "TAPECODE",     // tape code
	"3503": "RESERVED",     // type was used for NUMTYPES but was not required
	"3602": "FORMAT",       // format type
	"3706": "MASK",         // list of layers
	"3800": "ENDMASKS",     // end of MASK
	"3902": "LIBDIRSIZE",   // contains the number of pages in the Library directory
	"3A06": "SRFNAME",      // contains the name of the Sticks Rules File, if one is bound to the library
	"3B02": "LIBSECUR",     // contains an array of Access Control List (ACL) data
}

var RecordTypesBytes map[string][]byte = map[string][]byte{
	"HEADER":       {0x00, 0x02}, // version number
	"BGNLIB":       {0x01, 0x02}, // begin of library, last modification data
	"LIBNAME":      {0x02, 0x06}, // name of library
	"UNITS":        {0x03, 0x05}, // user and database units, contains two eight-byte real numbers. The first number is the size of a database unit in user units. The second number is the size of a database unit in meters.
	"ENDLIB":       {0x04, 0x00}, // end of library
	"BGNSTR":       {0x05, 0x02}, // begin of structure + creation and modification time
	"STRNAME":      {0x06, 0x06}, // name of structure
	"ENDSTR":       {0x07, 0x00}, // end of structure
	"BOUNDARY":     {0x08, 0x00}, // begin of boundary element
	"PATH":         {0x09, 0x00}, // begin of path element
	"SREF":         {0x0a, 0x00}, // begin of structure reference element
	"AREF":         {0x0b, 0x00}, // begin of array reference element
	"TEXT":         {0x0c, 0x00}, // begin of text element
	"LAYER":        {0x0d, 0x02}, // layer number of element
	"DATATYPE":     {0x0e, 0x02}, // Datatype number of element
	"WIDTH":        {0x0f, 0x03}, // width of element in db units
	"XY":           {0x10, 0x03}, // list of xy coordinates in db units
	"ENDEL":        {0x11, 0x00}, // end of element
	"SNAME":        {0x12, 0x06}, // name of structure reference
	"COLROW":       {0x13, 0x02}, // number of colrow[0]=columns and colrow[1]=rows in array reference
	"NODE":         {0x15, 0x00}, // begin of node element
	"TEXTTYPE":     {0x16, 0x02}, // texttype number
	"PRESENTATION": {0x17, 0x01}, // text presentation, font
	"STRINGBODY":   {0x19, 0x06}, // character string for text element
	"STRANS":       {0x1a, 0x01}, // array reference, structure reference and text transform flags
	"MAG":          {0x1b, 0x05}, // magnification factor for text and references
	"ANGLE":        {0x1c, 0x05}, // rotation angle for text and references
	"REFLIBS":      {0x1f, 0x06}, // name of referenced libraries
	"FONTS":        {0x20, 0x06}, // name of text fonts definition files
	"PATHTYPE":     {0x21, 0x02}, // type of PATH element end ( rounded, square)
	"GENERATIONS":  {0x22, 0x02}, // number of deleted structure ?????
	"ATTRTABLE":    {0x23, 0x06}, // attribute table, used in combination with element properties
	"ELFLAGS":      {0x26, 0x01}, // template data
	"NODETYPE":     {0x2a, 0x02}, // node type number for NODE element
	"PROPATTR":     {0x2b, 0x02}, // attribute number
	"PROPVALUE":    {0x2c, 0x06}, // attribute name
	"BOX":          {0x2d, 0x00}, // begin of box element
	"BOXTYPE":      {0x2e, 0x02}, // boxtype for box element
	"PLEX":         {0x2f, 0x03}, // plex number
	"BGNEXTN":      {0x30, 0x03}, // path type 4 extension start
	"ENDEXTN":      {0x31, 0x03}, // path type 4 extension end
	"TAPENUM":      {0x32, 0x02}, // tape number
	"TAPECODE":     {0x33, 0x02}, // tape code
	"FORMAT":       {0x36, 0x02}, // format type
	"MASK":         {0x37, 0x06}, // list of layers
	"ENDMASKS":     {0x38, 0x00}, // end of MASK
	"LIBDIRSIZE":   {0x39, 0x02}, // contains the number of pages in the Library directory
	"SRFNAME":      {0x3a, 0x06}, // contains the name of the Sticks Rules File, if one is bound to the library
	"LIBSECUR":     {0x3b, 0x02}, // contains an array of Access Control List (ACL) data
}

type ElementType int

const (
	PolygonType ElementType = iota
	PathType
	LabelType
	SRefType
	ARefType
	UnsupportedType
)

const HEADERSIZE = 4

type Record struct {
	Size     uint16
	Datatype string
	Data     []byte
}

func (r Record) GetData() (any, error) {
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
		return "No data", nil
	case "BGNSTR":
		return getDataSlice[int16](r)
	case "STRNAME":
		return getDataString(r)
	case "ENDSTR":
		return "No data", nil
	case "BOUNDARY":
		return "No data", nil
	case "PATH":
		return "No data", nil
	case "SREF":
		return "No data", nil
	case "AREF":
		return "No data", nil
	case "TEXT":
		return "No data", nil
	case "LAYER":
		return getDataPoint[int16](r)
	case "DATATYPE":
		return getDataPoint[int16](r)
	case "WIDTH":
		return getDataPoint[int32](r)
	case "XY":
		return getDataSlice[int32](r)
	case "ENDEL":
		return "No data", nil
	case "SNAME":
		return getDataString(r)
	case "COLROW":
		return getDataSlice[int16](r)
	case "NODE":
		return "No data", nil
	case "TEXTTYPE":
		return getDataPoint[int16](r)
	case "PRESENTATION":
		return getDataPoint[uint16](r)
	case "STRINGBODY":
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
		return "No data", nil
	case "BOXTYPE":
		return getDataPoint[int16](r)
	case "PLEX":
		return getDataPoint[int32](r)
	case "BGNEXTN":
		return getDataPoint[int32](r)
	case "ENDEXTN":
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
		return "No data", nil
	default:
		panic("unexpected datatype")
	}
}

func (r Record) String() string {
	data, err := r.GetData()
	if err != nil {
		return fmt.Sprintf("error occurred: %v", err)
	}
	return fmt.Sprintf("%s, Bytes: %d Data: %v", r.Datatype, r.Size, data)
}

func (r Record) Bytes() []byte {
	resultBytes := []byte{}
	resultBytes = append(resultBytes, byte(r.Size>>8), byte(r.Size))
	resultBytes = append(resultBytes, RecordTypesBytes[r.Datatype]...)
	resultBytes = append(resultBytes, r.Data...)
	return resultBytes
}

type Element interface {
	String() string
	GetData() any
	Records() ([]Record, error)
	GetLayer() string
	Type() ElementType
}

type Library struct {
	Header     int16
	BgnLib     []int16
	LibName    string
	Units      []float64
	Structures map[string]*Structure
}

func (l Library) String() string {
	structureInfo := "\n"
	for _, structure := range l.Structures {
		structureInfo += "      " + structure.StrName + "\n"
		structureElements := structure.ListElements()
		structureInfo += structureElements
	}
	return fmt.Sprintf(`Library:
   Version: %d
   Name: %s
   Units: %v
   Structures:%s`, l.Header, l.LibName, l.Units, structureInfo)
}
func (l Library) Records() ([]Record, error) {
	records, err := fieldsToRecords(l)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for library: %v", err)
	}
	return wrapStartEnd("BGNLIB", records), nil
}

type Structure struct {
	BgnStr   []int16
	StrName  string
	Elements []Element
}

func (s Structure) String() string {
	return fmt.Sprintf("%s, %v", s.StrName, s.Elements)
}

func (s Structure) ListElements() string {
	result := ""
	for _, v := range s.Elements {
		result += "         " + v.String() + "\n"
	}
	return result
}
func (s Structure) Records() ([]Record, error) {
	records, err := fieldsToRecords(s)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for structure: %v", err)
	}
	return wrapStartEnd("BGNSTR", records), nil
}

type Boundary struct {
	ElFlags  uint16
	Plex     int32
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
func (b Boundary) Records() ([]Record, error) {
	records, err := fieldsToRecords(b)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for boundary: %v", err)
	}
	return wrapStartEnd("BOUNDARY", records), nil
}
func (b Boundary) GetLayer() string {
	return fmt.Sprintf("%d/%d", b.Layer, b.Datatype)
}
func (b Boundary) Type() ElementType {
	return PolygonType
}
func (b Boundary) GetPoints() []int32 {
	return b.XY
}

type Path struct {
	ElFlags  uint16
	Plex     int32
	Layer    int16
	Datatype int16
	Pathtype int16
	Bgnextn  int32
	Endextn  int32
	Width    int32
	XY       []int32
}

func (p Path) GetData() any {
	return p.XY
}
func (p Path) GetWidth() int32 {
	return p.Width
}
func (p Path) GetPathType() int16 {
	return p.Pathtype
}

func (p Path) String() string {
	return fmt.Sprintf("Path - ElFlags: %v, Plex: %v, Layer: %v, Datatype: %v, Pathtype: %v, Width: %v, XY: %v",
		p.ElFlags, p.Plex, p.Layer, p.Datatype, p.Pathtype, p.Width, p.XY)
}
func (p Path) Records() ([]Record, error) {
	records, err := fieldsToRecords(p)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for path: %v", err)
	}
	return wrapStartEnd("PATH", records), nil
}
func (p Path) GetLayer() string {
	return fmt.Sprintf("%d/%d", p.Layer, p.Datatype)
}
func (p Path) Type() ElementType {
	return PathType
}

type Text struct {
	ElFlags      uint16
	Plex         int32
	Layer        int16
	Texttype     int16
	Presentation uint16
	Strans       uint16
	Mag          float64
	Angle        float64
	XY           []int32
	StringBody   string
}

func (t Text) GetData() any {
	return t.StringBody
}
func (t Text) String() string {
	return fmt.Sprintf("Text - ElFlags: %v, Plex: %v, Layer: %v, XY: %v, String: %v", t.ElFlags, t.Plex, t.Layer, t.XY, t.StringBody)
}
func (t Text) Records() ([]Record, error) {
	records, err := fieldsToRecords(t)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for text: %v", err)
	}
	return wrapStartEnd("TEXT", records), err
}
func (t Text) GetLayer() string {
	return fmt.Sprintf("%d/%d", t.Layer, t.Texttype)
}
func (t Text) Type() ElementType {
	return LabelType
}

type Node struct {
	ElFlags  uint16
	Plex     int32
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
func (n Node) Records() ([]Record, error) {
	records, err := fieldsToRecords(n)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for node: %v", err)
	}
	return wrapStartEnd("NODE", records), nil
}
func (n Node) GetLayer() string {
	return fmt.Sprintf("%d/%d", n.Layer, n.Nodetype)
}
func (n Node) Type() ElementType {
	return UnsupportedType
}

type Box struct {
	ElFlags uint16
	Plex    int32
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
func (b Box) Records() ([]Record, error) {
	records, err := fieldsToRecords(b)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for box: %v", err)
	}
	return wrapStartEnd("BOX", records), nil
}
func (b Box) GetLayer() string {
	return fmt.Sprintf("%d/%d", b.Layer, b.Boxtype)
}
func (b Box) Type() ElementType {
	return PolygonType
}
func (b Box) GetPoints() []int32 {
	return b.XY
}

type SRef struct {
	ElFlags uint16
	Plex    int32
	Sname   string
	Strans  uint16 // Strans flags do nothing yet
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
func (s SRef) Records() ([]Record, error) {
	records, err := fieldsToRecords(s)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for sref: %v", err)
	}
	return wrapStartEnd("SREF", records), nil
}
func (s SRef) GetLayer() string {
	return "cellref"
}
func (s SRef) Type() ElementType {
	return SRefType
}

type ARef struct {
	ElFlags uint16
	Plex    int32
	Sname   string
	Strans  uint16 // Strans flags do nothing yet
	Mag     float64
	Angle   float64
	Colrow  []int16
	XY      []int32
}

func (a ARef) GetData() any {
	return a.XY
}
func (a ARef) String() string {
	return fmt.Sprintf("ARef - ElFlags: %v, Plex: %v, Sname: %s, Strans: %v, Mag: %v, Angle: %v, Colrow: %v, XY: %v",
		a.ElFlags, a.Plex, a.Sname, a.Strans, a.Mag, a.Angle, a.Colrow, a.XY)
}
func (a ARef) Records() ([]Record, error) {
	records, err := fieldsToRecords(a)
	if err != nil {
		return []Record{}, fmt.Errorf("could not produce records for aref: %v", err)
	}
	return wrapStartEnd("AREF", records), nil
}
func (a ARef) GetLayer() string {
	return "cellref"
}
func (a ARef) Type() ElementType {
	return ARefType
}

// Wraps a record slice with their start record "{ELEMENTTYPE}" and end record "ENDEL"
func wrapStartEnd(elementType string, records []Record) []Record {
	wrappedRecords := []Record{}
	if elementType == "BGNSTR" {
		return append(records, Record{Size: 4, Datatype: "ENDSTR", Data: []byte{}})
	} else if elementType == "BGNLIB" {
		return append(records, Record{Size: 4, Datatype: "ENDLIB", Data: []byte{}})
	} else {
		wrappedRecords = append(wrappedRecords, Record{Size: 4, Datatype: elementType, Data: []byte{}})
		wrappedRecords = append(wrappedRecords, records...)
		wrappedRecords = append(wrappedRecords, Record{Size: 4, Datatype: "ENDEL", Data: []byte{}})
	}
	return wrappedRecords
}
