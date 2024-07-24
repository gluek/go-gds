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

func TestReadRecords(t *testing.T) {
	testFile := "klayout_test.gds"

	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	records, err := ReadRecords(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	for _, v := range records {
		fmt.Println(v)
	}
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
	err = WriteGDS(fhWrite, library)
	if err != nil {
		t.Fatalf("could not write library to gds file: %v", err)
	}
	fhWrite.Close()
	err = os.Remove("gds_test.gds")
	if err != nil {
		t.Fatalf("could not delete gds_test file")
	}
}
