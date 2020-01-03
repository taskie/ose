package osplus

import (
	"testing"
)

func TestBufIOContainerIn(t *testing.T) {
	io := NewBufIOContainer()
	io.InBuf.WriteByte(42)
	bs := []byte{0}
	io.In().Read(bs)
	if bs[0] != 42 {
		t.Fatal("In() failed")
	}
}

func TestBufIOContainerOut(t *testing.T) {
	io := NewBufIOContainer()
	io.Out().Write([]byte{42})
	b, err := io.OutBuf.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 42 {
		t.Fatal("Out() failed")
	}
}

func TestBufIOContainerErr(t *testing.T) {
	io := NewBufIOContainer()
	io.Err().Write([]byte{42})
	b, err := io.ErrBuf.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 42 {
		t.Fatal("Err() failed")
	}
}
