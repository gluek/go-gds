package gds

import "fmt"

type GDSData struct {
	Layers   []string                 `json:"layers"`
	Polygons map[string]*PolygonLayer `json:"polygons"`
	Paths    map[string]*PathLayer    `json:"paths"`
	Labels   map[string]*LabelLayer   `json:"labels"`
}

type PolygonLayer struct {
	Enabled  bool      `json:"enable"`
	Polygons [][]int32 `json:"polygons"`
}

func (p *PolygonLayer) appendPolygon(poly []int32) [][]int32 {
	p.Polygons = append(p.Polygons, poly)
	return p.Polygons
}
func (p PolygonLayer) String() string {
	return fmt.Sprintf("%v, %v", p.Enabled, p.Polygons)
}

type PathLayer struct {
	Enabled   bool      `json:"enable"`
	PathTypes []int16   `json:"types"`
	Widths    []int32   `json:"widths"`
	Paths     [][]int32 `json:"paths"`
}

func (p *PathLayer) appendPath(xy []int32, pathtype int16, width int32) ([][]int32, []int16, []int32) {
	p.Paths = append(p.Paths, xy)
	p.PathTypes = append(p.PathTypes, pathtype)
	p.Widths = append(p.Widths, width)
	return p.Paths, p.PathTypes, p.Widths
}
func (p PathLayer) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", p.Enabled, p.PathTypes, p.Widths, p.Paths)
}

type LabelLayer struct {
	Enabled     bool      `json:"enable"`
	Labels      []string  `json:"labels"`
	LabelCoords [][]int32 `json:"labelxy"`
}

func (l *LabelLayer) appendLabel(xy []int32, text string) ([][]int32, []string) {
	l.Labels = append(l.Labels, text)
	l.LabelCoords = append(l.LabelCoords, xy)
	return l.LabelCoords, l.Labels
}
func (l LabelLayer) String() string {
	return fmt.Sprintf("%v, %v, %v", l.Enabled, l.Labels, l.LabelCoords)
}

// Polygon interface includes Element Type BOUNDARY and BOX
type Polygon interface {
	GetPoints() []int32
}

type Reference interface {
	GetSname() string
}
