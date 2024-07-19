package gds

import (
	"bufio"
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
