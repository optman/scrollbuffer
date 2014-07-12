package scrollbuffer

import (
	"testing"
)

func testReadWrite(t *testing.T, pos int) {

	buf := New(1024, 0)

	data := make([]byte, 10)
	text := "123456789"

	buf.Write(pos, ([]byte)(text))
	buf.Read(pos, data[0:len(text)])

	if string(data[0:len(text)]) != text {
		t.Fatal()
	}
}

func TestReadWrite(t *testing.T) {

	testReadWrite(t, 0)
	testReadWrite(t, 10)
	testReadWrite(t, 1022)
	testReadWrite(t, 1023)
	testReadWrite(t, 1024)
	testReadWrite(t, 1025)
	testReadWrite(t, 1025*1000)
}

//balance forward and backward.
func TestScollBackAndForward(t *testing.T) {

	buf := New(20, 10)
	buf.Write(3, ([]byte)("456"))
	buf.Write(0, ([]byte)("123"))
	buf.Write(6, ([]byte)("789"))

	if buf.DataRanges().ContainRange(MakeRange(0, 9)) == false {
		t.Fatal()
	}

	data := make([]byte, 9)
	buf.Read(0, data)
	if string(data) != "123456789" {
		t.Fatal()
	}
}

//optimize for forward
func TestScollForward(t *testing.T) {

	buf := New(9, 0)
	buf.Write(0, ([]byte)("123"))
	buf.Write(3, ([]byte)("456"))
	buf.Write(6, ([]byte)("789"))

	if buf.DataRanges().ContainRange(MakeRange(0, 9)) == false {
		t.Fatal()
	}

	data := make([]byte, 9)
	buf.Read(0, data)
	if string(data) != "123456789" {
		t.Fatal()
	}
}

//optimize for backward
func TestScollBackward(t *testing.T) {

	buf := New(12, 9)
	buf.Write(6, ([]byte)("789"))
	buf.Write(3, ([]byte)("456"))
	buf.Write(0, ([]byte)("123"))

	if buf.DataRanges().ContainRange(MakeRange(0, 9)) == false {
		t.Fatal()
	}

	data := make([]byte, 9)
	buf.Read(0, data)
	if string(data) != "123456789" {
		t.Fatal()
	}
}
