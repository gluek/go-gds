package gds

import "math"

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

func resolveSRefPolygons(lib *Library, ref SRef) map[string][][]int32 {
	result := map[string][][]int32{}
	for _, element := range lib.Structures[ref.Sname].Elements {
		if element.Type() == PolygonType {
			points := transformPoints(element.(Polygon).GetPoints(), ref.XY[0], ref.XY[1], ref.Strans, ref.Mag, ref.Angle)
			_, ok := result[element.GetLayer()]
			if ok {
				result[element.GetLayer()] = append(result[element.GetLayer()], points)
			} else {
				result[element.GetLayer()] = [][]int32{points}
			}
		}
	}
	return result
}

func transformPoints(array []int32, xshift int32, yshift int32, strans uint16, mag float64, angle float64) []int32 {
	radians := angle * math.Pi / 180
	transformedArray := make([]int32, len(array))
	for i := 0; i < len(array); i += 2 {
		var x, y int32
		// y-Axis mirroring
		x = array[i]
		y = array[i+1] * int32((0.5-float32((strans>>15)))*2)
		// rotation + magnification
		x = int32((float64(x)*math.Cos(radians) - float64(x)*math.Sin(radians)) * mag)
		y = int32((float64(y)*math.Cos(radians) + float64(y)*math.Sin(radians)) * mag)
		// shift
		transformedArray[i] = int32(x) + xshift
		transformedArray[i+1] = int32(y) + yshift

	}
	return transformedArray
}
