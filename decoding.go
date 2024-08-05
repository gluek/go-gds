package gds

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

// Convert bits which represents 8-byte real with 1-bit sign, 7-bit exponent and 56-bit mantissa to IEEE754 float64
func decodeReal(bits uint64) float64 {
	sign := 1.0
	if uint64(bits&0x80_00_00_00_00_00_00_00) > 0 {
		sign = -1.0
	}
	exponent := int8(bits >> 56)
	rangingFactor := float64(uint64(0b00000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000))
	mantissa := float64(bits&0x00_ff_ff_ff_ff_ff_ff_ff) / rangingFactor
	return sign * mantissa * math.Pow(16, math.Abs(float64(exponent))-64)
}

func getRealSlice(data Record) ([]float64, error) {
	initSlice := make([]uint64, int((data.Size-HEADERSIZE)/8))
	finalSlice := make([]float64, len(initSlice))

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &initSlice)
	if err != nil {
		return finalSlice, fmt.Errorf("could not read binary data: %v", err)
	}
	for i, number := range initSlice {
		finalSlice[i] = decodeReal(number)
	}
	return finalSlice, nil
}

func getDataSlice[T any](data Record) ([]T, error) {
	var typeInit T
	typeSize := reflect.TypeOf(typeInit).Size()
	result := make([]T, int((data.Size-HEADERSIZE)/uint16(typeSize)))
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		return result, fmt.Errorf("could not read binary data: %v", err)
	}
	return result, nil
}

func getRealPoint(data Record) (float64, error) {
	var number uint64

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &number)
	if err != nil {
		return float64(0), fmt.Errorf("could not read binary data: %v", err)
	}
	return decodeReal(number), nil
}

func getDataPoint[T any](data Record) (T, error) {
	var result T
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		return result, fmt.Errorf("could not read binary data. RecordType: %s, %v", data.Datatype, err)
	}
	return result, nil
}

func getDataString(data Record) (string, error) {
	return strings.TrimRight(string(data.Data), string(byte(0))), nil
}

func decodeRecord(reader *bufio.Reader) (*Record, error) {
	var n int
	var err error

	bSize := make([]byte, 2)
	n, err = reader.Read(bSize)
	if err != nil {
		return nil, fmt.Errorf("could not read size bytes: %v", err)
	}
	if n != 2 {
		return nil, fmt.Errorf("wrong number of size bytes")
	}

	size := binary.BigEndian.Uint16(bSize)

	bDatatype := make([]byte, 2)
	n, err = reader.Read(bDatatype)
	if err != nil {
		return nil, fmt.Errorf("could not read datatype bytes: %v", err)
	}
	if n != 2 {
		return nil, fmt.Errorf("wrong number of datatype bytes")
	}
	datatype := hex.EncodeToString(bDatatype)
	if size < 4 {
		return nil, fmt.Errorf("size smaller than 4 bytes")
	}

	datatypeString, ok := RecordTypes[datatype]
	if !ok {
		return nil, fmt.Errorf("unknown datatype: %s", datatype)
	}

	bData := make([]byte, size-4)

	n, err = io.ReadFull(reader, bData)
	if n != int(size-4) {
		return nil, fmt.Errorf("wrong number of data bytes for %s/%s. expected: %d got: %d", datatype, datatypeString, size-4, n)
	}
	if err != nil {
		return nil, fmt.Errorf("could not read data bytes: %v", err)
	}
	return &Record{Size: size, Datatype: datatypeString, Data: bData}, nil
}

func decodeBoundary(reader *bufio.Reader) (*Boundary, error) {
	boundary := Boundary{
		ElFlags:  0,
		Plex:     0,
		Layer:    -1,
		Datatype: -1,
		XY:       []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Boundary/%s: %v", newRecord.Datatype, err)
			}
			boundary.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Boundary/%s: %v", newRecord.Datatype, err)
			}
			boundary.Plex = data.(int32)
		case "LAYER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Boundary/%s: %v", newRecord.Datatype, err)
			}
			boundary.Layer = data.(int16)
		case "DATATYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Boundary/%s: %v", newRecord.Datatype, err)
			}
			boundary.Datatype = data.(int16)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Boundary/%s: %v", newRecord.Datatype, err)
			}
			boundary.XY = data.([]int32)
		}
	}
	return &boundary, nil
}

