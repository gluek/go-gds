package gds

import (
	"bufio"
	"fmt"
	"os"
)

// https://www.artwork.com/gdsii/gdsii/index.htm
// https://boolean.klaasholwerda.nl/interface/bnf/gdsformat.html

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

func (l Library) GetLayers() map[string][]Element {
	result := map[string][]Element{}
	for _, structure := range l.Structures {
		for _, element := range structure.Elements {
			elementSlice, ok := result[element.GetLayer()]
			if ok {
				result[element.GetLayer()] = append(elementSlice, element)
			} else {
				result[element.GetLayer()] = []Element{element}
			}
		}
	}
	return result
}

func (l Library) GetStructures() map[string]Structure {
	result := map[string]Structure{}
	for _, structure := range l.Structures {
		result[structure.StrName] = structure
	}
	return result
}
