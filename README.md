# go-gds

go-gds is a library to encode/decode Calma GDSII binary files written in Go.

## Features

- Decoding binary files to go types
- Encoding go types to binary
- High-level api functions to extract geometries, separated into cells and layers

## Missing

- Functions to extract labels and paths
- Functionality to manipulate geometries, paths, labels, cells and layers
- Entire path starting from geometries to binary

## Example

```go
fh, err := os.Open("mygds.gds")
if err != nil {
    log.Printf("could not open test gds file: %v", err)
}
defer fh.Close()

library, err := ReadGDS(fh)
if err != nil {
    log.Printf("could not parse gds file: %v", err)
}
// Prints basic library information, cells and contained elements
fmt.Print(library)


// Returns map of layer/datatype -> polygons for specified cell
layermap, err := library.GetLayermapPolygons("mycell")
if err != nil {
    log.Printf("could not extract layermap polygons: %v", %v)
}

// Show layers and polygons
fmt.Print(layermap)
```