func decodePath(reader *bufio.Reader) (*Path, error) {
	path := Path{
		ElFlags:  0,
		Plex:     0,
		Layer:    -1,
		Datatype: -1,
		Pathtype: -1,
		Width:    -1,
		XY:       []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.Plex = data.(int32)
		case "LAYER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.Layer = data.(int16)
		case "DATATYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.Datatype = data.(int16)
		case "PATHTYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.Pathtype = data.(int16)
		case "WIDTH":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.Width = data.(int32)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Path/%s: %v", newRecord.Datatype, err)
			}
			path.XY = data.([]int32)
		case "BGNEXTN":
			continue
		case "ENDEXTN":
			continue
		default:
			return nil, fmt.Errorf("could not decode Path/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &path, nil
}

func decodeSREF(reader *bufio.Reader) (*SRef, error) {
	sref := SRef{
		ElFlags: 0,
		Plex:    0,
		Sname:   "",
		Strans:  0,
		Mag:     1,
		Angle:   0,
		XY:      []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.Plex = data.(int32)
		case "SNAME":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.Sname = data.(string)
		case "STRANS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.Strans = data.(uint16)
		case "MAG":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.Mag = data.(float64)
		case "ANGLE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.Angle = data.(float64)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Sref/%s: %v", newRecord.Datatype, err)
			}
			sref.XY = data.([]int32)
		default:
			return nil, fmt.Errorf("could not decode Sref/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &sref, nil
}

func decodeAREF(reader *bufio.Reader) (*ARef, error) {
	aref := ARef{
		ElFlags: 0,
		Plex:    0,
		Sname:   "",
		Strans:  0,
		Mag:     1,
		Angle:   0,
		Colrow:  []int16{},
		XY:      []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Plex = data.(int32)
		case "SNAME":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Sname = data.(string)
		case "STRANS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Strans = data.(uint16)
		case "MAG":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Mag = data.(float64)
		case "ANGLE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Angle = data.(float64)
		case "COLROW":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.Colrow = data.([]int16)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Aref/%s: %v", newRecord.Datatype, err)
			}
			aref.XY = data.([]int32)
		default:
			return nil, fmt.Errorf("could not decode Aref/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &aref, nil
}

func decodeText(reader *bufio.Reader) (*Text, error) {
	text := Text{
		ElFlags:      0,
		Plex:         0,
		Layer:        -1,
		Texttype:     -1,
		Presentation: 0,
		Strans:       0,
		Mag:          1,
		Angle:        0,
		StringBody:   "",
		XY:           []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Plex = data.(int32)
		case "LAYER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Layer = data.(int16)
		case "TEXTTYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Texttype = data.(int16)
		case "PRESENTATION":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Presentation = data.(uint16)
		case "STRANS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Strans = data.(uint16)
		case "MAG":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Mag = data.(float64)
		case "ANGLE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.Angle = data.(float64)
		case "STRINGBODY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.StringBody = data.(string)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Text/%s: %v", newRecord.Datatype, err)
			}
			text.XY = data.([]int32)
		default:
			return nil, fmt.Errorf("could not decode Text/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &text, nil
}

func decodeNode(reader *bufio.Reader) (*Node, error) {
	node := Node{
		ElFlags:  0,
		Plex:     0,
		Layer:    -1,
		Nodetype: -1,
		XY:       []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Node/%s: %v", newRecord.Datatype, err)
			}
			node.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Node/%s: %v", newRecord.Datatype, err)
			}
			node.Plex = data.(int32)
		case "LAYER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Node/%s: %v", newRecord.Datatype, err)
			}
			node.Layer = data.(int16)
		case "NODETYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Node/%s: %v", newRecord.Datatype, err)
			}
			node.Nodetype = data.(int16)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Node/%s: %v", newRecord.Datatype, err)
			}
			node.XY = data.([]int32)
		default:
			return nil, fmt.Errorf("could not decode Node/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &node, nil
}

