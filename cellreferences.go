package gds

import (
	"math"
	"reflect"
)

func resolveSRef(lib *Library, layermap any, ref *SRef) {
	for _, element := range lib.Structures[ref.Sname].Elements {
		if element.Type() == PolygonType && reflect.TypeOf(layermap) == reflect.TypeOf(map[string]*PolygonLayer{}) {
			points := transformPoints(element.(Polygon).GetPoints(), ref.XY[0], ref.XY[1], ref.Strans, ref.Mag, ref.Angle)
			layermap := layermap.(map[string]*PolygonLayer)
			layer, ok := layermap[element.GetLayer()]
			if ok {
				layer.appendPolygon(points)
			} else {
				layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
			}
		} else if element.Type() == PathType && reflect.TypeOf(layermap) == reflect.TypeOf(map[string]*PathLayer{}) {
			points := transformPoints(element.(*Path).XY, ref.XY[0], ref.XY[1], ref.Strans, ref.Mag, ref.Angle)
			layermap := layermap.(map[string]*PathLayer)
			layer, ok := layermap[element.GetLayer()]
			if ok {
				layer.appendPath(points, element.(*Path).GetPathType(), int32(float64(element.(*Path).GetWidth())*ref.Mag))
			} else {
				layermap[element.GetLayer()] = &PathLayer{
					Enabled:   true,
					Paths:     [][]int32{points},
					PathTypes: []int16{element.(*Path).GetPathType()},
					Widths:    []int32{int32(float64(element.(*Path).GetWidth()) * ref.Mag)},
				}
			}
		} else if element.Type() == SRefType {
			resolveSRef(lib, layermap, element.(*SRef))
		} else if element.Type() == ARefType {
			resolveARef(lib, layermap, element.(*ARef))
		}
	}
}

func resolveARef(lib *Library, layermap any, ref *ARef) {
	var xshift, yshift int32

	nCol := ref.Colrow[0]
	nRow := ref.Colrow[1]

	refPoint := ref.XY[:2]
	mulColSpacing := ref.XY[2:4]
	mulRowSpacing := ref.XY[4:]
	mulColSpacing = []int32{mulColSpacing[0] - ref.XY[0], mulColSpacing[1] - ref.XY[1]}
	mulRowSpacing = []int32{mulRowSpacing[0] - ref.XY[0], mulRowSpacing[1] - ref.XY[1]}
	for i := range nCol {
		for j := range nRow {
			xshift = int32(math.Round(float64(refPoint[0]) + float64(i)*float64(mulColSpacing[0])/float64(nCol) + float64(j)*float64(mulRowSpacing[0])/float64(nRow)))
			yshift = int32(math.Round(float64(refPoint[1]) + float64(i)*float64(mulColSpacing[1])/float64(nCol) + float64(j)*float64(mulRowSpacing[1])/float64(nRow)))
			for _, element := range lib.Structures[ref.Sname].Elements {
				if element.Type() == PolygonType && reflect.TypeOf(layermap) == reflect.TypeOf(map[string]*PolygonLayer{}) {
					layermap := layermap.(map[string]*PolygonLayer)
					points := transformPoints(element.(Polygon).GetPoints(),
						xshift,
						yshift,
						ref.Strans, ref.Mag, ref.Angle)
					layer, ok := layermap[element.GetLayer()]
					if ok {
						layer.appendPolygon(points)
					} else {
						layermap[element.GetLayer()] = &PolygonLayer{Enabled: true, Polygons: [][]int32{points}}
					}
				} else if element.Type() == PathType && reflect.TypeOf(layermap) == reflect.TypeOf(map[string]*PathLayer{}) {
					points := transformPoints(element.(*Path).XY, xshift, yshift, ref.Strans, ref.Mag, ref.Angle)
					layermap := layermap.(map[string]*PathLayer)
					layer, ok := layermap[element.GetLayer()]
					if ok {
						layer.appendPath(points, element.(*Path).GetPathType(), int32(float64(element.(*Path).GetWidth())*ref.Mag))
					} else {
						layermap[element.GetLayer()] = &PathLayer{
							Enabled:   true,
							Paths:     [][]int32{points},
							PathTypes: []int16{element.(*Path).GetPathType()},
							Widths:    []int32{int32(float64(element.(*Path).GetWidth()) * ref.Mag)},
						}
					}
				} else if element.Type() == SRefType {
					resolveSRef(lib, layermap, element.(*SRef))
				} else if element.Type() == ARefType {
					resolveARef(lib, layermap, element.(*ARef))
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
