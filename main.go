package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// https://www.artwork.com/gdsii/gdsii/index.htm
// https://boolean.klaasholwerda.nl/interface/bnf/gdsformat.html

func readGDS(f *os.File) (*Library, error) {
	var err error

	reader := bufio.NewReader(f)
	library, err := libraryFactory(reader)
	if err != nil {
		return nil, err
	}
	return library, nil
}

func main() {
	file, err := os.Open("klayout_test.gds")
	if err != nil {
		log.Fatalf("could not open gds file: %v", err)
	}
	defer file.Close()

	library, err := readGDS(file)
	if err != nil {
		log.Fatalf("could not read gds: %v", err)
	}
	fmt.Print(library)
}
