package osplus

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

type ComposedReadCloser struct {
	Reader    io.Reader
	CloseFunc func(r io.Reader) error
}

func NewReadCloser(r io.Reader, closeFunc func(r io.Reader) error) io.ReadCloser {
	return &ComposedReadCloser{
		Reader:    r,
		CloseFunc: closeFunc,
	}
}

func (crc *ComposedReadCloser) Read(p []byte) (int, error) {
	return crc.Reader.Read(p)
}

func (crc *ComposedReadCloser) Close() error {
	return crc.CloseFunc(crc.Reader)
}

type ComposedWriteCloser struct {
	Writer    io.Writer
	CloseFunc func(w io.Writer) error
}

func NewWriteCloser(w io.Writer, closeFunc func(w io.Writer) error) io.WriteCloser {
	return &ComposedWriteCloser{
		Writer:    w,
		CloseFunc: closeFunc,
	}
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return NewWriteCloser(w, func(_ io.Writer) error { return nil })
}

func (cwc *ComposedWriteCloser) Write(p []byte) (int, error) {
	return cwc.Writer.Write(p)
}

func (cwc *ComposedWriteCloser) Close() error {
	return cwc.CloseFunc(cwc.Writer)
}

type Opener struct {
	FallbackReader        io.Reader
	FallbackWriter        io.Writer
	TreatHyphenAsFileName bool
	Unbuffered            bool
	MoveOptions           *MoveOptions
}

func NewOpener() *Opener {
	return &Opener{
		FallbackReader: os.Stdin,
		FallbackWriter: os.Stdout,
	}
}

func (opener *Opener) openFile(name string, ff func(name string) (*os.File, error)) (io.ReadCloser, error) {
	if name == "" || (!opener.TreatHyphenAsFileName && name == "-") {
		if opener.Unbuffered {
			return ioutil.NopCloser(opener.FallbackReader), nil
		}
		return ioutil.NopCloser(bufio.NewReader(opener.FallbackReader)), nil
	}
	f, err := ff(name)
	if err != nil {
		return nil, err
	}
	if opener.Unbuffered {
		return f, nil
	}
	return NewReadCloser(bufio.NewReader(f), func(_ io.Reader) error {
		return f.Close()
	}), nil
}

func (opener *Opener) OpenFile(name string, flag int, perm os.FileMode) (io.ReadCloser, error) {
	return opener.openFile(name, func(name string) (*os.File, error) { return os.OpenFile(name, flag, perm) })
}

func (opener *Opener) Open(name string) (io.ReadCloser, error) {
	return opener.openFile(name, os.Open)
}

func (opener *Opener) createFile(name string, ff func(name string) (*os.File, error)) (io.WriteCloser, error) {
	if name == "" || (!opener.TreatHyphenAsFileName && name == "-") {
		if opener.Unbuffered {
			return NopWriteCloser(opener.FallbackWriter), nil
		}
		bw := bufio.NewWriter(opener.FallbackWriter)
		return NewWriteCloser(bw, func(_ io.Writer) error {
			return bw.Flush()
		}), nil
	}
	f, err := ff(name)
	if err != nil {
		return nil, err
	}
	if opener.Unbuffered {
		return f, nil
	}
	bw := bufio.NewWriter(f)
	return NewWriteCloser(bw, func(_ io.Writer) error {
		err := bw.Flush()
		if err != nil {
			return err
		}
		return f.Close()
	}), nil
}

func (opener *Opener) CreateFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return opener.createFile(name, func(name string) (*os.File, error) { return os.OpenFile(name, flag, perm) })
}

func (opener *Opener) Create(name string) (io.WriteCloser, error) {
	return opener.createFile(name, os.Create)
}

var nop = func() {}

func (opener *Opener) CreateViaTempFile(name, tempFileDir, tempFilePattern string) (io.WriteCloser, func(), error) {
	if name == "" || (!opener.TreatHyphenAsFileName && name == "-") {
		if opener.Unbuffered {
			return NopWriteCloser(opener.FallbackWriter), nop, nil
		}
		bw := bufio.NewWriter(opener.FallbackWriter)
		return NewWriteCloser(bw, func(_ io.Writer) error {
			return bw.Flush()
		}), nop, nil
	}
	f, err := ioutil.TempFile(tempFileDir, tempFilePattern)
	if err != nil {
		return nil, nop, err
	}
	var wc io.WriteCloser = f
	if !opener.Unbuffered {
		bw := bufio.NewWriter(f)
		wc = NewWriteCloser(bw, func(_ io.Writer) error {
			err := bw.Flush()
			if err != nil {
				return err
			}
			return f.Close()
		})
	}
	canceled := false
	return NewWriteCloser(wc, func(_ io.Writer) error {
		err := wc.Close()
		if err != nil {
			return err
		}
		if canceled {
			os.Remove(f.Name())
			return nil
		}
		err = MoveFile(f.Name(), name, opener.MoveOptions)
		if err != nil {
			os.Remove(f.Name())
			return err
		}
		return nil
	}), func() { canceled = true }, err
}
