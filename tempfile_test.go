package osplus

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	tf, err := CreateTempFile("", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("ABC")
	_, err = tf.Write(expected)
	if err != nil {
		t.Fatal(err)
	}
	err = tf.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(tf.File.Name())
	if err == nil {
		t.Fatalf("%s must not exist", tf.File.Name())
	}
}

func TestCreateTempFileWithDestination(t *testing.T) {
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
	tf, err := CreateTempFileWithDestination(fooPath, "", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("ABC")
	_, err = tf.Write(expected)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not exist", fooPath)
	}
	err = tf.Close()
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadFile(fooPath)
	if err != nil {
		t.Fatalf("%s must exist: %s", fooPath, err.Error())
	}
	if !bytes.Equal(bs, expected) {
		t.Fatalf("invalid value: %v (expected: %v)", bs, expected)
	}
	_, err = os.Stat(tf.File.Name())
	if err == nil {
		t.Fatalf("%s must not exist", tf.File.Name())
	}
}
