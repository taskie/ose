package osplus

import (
	"testing"
)

func TestFakeIOIn(t *testing.T) {
	io := NewFakeIO()
	io.InBuf.WriteByte(42)
	bs := []byte{0}
	io.In().Read(bs)
	if bs[0] != 42 {
		t.Fatal("In() failed")
	}
}

func TestFakeIOOut(t *testing.T) {
	io := NewFakeIO()
	io.Out().Write([]byte{42})
	b, err := io.OutBuf.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 42 {
		t.Fatal("Out() failed")
	}
}

func TestFakeIOErr(t *testing.T) {
	io := NewFakeIO()
	io.Err().Write([]byte{42})
	b, err := io.ErrBuf.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 42 {
		t.Fatal("Err() failed")
	}
}
