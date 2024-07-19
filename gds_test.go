package gds

import (
	"fmt"
	"os"
	"testing"
)

func TestReadGDS(t *testing.T) {
	testFile := "klayout_test.gds"

	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	fmt.Print(library)
}

func TestWriteGDS(t *testing.T) {
	testFile := "klayout_test.gds"

	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}

	fhWrite, err := os.Create("gds_test.gds")
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fhWrite.Close()
	err = WriteGDS(fhWrite, library)
	if err != nil {
		t.Fatalf("could not write library to gds file: %v", err)
	}
}
