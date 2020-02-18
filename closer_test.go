package ose_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/taskie/ose"
)

func TestReadCloser(t *testing.T) {
	expected := []byte("ABC")
	buf := bytes.NewBuffer(expected)
	closed := false
	rc := ose.NewReadCloser(buf, func(r io.Reader) error {
		closed = true
		return nil
	})
	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, bs) {
		t.Fatalf("invalid value: %v (expected: %v)", bs, expected)
	}
	if closed {
		t.Fatal("must not be closed")
	}
	rc.Close()
	if !closed {
		t.Fatal("must be closed")
	}
}

func TestNopReadCloser(t *testing.T) {
	expected := []byte("ABC")
	buf := bytes.NewBuffer(expected)
	closed := false
	rc := ose.NewReadCloser(buf, func(r io.Reader) error {
		closed = true
		return nil
	})
	nrc := ose.NopReadCloser(rc)
	nrc.Close()
	if closed {
		t.Fatal("must not be closed")
	}
}

func TestWriteCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	closed := false
	wc := ose.NewWriteCloser(buf, func(w io.Writer) error {
		closed = true
		return nil
	})
	expected := []byte("ABC")
	n, err := wc.Write(expected)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(expected) {
		t.Fatalf("invalid length: %d", n)
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("invalid value: %v (expected: %v)", buf.Bytes(), expected)
	}
	if closed {
		t.Fatal("must not be closed")
	}
	wc.Close()
	if !closed {
		t.Fatal("must be closed")
	}
}

func TestNopWriteCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	closed := false
	wc := ose.NewWriteCloser(buf, func(w io.Writer) error {
		closed = true
		return nil
	})
	nwc := ose.NopWriteCloser(wc)
	nwc.Close()
	if closed {
		t.Fatal("must not be closed")
	}
}

func TestExtendedReadCloser(t *testing.T) {
	expected := []byte("ABC")
	buf := bytes.NewBuffer(expected)
	closedBase := false
	rcBase := ose.NewReadCloser(buf, func(r io.Reader) error {
		closedBase = true
		return nil
	})
	closed := false
	rc := ose.ExtendReadCloser(rcBase, func(rc io.ReadCloser) error {
		closed = true
		rc.Close()
		return nil
	})
	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, bs) {
		t.Fatalf("invalid value: %v (expected: %v)", bs, expected)
	}
	if closed || closedBase {
		t.Fatal("must not be closed")
	}
	rc.Close()
	if !closed || !closedBase {
		t.Fatal("must be closed")
	}
}

func TestExtendedWriteCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	closedBase := false
	wcBase := ose.NewWriteCloser(buf, func(w io.Writer) error {
		closedBase = true
		return nil
	})
	closed := false
	wc := ose.ExtendWriteCloser(wcBase, func(wc io.WriteCloser) error {
		closed = true
		wc.Close()
		return nil
	})
	expected := []byte("ABC")
	n, err := wc.Write(expected)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(expected) {
		t.Fatalf("invalid length: %d", n)
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("invalid value: %v (expected: %v)", buf.Bytes(), expected)
	}
	if closed || closedBase {
		t.Fatal("must not be closed")
	}
	wc.Close()
	if !closed || !closedBase {
		t.Fatal("must be closed")
	}
}

func TestConditionalReadCloser(t *testing.T) {
	expected := []byte("ABC")
	buf := bytes.NewBuffer(expected)
	closed := false
	rc := ose.NewReadCloser(buf, func(r io.Reader) error {
		closed = true
		return nil
	})

	cc := ose.NewConditionalReadCloser(rc, false)
	if cc.Called {
		t.Fatal("must not be called")
	}
	cc.Close()
	if closed || cc.Closed {
		t.Fatal("must not be closed")
	}
	if !cc.Called {
		t.Fatal("must be called")
	}

	bs, err := ioutil.ReadAll(cc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, bs) {
		t.Fatalf("invalid value: %v (expected: %v)", bs, expected)
	}

	cc.Enabled = true
	cc.Close()
	if !closed || !cc.Closed {
		t.Fatal("must be closed")
	}
}
func TestConditionalWriteCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	closed := false
	wc := ose.NewWriteCloser(buf, func(w io.Writer) error {
		closed = true
		return nil
	})
	expected := []byte("ABC")

	cc := ose.NewConditionalWriteCloser(wc, false)
	if cc.Called {
		t.Fatal("must not be called")
	}
	cc.Close()
	if closed || cc.Closed {
		t.Fatal("must not be closed")
	}
	if !cc.Called {
		t.Fatal("must be called")
	}

	n, err := cc.Write(expected)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(expected) {
		t.Fatalf("invalid length: %d", n)
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf("invalid value: %v (expected: %v)", buf.Bytes(), expected)
	}

	cc.Enabled = true
	cc.Close()
	if !closed || !cc.Closed {
		t.Fatal("must be closed")
	}
}
