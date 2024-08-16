package gds

import (
	"math"
	"reflect"
)

func resolveSRef(lib *Library, container any, references []Reference) {
	callingRef := references[len(references)-1]
	for _, element := range lib.Structures[callingRef.GetSname()].Elements {
		if element.Type() == PolygonType {
			// Basically checks if the calling function is GetCellData or GetLayermapPolygons
			var layermap map[string]*PolygonLayer
			if reflect.TypeOf(container) == reflect.TypeOf(map[string]*PolygonLayer{}) {
				layermap = container.(map[string]*PolygonLayer)
			} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
				layermap = container.(*CellData).Polygons
			} else {
				continue
			}

			var points []int32 = element.(Polygon).GetPoints()
			for i := range len(references) {
				ref := references[len(references)-i-1] // resolve from last to first
				points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
			}

			layer, ok := layermap[element.GetLayer()]
			if ok {
				layer.appendPolygon(points)
			} else {
				layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
			}
		} else if element.Type() == PathType {
			// Basically checks if the calling function is GetCellData or GetLayermapPaths
			var layermap map[string]*PathLayer
			if reflect.TypeOf(container) == reflect.TypeOf(map[string]*PathLayer{}) {
				layermap = container.(map[string]*PathLayer)
			} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
				layermap = container.(*CellData).Paths
			} else {
				continue
			}

			var points []int32 = element.(*Path).XY
			var width float64 = float64(element.(*Path).Width)
			for i := range len(references) {
				ref := references[len(references)-i-1] // resolve from last to first
				points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
				width = width * ref.(*SRef).Mag
			}

			layer, ok := layermap[element.GetLayer()]
			if ok {
				layer.appendPath(points, element.(*Path).Pathtype, int32(width), element.(*Path).Bgnextn, element.(*Path).Endextn)
			} else {
				layermap[element.GetLayer()] = &PathLayer{
					Enabled:   true,
					Paths:     [][]int32{points},
					PathTypes: []int16{element.(*Path).Pathtype},
					ExtBegin:  []int32{element.(*Path).Bgnextn},
					ExtEnd:    []int32{element.(*Path).Endextn},
					Widths:    []int32{int32(width)},
				}
			}
		} else if element.Type() == LabelType {
			// Basically checks if the calling function is GetCellData or GetLayermapLabels
			var layermap map[string]*LabelLayer
			if reflect.TypeOf(container) == reflect.TypeOf(map[string]*LabelLayer{}) {
				layermap = container.(map[string]*LabelLayer)
			} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
				layermap = container.(*CellData).Labels
			} else {
				continue
			}

			layer, ok := layermap[element.GetLayer()]
			var points []int32 = element.(*Text).XY
			for i := range len(references) {
				ref := references[len(references)-i-1] // resolve from last to first
				points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
			}
			if ok {
				layer.appendLabel(points, element.(*Text).StringBody)
			} else {
				layermap[element.GetLayer()] = &LabelLayer{
					Enabled:     true,
					Labels:      []string{element.(*Text).StringBody},
					LabelCoords: [][]int32{points},
				}
			}
		} else if element.Type() == SRefType {
			newReferences := append(references, element.(*SRef))
			resolveSRef(lib, container, newReferences)
		} else if element.Type() == ARefType {
			newReferences := append(references, element.(*ARef))
			resolveARef(lib, container, newReferences)
		}
	}
}

