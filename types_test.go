package main

import (
	"testing"
)

func assertEqual(t *testing.T, a any, b any) {
	if a == b {
		return
	}
	t.Fatalf("%v != %v", a, b)
}

func TestGetRealSlice(t *testing.T) {
	data := Record{
		Size:     4 + 5*8,
		Datatype: "",
		Data: []byte{
			byte(0b01000000), byte(0b10000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), // 0.5
			byte(0b11000000), byte(0b10000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), // -0.5
			byte(0b01000001), byte(0b00011000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), // 1.5
			byte(0b01000010), byte(0b01100100), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), // 100.0
			byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000), // 0.0
		},
	}
	result := getRealSlice(data)
	assertEqual(t, result[0], 0.5)
	assertEqual(t, result[1], -0.5)
	assertEqual(t, result[2], 1.5)
	assertEqual(t, result[3], 100.0)
	assertEqual(t, result[4], 0.0)
}
