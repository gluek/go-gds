package gds

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

func recordFactory(reader *bufio.Reader) (*Record, error) {
	var n int
	var err error

	bSize := make([]byte, 2)
	n, err = reader.Read(bSize)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		log.Println("wrong number of size bytes")
	}

	size := binary.BigEndian.Uint16(bSize)

	bDatatype := make([]byte, 2)
	n, err = reader.Read(bDatatype)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		log.Println("wrong number of datatype bytes")
	}
	datatype := hex.EncodeToString(bDatatype)
	if size < 4 {
		log.Println("size smaller 4")
	}
	bData := make([]byte, size-4)

	n, err = io.ReadFull(reader, bData)
	if n != int(size-4) {
		log.Printf("wrong number of data bytes. expected: %d got: %d\n", size-4, n)
	}
	if err != nil {
		return nil, err
	}
	return &Record{Size: size, Datatype: RecordTypes[datatype], Data: bData}, nil
}

func boundaryFactory(reader *bufio.Reader) (*Boundary, error) {
	boundary := Boundary{
		ElFlags:  0,
		Plex:     0,
		Layer:    -1,
		Datatype: -1,
		XY:       []int32{},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			boundary.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			boundary.Plex = newRecord.GetData().(int16)
		case "LAYER":
			boundary.Layer = newRecord.GetData().(int16)
		case "DATATYPE":
			boundary.Datatype = newRecord.GetData().(int16)
		case "XY":
			boundary.XY = newRecord.GetData().([]int32)
		}
	}
	return &boundary, nil
}

func pathFactory(reader *bufio.Reader) (*Path, error) {
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
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			path.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			path.Plex = newRecord.GetData().(int16)
		case "LAYER":
			path.Layer = newRecord.GetData().(int16)
		case "DATATYPE":
			path.Datatype = newRecord.GetData().(int16)
		case "PATHTYPE":
			path.Pathtype = newRecord.GetData().(int16)
		case "WIDTH":
			path.Width = newRecord.GetData().(int32)
		case "XY":
			path.XY = newRecord.GetData().([]int32)
		}
	}
	return &path, nil
}

func srefFactory(reader *bufio.Reader) (*SRef, error) {
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
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			sref.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			sref.Plex = newRecord.GetData().(int16)
		case "SNAME":
			sref.Sname = newRecord.GetData().(string)
		case "STRANS":
			sref.Strans = newRecord.GetData().(uint16)
		case "MAG":
			sref.Mag = newRecord.GetData().(float64)
		case "ANGLE":
			sref.Angle = newRecord.GetData().(float64)
		case "XY":
			sref.XY = newRecord.GetData().([]int32)
		}
	}
	return &sref, nil
}

func arefFactory(reader *bufio.Reader) (*ARef, error) {
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
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			aref.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			aref.Plex = newRecord.GetData().(int16)
		case "SNAME":
			aref.Sname = newRecord.GetData().(string)
		case "STRANS":
			aref.Strans = newRecord.GetData().(uint16)
		case "MAG":
			aref.Mag = newRecord.GetData().(float64)
		case "ANGLE":
			aref.Angle = newRecord.GetData().(float64)
		case "COLROW":
			aref.Colrow = newRecord.GetData().([]int16)
		case "XY":
			aref.XY = newRecord.GetData().([]int32)
		}
	}
	return &aref, nil
}

func textFactory(reader *bufio.Reader) (*Text, error) {
	text := Text{
		ElFlags: 0,
		Plex:    0,
		Layer:   -1,
		Textbody: Textbody{
			Texttype:     -1,
			Presentation: 0,
			Strans:       0,
			Mag:          1,
			Angle:        0,
			StringBody:   "",
			XY:           []int32{},
		},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			text.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			text.Plex = newRecord.GetData().(int16)
		case "LAYER":
			text.Layer = newRecord.GetData().(int16)
		case "TEXTTYPE":
			text.Texttype = newRecord.GetData().(int16)
		case "PRESENTATION":
			text.Presentation = newRecord.GetData().(uint16)
		case "STRANS":
			text.Strans = newRecord.GetData().(uint16)
		case "MAG":
			text.Mag = newRecord.GetData().(float64)
		case "ANGLE":
			text.Angle = newRecord.GetData().(float64)
		case "STRING":
			text.StringBody = newRecord.GetData().(string)
		case "XY":
			text.XY = newRecord.GetData().([]int32)
		}
	}
	return &text, nil
}

func nodeFactory(reader *bufio.Reader) (*Node, error) {
	node := Node{
		ElFlags:  0,
		Plex:     0,
		Layer:    -1,
		Nodetype: -1,
		XY:       []int32{},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			node.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			node.Plex = newRecord.GetData().(int16)
		case "LAYER":
			node.Layer = newRecord.GetData().(int16)
		case "NODETYPE":
			node.Nodetype = newRecord.GetData().(int16)
		case "XY":
			node.XY = newRecord.GetData().([]int32)
		}
	}
	return &node, nil
}

func boxFactory(reader *bufio.Reader) (*Box, error) {
	box := Box{
		ElFlags: 0,
		Plex:    0,
		Layer:   -1,
		Boxtype: -1,
		XY:      []int32{},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDEL":
			break OuterLoop
		case "ELFLAGS":
			box.ElFlags = newRecord.GetData().(uint16)
		case "PLEX":
			box.Plex = newRecord.GetData().(int16)
		case "LAYER":
			box.Layer = newRecord.GetData().(int16)
		case "BOXTYPE":
			box.Boxtype = newRecord.GetData().(int16)
		case "XY":
			box.XY = newRecord.GetData().([]int32)
		}
	}
	return &box, nil
}

func structureFactory(reader *bufio.Reader, bgnStrRecord *Record) (*Structure, error) {
	structure := Structure{
		BgnStr:   bgnStrRecord.GetData().([]int16),
		StrName:  "Unknown",
		Elements: []Element{},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDSTR":
			break OuterLoop
		case "STRNAME":
			structure.StrName = newRecord.GetData().(string)
		case "BOUNDARY":
			element, err := boundaryFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "PATH":
			element, err := pathFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "SREF":
			element, err := srefFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "AREF":
			element, err := arefFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "TEXT":
			element, err := textFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "NODE":
			element, err := nodeFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		case "BOX":
			element, err := boxFactory(reader)
			if err != nil {
				return nil, err
			}
			structure.Elements = append(structure.Elements, element)
		}
	}
	return &structure, nil
}

func libraryFactory(reader *bufio.Reader) (*Library, error) {
	library := Library{
		Header:     0,
		BgnLib:     []int16{},
		LibName:    "Unknown",
		Units:      []float64{},
		Structures: []*Structure{},
	}
OuterLoop:
	for {
		newRecord, err := recordFactory(reader)
		if err != nil {
			return nil, err
		}
		switch newRecord.Datatype {
		case "ENDLIB":
			break OuterLoop
		case "BGNLIB":
			library.BgnLib = newRecord.GetData().([]int16)
		case "LIBNAME":
			library.LibName = newRecord.GetData().(string)
		case "UNITS":
			library.Units = newRecord.GetData().([]float64)
		case "BGNSTR":
			element, err := structureFactory(reader, newRecord)
			if err != nil {
				return nil, err
			}
			library.Structures = append(library.Structures, element)
		}
	}
	return &library, nil
}
