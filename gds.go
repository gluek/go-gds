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
