package gds

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type CellData struct {
	Layers   []string                 `json:"layers"`
	Polygons map[string]*PolygonLayer `json:"polygons"`
	Paths    map[string]*PathLayer    `json:"paths"`
	Labels   map[string]*LabelLayer   `json:"labels"`
}

// implement sort interface
type ByLayer []string

func errHandlerAtoi(a int, e error) int {
	if e != nil {
		log.Println("could not convert str to int")
	}
	return a
}

func compareLayers(i string, j string) int {
	iSplit := strings.Split(i, "/")
	jSplit := strings.Split(j, "/")
	iNum := errHandlerAtoi(strconv.Atoi(iSplit[0]))*1000 + errHandlerAtoi(strconv.Atoi(iSplit[1]))
	jNum := errHandlerAtoi(strconv.Atoi(jSplit[0]))*1000 + errHandlerAtoi(strconv.Atoi(jSplit[1]))
	if iNum < jNum {
		return -1
	} else if iNum > jNum {
		return 1
	} else {
		return 0
	}
}

func (c *CellData) PopulateLayers() {
	layerSet := map[string]bool{}
	for k, _ := range c.Polygons {
		_, ok := layerSet[k]
		if !ok {
			layerSet[k] = true
		}
	}
	for k, _ := range c.Paths {
		_, ok := layerSet[k]
		if !ok {
			layerSet[k] = true
		}
	}
	for k, _ := range c.Labels {
		_, ok := layerSet[k]
		if !ok {
			layerSet[k] = true
		}
	}
	for k, _ := range layerSet {
		c.Layers = append(c.Layers, k)
	}
	slices.SortFunc(c.Layers, compareLayers)
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
	ExtBegin  []int32   `json:"extbegin"`
	ExtEnd    []int32   `json:"extend"`
	Paths     [][]int32 `json:"paths"`
}

func (p *PathLayer) appendPath(xy []int32, pathtype int16, width int32, extbegin int32, extend int32) ([][]int32, []int16, []int32, []int32, []int32) {
	p.Paths = append(p.Paths, xy)
	p.PathTypes = append(p.PathTypes, pathtype)
	p.Widths = append(p.Widths, width)
	p.ExtBegin = append(p.ExtBegin, extbegin)
	p.ExtEnd = append(p.ExtEnd, extend)
	return p.Paths, p.PathTypes, p.Widths, p.ExtBegin, p.ExtEnd
}
func (p PathLayer) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", p.Enabled, p.PathTypes, p.Widths, p.Paths)
}

// TODO: Include vertical/horizontal anchor
type LabelLayer struct {
	Enabled     bool      `json:"enable"`
	Labels      []string  `json:"labels"`
	LabelCoords [][]int32 `json:"xy"`
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
