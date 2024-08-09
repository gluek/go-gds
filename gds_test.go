package gds

import (
	"fmt"
	"os"
	"testing"

	svg "github.com/ajstarks/svgo"
)

const testFile = "klayout_test.gds"

func TestReadGDS(t *testing.T) {
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

func TestGetLayerPolygons(t *testing.T) {
	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	polygons, err := library.GetLayermapPolygons("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}
	fmt.Print(polygons)
}

func TestGetLayerPaths(t *testing.T) {
	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	paths, err := library.GetLayermapPaths("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}
	fmt.Print(paths)
}

func TestGetLayerLabels(t *testing.T) {
	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	labels, err := library.GetLayermapLabels("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}
	fmt.Print(labels)
}

func TestGetCellData(t *testing.T) {
	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}
	celldata, err := library.GetCellData("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}
	fmt.Print(celldata.Polygons, "\n")
	fmt.Print(celldata.Paths, "\n")
	fmt.Print(celldata.Labels, "\n")
}

func TestDrawPolygons(t *testing.T) {
	fh, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("could not open test gds file: %v", err)
	}
	defer fh.Close()

	library, err := ReadGDS(fh)
	if err != nil {
		t.Fatalf("could not parse gds file: %v", err)
	}

	layermapPolygons, err := library.GetLayermapPolygons("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}

	fhSVG, err := os.Create("test.svg")
	if err != nil {
		t.Fatalf("could not generate svg")
	}
	colormap := []string{"black", "blue", "red", "yellow", "cyan", "magenta", "purple", "green", "orange"}
	width := 25000
	height := 25000
	canvas := svg.New(fhSVG)
	canvas.Start(width, height)
	j := 0
	for _, v := range layermapPolygons {
		for _, poly := range v.Polygons {
			var x, y []int
			for i := 0; i < len(poly); i += 2 {
				x = append(x, int(poly[i]))
				y = append(y, int(poly[i+1]))
			}
			canvas.Polygon(x, y, fmt.Sprintf("stroke-width:0.1%%;fill:none;stroke:%s", colormap[j]))
		}
		j++
	}

	layermapLabels, err := library.GetLayermapLabels("top")
	if err != nil {
		t.Fatalf("could not get layermap polygons: %v", err)
	}
	for _, v := range layermapLabels {
		for i := range len(v.Labels) {
			canvas.Text(int(v.LabelCoords[i][0]), int(v.LabelCoords[i][1]), v.Labels[i], "dominant-baseline:middle;text-anchor:middle;font-size:100px;fill:white")
		}
	}

	canvas.End()
}

func TestReadRecords(t *testing.T) {
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