func decodeBox(reader *bufio.Reader) (*Box, error) {
	box := Box{
		ElFlags: 0,
		Plex:    0,
		Layer:   -1,
		Boxtype: -1,
		XY:      []int32{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("error decoding record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Box/%s: %v", newRecord.Datatype, err)
			}
			box.ElFlags = data.(uint16)
		case "PLEX":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Box/%s: %v", newRecord.Datatype, err)
			}
			box.Plex = data.(int32)
		case "LAYER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Box/%s: %v", newRecord.Datatype, err)
			}
			box.Layer = data.(int16)
		case "BOXTYPE":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Box/%s: %v", newRecord.Datatype, err)
			}
			box.Boxtype = data.(int16)
		case "XY":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Box/%s: %v", newRecord.Datatype, err)
			}
			box.XY = data.([]int32)
		default:
			return nil, fmt.Errorf("could not decode Box/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &box, nil
}

func decodeStructure(reader *bufio.Reader, bgnStrRecord *Record) (*Structure, error) {
	data, err := bgnStrRecord.GetData()
	if err != nil {
		return nil, fmt.Errorf("could not decode Structure/BGNSTR: %v", err)
	}
	structure := Structure{
		BgnStr:   data.([]int16),
		StrName:  "Unknown",
		Elements: []Element{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDSTR":
			break OuterLoop
		case "BGNSTR":
			continue //should not happen outside of tests
		case "STRNAME":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.StrName = data.(string)
		case "BOUNDARY":
			element, err := decodeBoundary(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "PATH":
			element, err := decodePath(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "SREF":
			element, err := decodeSREF(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "AREF":
			element, err := decodeAREF(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "TEXT":
			element, err := decodeText(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "NODE":
			element, err := decodeNode(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		case "BOX":
			element, err := decodeBox(reader)
			if err != nil {
				return nil, fmt.Errorf("could not decode Structure/%s: %v", newRecord.Datatype, err)
			}
			structure.Elements = append(structure.Elements, element)
		default:
			return nil, fmt.Errorf("could not decode Structure/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &structure, nil
}

func decodeLibrary(reader *bufio.Reader) (*Library, error) {
	library := Library{
		Header:     0,
		BgnLib:     []int16{},
		LibName:    "Unknown",
		Units:      []float64{},
		Structures: map[string]*Structure{},
	}
OuterLoop:
	for {
		newRecord, err := decodeRecord(reader)
		if err != nil {
			return nil, fmt.Errorf("could not decode record: %v", err)
		}
		switch newRecord.Datatype {
		case "ENDLIB":
			break OuterLoop
		case "HEADER":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Library/%s: %v", newRecord.Datatype, err)
			}
			library.Header = data.(int16)
		case "BGNLIB":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Library/%s: %v", newRecord.Datatype, err)
			}
			library.BgnLib = data.([]int16)
		case "LIBNAME":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Library/%s: %v", newRecord.Datatype, err)
			}
			library.LibName = data.(string)
		case "UNITS":
			data, err := newRecord.GetData()
			if err != nil {
				return nil, fmt.Errorf("could not decode Library/%s: %v", newRecord.Datatype, err)
			}
			library.Units = data.([]float64)
		case "BGNSTR":
			element, err := decodeStructure(reader, newRecord)
			if err != nil {
				return nil, fmt.Errorf("could not decode Library/%s: %v", newRecord.Datatype, err)
			}
			library.Structures[element.StrName] = element
		default:
			return nil, fmt.Errorf("could not decode Library/%s: unknown datatype", newRecord.Datatype)
		}
	}
	return &library, nil
}
