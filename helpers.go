package gds

import (
	"math"
)

// Wraps a record slice with their start record "{ELEMENTTYPE}" and end record "ENDEL"
func wrapStartEnd(elementType string, records []Record) []Record {
	wrappedRecords := []Record{}
	if elementType == "BGNSTR" {
		return append(records, Record{Size: 4, Datatype: "ENDSTR", Data: []byte{}})
	} else if elementType == "BGNLIB" {
		return append(records, Record{Size: 4, Datatype: "ENDLIB", Data: []byte{}})
	} else {
		wrappedRecords = append(wrappedRecords, Record{Size: 4, Datatype: elementType, Data: []byte{}})
		wrappedRecords = append(wrappedRecords, records...)
		wrappedRecords = append(wrappedRecords, Record{Size: 4, Datatype: "ENDEL", Data: []byte{}})
	}
	return wrappedRecords
}

func resolveSRefPolygons(lib *Library, layermap map[string]*PolygonLayer, ref *SRef) {
	for _, element := range lib.Structures[ref.Sname].Elements {
		if element.Type() == PolygonType {
			points := transformPoints(element.(Polygon).GetPoints(), ref.XY[0], ref.XY[1], ref.Strans, ref.Mag, ref.Angle)
			layer, ok := layermap[element.GetLayer()]
			if ok {
				layer.appendPolygon(points)
			} else {
				layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
			}
		}
	}
}

func resolveARefPolygons(lib *Library, layermap map[string]*PolygonLayer, ref *ARef) {
	nCol := ref.Colrow[0]
	nRow := ref.Colrow[1]

	refPoint := ref.XY[:2]
	mulColSpacing := ref.XY[2:4]
	mulRowSpacing := ref.XY[4:]
	mulColSpacing = []int32{mulColSpacing[0] - ref.XY[0], mulColSpacing[1] - ref.XY[1]}
	mulRowSpacing = []int32{mulRowSpacing[0] - ref.XY[0], mulRowSpacing[1] - ref.XY[1]}
	for i := range nCol {
		for j := range nRow {
			for _, element := range lib.Structures[ref.Sname].Elements {
				if element.Type() == PolygonType {
					points := transformPoints(element.(Polygon).GetPoints(),
						int32(math.Round(float64(refPoint[0])+float64(i)*float64(mulColSpacing[0])/float64(nCol)+float64(j)*float64(mulRowSpacing[0])/float64(nRow))),
						int32(math.Round(float64(refPoint[1])+float64(i)*float64(mulColSpacing[1])/float64(nCol)+float64(j)*float64(mulRowSpacing[1])/float64(nRow))),
						ref.Strans, ref.Mag, ref.Angle)
					layer, ok := layermap[element.GetLayer()]
					if ok {
						layer.appendPolygon(points)
					} else {
						layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
					}
				}
			}
		}
	}
}

func transformPoints(array []int32, xshift int32, yshift int32, strans uint16, mag float64, angle float64) []int32 {
	radians := angle * math.Pi / 180
	transformedArray := make([]int32, len(array))
	for i := 0; i < len(array); i += 2 {
		var x, y float64
		// y-Axis mirroring
		x = float64(array[i])
		y = float64(array[i+1]) * (0.5 - float64((strans >> 15))) * 2
		// rotation + magnification
		x_temp := (float64(x)*math.Cos(radians) - float64(y)*math.Sin(radians)) * mag
		y = (float64(x)*math.Sin(radians) + float64(y)*math.Cos(radians)) * mag
		x = x_temp
		// shift
		transformedArray[i] = int32(math.Round(x + float64(xshift)))
		transformedArray[i+1] = int32(math.Round(y + float64(yshift)))

	}
	return transformedArray
}
