package osplus

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadCloser(t *testing.T) {
	expected := []byte("ABC")
	buf := bytes.NewBuffer(expected)
	closed := false
	rc := NewReadCloser(buf, func(r io.Reader) error {
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

func TestWriteCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	closed := false
	wc := NewWriteCloser(buf, func(w io.Writer) error {
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

func TestOpener(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}()

	fooPath := filepath.Join(tmp, "foo")
	o := NewOpener()
	wc, err := o.Create(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	wc.Close()

	rc, err := o.Open(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	rc.Close()
	if string(bs) != "foo" {
		t.Fatalf("invalid content: %v", bs)
	}
}

func TestOpenerViaTempFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}()

	fooPath := filepath.Join(tmp, "foo")
	o := NewOpener()
	wc, _, err := o.CreateViaTempFile(fooPath, "", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
	wc.Close()
	bs, err := ioutil.ReadFile(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "foo" {
		t.Fatalf("invalid content: %v", bs)
	}
}

func TestCancelOpenerViaTempFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}()

	fooPath := filepath.Join(tmp, "foo")
	o := NewOpener()
	wc, cancel, err := o.CreateViaTempFile(fooPath, "", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
	cancel()
	wc.Close()
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
}