func resolveARef(lib *Library, container any, references []Reference) {
	callingRef := references[len(references)-1].(*ARef)
	var xshift, yshift int32

	nCol := callingRef.Colrow[0]
	nRow := callingRef.Colrow[1]

	refPoint := callingRef.XY[:2]
	mulColSpacing := callingRef.XY[2:4]
	mulRowSpacing := callingRef.XY[4:]
	mulColSpacing = []int32{mulColSpacing[0] - callingRef.XY[0], mulColSpacing[1] - callingRef.XY[1]}
	mulRowSpacing = []int32{mulRowSpacing[0] - callingRef.XY[0], mulRowSpacing[1] - callingRef.XY[1]}
	for i := range nCol {
		for j := range nRow {
			xshift = int32(math.Round(float64(refPoint[0]) + float64(i)*float64(mulColSpacing[0])/float64(nCol) + float64(j)*float64(mulRowSpacing[0])/float64(nRow)))
			yshift = int32(math.Round(float64(refPoint[1]) + float64(i)*float64(mulColSpacing[1])/float64(nCol) + float64(j)*float64(mulRowSpacing[1])/float64(nRow)))
			newReferences := references

			for _, element := range lib.Structures[callingRef.GetSname()].Elements {
				newSref := &SRef{
					ElFlags: callingRef.ElFlags,
					Plex:    callingRef.Plex,
					Sname:   callingRef.Sname,
					Strans:  callingRef.Strans,
					Mag:     callingRef.Mag,
					Angle:   callingRef.Angle,
					XY:      []int32{xshift, yshift},
				}
				newReferences[len(newReferences)-1] = newSref

				if element.Type() == PolygonType {
					// Basically checks if the calling function is GetCellData or GetLayermapPolygons
					var layermap map[string]*PolygonLayer
					if reflect.TypeOf(container) == reflect.TypeOf(map[string]*PolygonLayer{}) {
						layermap = container.(map[string]*PolygonLayer)
					} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
						layermap = container.(*CellData).Polygons
					} else {
						continue
					}

					var points []int32 = element.(Polygon).GetPoints()
					for i := range len(newReferences) {
						ref := references[len(newReferences)-i-1] // resolve from last to first
						points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
					}

					layer, ok := layermap[element.GetLayer()]
					if ok {
						layer.appendPolygon(points)
					} else {
						layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
					}
				} else if element.Type() == PathType {
					// Basically checks if the calling function is GetCellData or GetLayermapPaths
					var layermap map[string]*PathLayer
					if reflect.TypeOf(container) == reflect.TypeOf(map[string]*PathLayer{}) {
						layermap = container.(map[string]*PathLayer)
					} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
						layermap = container.(*CellData).Paths
					} else {
						continue
					}

					var points []int32 = element.(*Path).XY
					var width float64 = float64(element.(*Path).Width)
					for i := range len(newReferences) {
						ref := references[len(newReferences)-i-1] // resolve from last to first
						points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
						width = width * ref.(*SRef).Mag
					}

					layer, ok := layermap[element.GetLayer()]
					if ok {
						layer.appendPath(points, element.(*Path).Pathtype, int32(width), element.(*Path).Bgnextn, element.(*Path).Endextn)
					} else {
						layermap[element.GetLayer()] = &PathLayer{
							Enabled:   true,
							Paths:     [][]int32{points},
							PathTypes: []int16{element.(*Path).Pathtype},
							ExtBegin:  []int32{element.(*Path).Bgnextn},
							ExtEnd:    []int32{element.(*Path).Endextn},
							Widths:    []int32{int32(width)},
						}
					}
				} else if element.Type() == LabelType {
					// Basically checks if the calling function is GetCellData or GetLayermapLabels
					var layermap map[string]*LabelLayer
					if reflect.TypeOf(container) == reflect.TypeOf(map[string]*LabelLayer{}) {
						layermap = container.(map[string]*LabelLayer)
					} else if reflect.TypeOf(container) == reflect.TypeOf(&CellData{}) {
						layermap = container.(*CellData).Labels
					} else {
						continue
					}

					layer, ok := layermap[element.GetLayer()]
					var points []int32 = element.(*Text).XY
					for i := range len(newReferences) {
						ref := references[len(newReferences)-i-1] // resolve from last to first
						points = transformPoints(points, ref.(*SRef).XY[0], ref.(*SRef).XY[1], ref.(*SRef).Strans, ref.(*SRef).Mag, ref.(*SRef).Angle)
					}

					if ok {
						layer.appendLabel(points, element.(*Text).StringBody)
					} else {
						layermap[element.GetLayer()] = &LabelLayer{
							Enabled:     true,
							Labels:      []string{element.(*Text).StringBody},
							LabelCoords: [][]int32{points},
						}
					}
				} else if element.Type() == SRefType {
					newReferences := append(references, element.(*SRef))
					resolveSRef(lib, container, newReferences)
				} else if element.Type() == ARefType {
					newReferences := append(references, element.(*ARef))
					resolveARef(lib, container, newReferences)
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
		// x-Axis mirroring
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
