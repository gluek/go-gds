package gds

import "fmt"

type PolygonLayer struct {
	Enabled  bool
	Polygons [][]int32
}

type PathLayer struct {
	Enabled  bool
	Polygons []Polygon
}

type LabelLayer struct {
	Enabled  bool
	Polygons []Polygon
}

func (p *PolygonLayer) appendPolygon(poly []int32) [][]int32 {
	p.Polygons = append(p.Polygons, poly)
	return p.Polygons
}
func (p PolygonLayer) String() string {
	return fmt.Sprintf("%v, %v", p.Enabled, p.Polygons)
}

type Polygon interface {
	GetPoints() []int32
}

type Reference interface {
	GetSname() string
}
