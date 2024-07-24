package gds

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

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
	bData := make([]byte, size-4)

	n, err = io.ReadFull(reader, bData)
	if n != int(size-4) {
		return nil, fmt.Errorf("wrong number of data bytes for %s/%s. expected: %d got: %d", datatype, RecordTypes[datatype], size-4, n)
	}
	if err != nil {
		return nil, fmt.Errorf("could not read data bytes: %v", err)
	}
	return &Record{Size: size, Datatype: RecordTypes[datatype], Data: bData}, nil
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
		default:
			return nil, fmt.Errorf("could not decode Path/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Sref/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Aref/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Text/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Node/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Box/%s: unkown datatype", newRecord.Datatype)
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
			return nil, fmt.Errorf("could not decode Structure/%s: unkown datatype", newRecord.Datatype)
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
		Structures: []Structure{},
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
			library.Structures = append(library.Structures, *element)
		default:
			return nil, fmt.Errorf("could not decode Library/%s: unkown datatype", newRecord.Datatype)
		}
	}
	return &library, nil
}
