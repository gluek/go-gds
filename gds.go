package gds

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.artwork.com/gdsii/gdsii/index.htm
// https://boolean.klaasholwerda.nl/interface/bnf/gdsformat.html
// http://bitsavers.informatik.uni-stuttgart.de/pdf/calma/GDS_II_Stream_Format_Manual_6.0_Feb87.pdf

func ReadGDS(f *os.File) (*Library, error) {
	var err error

	reader := bufio.NewReader(f)
	library, err := decodeLibrary(reader)
	if err != nil {
		return nil, err
	}
	return library, nil
}

func ReadRecords(f *os.File) ([]Record, error) {
	records := []Record{}
	reader := bufio.NewReader(f)
OuterLoop:
	for {
		record, err := decodeRecord(reader)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
		if record.Datatype == "ENDLIB" {
			break OuterLoop
		}
	}
	return records, nil
}

func WriteGDS(f *os.File, lib *Library) error {
	writer := bufio.NewWriter(f)
	records, err := lib.Records()
	if err != nil {
		return fmt.Errorf("could not write GDSII file: %v", err)
	}
	for _, record := range records {
		_, err := writer.Write(record.Bytes())
		if err != nil {
			return fmt.Errorf("could not write record %v to file: %v", record, err)
		}
	}
	writer.Flush()
	return nil
}

func (l Library) GetLayermapPolygons(cell string) (map[string]*PolygonLayer, error) {
	result := map[string]*PolygonLayer{}
	structure, ok := l.Structures[cell]
	if !ok {
		return map[string]*PolygonLayer{}, fmt.Errorf("cell with name %s does not exist", cell)
	}
	for _, element := range structure.Elements {
		if element.Type() == PolygonType {
			layer, ok := result[element.GetLayer()]
			if ok {
				layer.appendPolygon(element.(Polygon).GetPoints())
			} else {
				result[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{element.(Polygon).GetPoints()}}
			}
		} else if element.Type() == SRefType {
			resolveSRef(&l, result, element.(*SRef))
		} else if element.Type() == ARefType {
			resolveARef(&l, result, element.(*ARef))
		}
	}
	return result, nil
}

func (l Library) GetLayermapPaths(cell string) (map[string]*PathLayer, error) {
	result := map[string]*PathLayer{}
	structure, ok := l.Structures[cell]
	if !ok {
		return map[string]*PathLayer{}, fmt.Errorf("cell with name %s does not exist", cell)
	}
	for _, element := range structure.Elements {
		if element.Type() == PathType {
			layer, ok := result[element.GetLayer()]
			if ok {
				layer.appendPath(element.(*Path).GetData().([]int32), element.(*Path).GetPathType(), element.(*Path).GetWidth())
			} else {
				result[element.GetLayer()] = &PathLayer{
					Enabled:   true,
					Paths:     [][]int32{element.(*Path).GetData().([]int32)},
					PathTypes: []int16{element.(*Path).GetPathType()},
					Widths:    []int32{element.(*Path).GetWidth()},
				}
			}
		} else if element.Type() == SRefType {
			resolveSRef(&l, result, element.(*SRef))
		} else if element.Type() == ARefType {
			resolveARef(&l, result, element.(*ARef))
		}
	}
	return result, nil
}
