package gds

import (
	"bytes"
	"testing"
)

func assertEqual(t *testing.T, a any, b any) {
	if a == b {
		return
	}
	t.Fatalf("%v != %v", a, b)
}

func assertEqualByteSlice(t *testing.T, a []byte, b []byte) {
	if bytes.Equal(a, b) {
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
	result, err := getRealSlice(data)
	if err != nil {
		t.Fatalf("could not get real slice: %v", err)
	}
	assertEqual(t, result[0], 0.5)
	assertEqual(t, result[1], -0.5)
	assertEqual(t, result[2], 1.5)
	assertEqual(t, result[3], 100.0)
	assertEqual(t, result[4], 0.0)
}

func TestEncodeReal(t *testing.T) {
	assertEqual(t, uint64(0b01000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000), ignoreError(encodeReal(0.5)))
	assertEqual(t, uint64(0b11000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000), ignoreError(encodeReal(-0.5)))
	assertEqual(t, uint64(0b01000001_00011000_00000000_00000000_00000000_00000000_00000000_00000000), ignoreError(encodeReal(1.5)))
	assertEqual(t, uint64(0b01000010_01100100_00000000_00000000_00000000_00000000_00000000_00000000), ignoreError(encodeReal(100.0)))
	assertEqual(t, uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000), ignoreError(encodeReal(0.0)))
}

func TestDecodeReal(t *testing.T) {
	assertEqual(t, 0.5, decodeReal(uint64(0b01000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000)))
	assertEqual(t, -0.5, decodeReal(uint64(0b11000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000)))
	assertEqual(t, 1.5, decodeReal(uint64(0b01000001_00011000_00000000_00000000_00000000_00000000_00000000_00000000)))
	assertEqual(t, 100.0, decodeReal(uint64(0b01000010_01100100_00000000_00000000_00000000_00000000_00000000_00000000)))
	assertEqual(t, 0.0, decodeReal(uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000)))
}

func TestBitsToByteArray(t *testing.T) {
	expected := []byte{byte(0b01000000), byte(0b10000000), byte(0b00000000), byte(0b00000000),
		byte(0b00000000), byte(0b00000000), byte(0b00000000), byte(0b00000000)}
	got := bitsToByteArray(uint64(0b01000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000))
	assertEqualByteSlice(t, expected, got)
}

func ignoreError[T any](val T, _ error) T {
	return val
}

func TestGotypeToBytes(t *testing.T) {
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f)}, ignoreError(gotypeToBytes(int16(0x0f_0f))))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f)}, ignoreError(gotypeToBytes(uint16(0x0f_0f))))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)}, ignoreError(gotypeToBytes(int32(0x0f_0f_0f_0f))))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)}, ignoreError(gotypeToBytes([]int16{0x0f_0f, 0x0f_0f})))
	assertEqualByteSlice(t,
		[]byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)},
		ignoreError(gotypeToBytes([]int32{0x0f_0f_0f_0f, 0x0f_0f_0f_0f})))
	assertEqualByteSlice(t,
		[]byte{byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00)},
		ignoreError(gotypeToBytes(float64(0.5))))
	assertEqualByteSlice(t,
		[]byte{byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00),
			byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00)},
		ignoreError(gotypeToBytes([]float64{0.5, 0.5})))
	// length always must be even, therefore 7 char strings are padded with a 0-byte
	testString := []byte("test123")
	testString = append(testString, byte(0))
	assertEqualByteSlice(t, testString, ignoreError(gotypeToBytes("test123")))
	assertEqualByteSlice(t, []byte("test1234"), ignoreError(gotypeToBytes("test1234")))
}
