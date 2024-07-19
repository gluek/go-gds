package gds

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func gotypeToBytes(value any) []byte {
	switch reflect.TypeOf(value) {
	case reflect.TypeOf(int16(0)):
		return []byte{byte(value.(int16) >> 8), byte(value.(int16))}
	case reflect.TypeOf(uint16(0)):
		return []byte{byte(value.(uint16) >> 8), byte(value.(uint16))}
	case reflect.TypeOf(int32(0)):
		return []byte{
			byte(value.(int32) >> 24),
			byte(value.(int32) >> 16),
			byte(value.(int32) >> 8),
			byte(value.(int32)),
		}
	case reflect.TypeOf(float64(0)):
		return bitsToByteArray(encodeReal(value.(float64)))
	case reflect.TypeOf(""):
		return []byte(value.(string))
	case reflect.TypeOf([]int16{}):
		returnSlice := []byte{}
		for _, v := range value.([]int16) {
			returnSlice = append(returnSlice, gotypeToBytes(v)...)
		}
		return returnSlice
	case reflect.TypeOf([]int32{}):
		returnSlice := []byte{}
		for _, v := range value.([]int32) {
			returnSlice = append(returnSlice, gotypeToBytes(v)...)
		}
		return returnSlice
	case reflect.TypeOf([]float64{}):
		returnSlice := []byte{}
		for _, v := range value.([]float64) {
			returnSlice = append(returnSlice, gotypeToBytes(v)...)
		}
		return returnSlice
	default:
		panic(fmt.Sprintf("datatype not supported by GDSII: %v", reflect.TypeOf(value)))
	}
}

func bitsToByteArray(i uint64) []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

// Convert bits which represents 8-byte real with 1-bit sign, 7-bit exponent and 56-bit mantissa to IEEE754 float64
func decodeReal(bits uint64) float64 {
	sign := 1.0
	if uint64(bits&0x80_00_00_00_00_00_00_00) > 0 {
		sign = -1.0
	}
	exponent := int8(bits >> 56)
	rangingFactor := float64(uint64(0b00000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000))
	mantissa := float64(bits&0x00_ff_ff_ff_ff_ff_ff_ff) / rangingFactor
	return sign * mantissa * math.Pow(16, math.Abs(float64(exponent))-64)
}

// Convert IEEE754 float64 to 8-byte real with 1-bit sign, 7-bit exponent and 56-bit mantissa
func encodeReal(fl float64) uint64 {
	if fl == 0.0 {
		return uint64(0x00_00_00_00_00_00_00_00)
	}
	valueScientifc := fmt.Sprintf("%b", math.Abs(fl))
	sign := uint64(0x00_00_00_00_00_00_00_00)
	if fl < 0 {
		sign = uint64(0x80_00_00_00_00_00_00_00)
	}
	parts := strings.Split(valueScientifc, "p")
	factor, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		log.Fatalf("could not parse magnitude: %v", err)
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalf("could not parse exponent: %v", err)
	}
	factor = factor << 11
	newExp := int64((exp + 53) / 4)
	expRemainder := (exp + 53) % 4
	if expRemainder > 0 {
		factor = factor >> (4 - expRemainder)
		newExp++
	} else if expRemainder < 0 {
		factor = factor >> (-expRemainder)
	}
	newExpUint := uint64(newExp + 64)
	return uint64(sign | (factor >> 8) | (newExpUint << 56))
}

func getRealSlice(data Record) []float64 {
	initSlice := make([]uint64, int((data.Size-HEADERSIZE)/8))
	finalSlice := make([]float64, len(initSlice))

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &initSlice)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	for i, number := range initSlice {
		finalSlice[i] = decodeReal(number)
	}
	return finalSlice
}

func getDataSlice[T any](data Record) []T {
	var typeInit T
	typeSize := reflect.TypeOf(typeInit).Size()
	result := make([]T, int((data.Size-HEADERSIZE)/uint16(typeSize)))
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	return result
}

func getRealPoint(data Record) float64 {
	var number uint64

	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &number)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	return decodeReal(number)
}

func getDataPoint[T any](data Record) T {
	var result T
	reader := bytes.NewReader(data.Data)
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		log.Fatalf("could not read binary data: %v", err)
	}
	return result
}

func getDataString(data Record) string {
	return string(data.Data)
}
