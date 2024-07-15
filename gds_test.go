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

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	fmt.Print(library)
}
