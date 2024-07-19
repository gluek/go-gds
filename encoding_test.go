package gds

import (
	"bufio"
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

func mockFilehandler(data []byte) *bufio.Reader {
	return bufio.NewReader(bytes.NewReader(data))
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

func TestEncodeReal(t *testing.T) {
	assertEqual(t, uint64(0b01000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000), encodeReal(0.5))
	assertEqual(t, uint64(0b11000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000), encodeReal(-0.5))
	assertEqual(t, uint64(0b01000001_00011000_00000000_00000000_00000000_00000000_00000000_00000000), encodeReal(1.5))
	assertEqual(t, uint64(0b01000010_01100100_00000000_00000000_00000000_00000000_00000000_00000000), encodeReal(100.0))
	assertEqual(t, uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000), encodeReal(0.0))
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

func TestGotypeToBytes(t *testing.T) {
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f)}, gotypeToBytes(int16(0x0f_0f)))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f)}, gotypeToBytes(uint16(0x0f_0f)))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)}, gotypeToBytes(int32(0x0f_0f_0f_0f)))
	assertEqualByteSlice(t, []byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)}, gotypeToBytes([]int16{0x0f_0f, 0x0f_0f}))
	assertEqualByteSlice(t,
		[]byte{byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f), byte(0x0f)},
		gotypeToBytes([]int32{0x0f_0f_0f_0f, 0x0f_0f_0f_0f}))
	assertEqualByteSlice(t,
		[]byte{byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00)},
		gotypeToBytes(float64(0.5)))
	assertEqualByteSlice(t,
		[]byte{byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00),
			byte(0x40), byte(0x80), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00), byte(0x00)},
		gotypeToBytes([]float64{0.5, 0.5}))
	assertEqualByteSlice(t, []byte("test123"), gotypeToBytes("test123"))
}

func TestRecordsToRecords(t *testing.T) {
	arefTest := ARef{
		ElFlags: 0,
		Plex:    0,
		Sname:   "Test",
		Strans:  0,
		Mag:     1.0,
		Angle:   0.0,
		Colrow:  []int16{0, 0},
		XY:      []int32{0, 0, 1, 1},
	}
	recordsAref := fieldsToRecords(arefTest)
	recordsAref = append(recordsAref, Record{Size: 4, Datatype: "ENDEL", Data: []byte{}})
	arefReader := mockFilehandler(recordsToBytes(recordsAref))
	arefNew, err := decodeAREF(arefReader)
	if err != nil {
		t.Fatalf("error decoding aref: %v", err)
	}
	if arefTest.String() != arefNew.String() {
		t.Fatalf("%v not equal to %v", arefTest, arefNew)
	}

	textTest := Text{
		ElFlags:      0,
		Plex:         0,
		Layer:        1,
		Texttype:     0,
		Presentation: 0,
		Strans:       0,
		Mag:          1.0,
		Angle:        0.0,
		XY:           []int32{0, 0, 1, 1},
		StringBody:   "Test",
	}
	recordsText := fieldsToRecords(textTest)
	recordsText = append(recordsText, Record{Size: 4, Datatype: "ENDEL", Data: []byte{}})
	textReader := mockFilehandler(recordsToBytes(recordsText))
	textNew, err := decodeText(textReader)
	if err != nil {
		t.Fatalf("error decoding aref: %v", err)
	}
	if textTest.String() != textNew.String() {
		t.Fatalf("%v not equal to %v", textNew, textNew)
	}

}
