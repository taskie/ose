package osplus

import (
	"bufio"
	"io"
	"os"
)

func newBufferedReader(rc io.ReadCloser) io.ReadCloser {
	br := bufio.NewReader(rc)
	return NewReadCloser(br, func(_ io.Reader) error {
		return rc.Close()
	})
}

func newBufferedWriter(wc io.WriteCloser) io.WriteCloser {
	bw := bufio.NewWriter(wc)
	return NewWriteCloser(bw, func(_ io.Writer) error {
		err1 := bw.Flush()
		err2 := wc.Close()
		if err1 != nil {
			return err1
		}
		return err2
	})
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

func (opener *Opener) shouldFallback(name string) bool {
	return name == "" || (!opener.TreatHyphenAsFileName && name == "-")
}

func (opener *Opener) openFile(name string, ff func(name string) (*os.File, error)) (io.ReadCloser, error) {
	if opener.shouldFallback(name) {
		if opener.Unbuffered {
			return NopReadCloser(opener.FallbackReader), nil
		}
		return NopReadCloser(bufio.NewReader(opener.FallbackReader)), nil
	}
	f, err := ff(name)
	if err != nil {
		return nil, err
	}
	if opener.Unbuffered {
		return f, nil
	}
	return newBufferedReader(f), nil

}

func (opener *Opener) OpenFile(name string, flag int, perm os.FileMode) (io.ReadCloser, error) {
	return opener.openFile(name, func(name string) (*os.File, error) { return os.OpenFile(name, flag, perm) })
}

func (opener *Opener) Open(name string) (io.ReadCloser, error) {
	return opener.openFile(name, os.Open)
}

func (opener *Opener) createFile(name string, ff func(name string) (*os.File, error)) (io.WriteCloser, error) {
	if opener.shouldFallback(name) {
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
	return newBufferedWriter(f), nil
}

func (opener *Opener) CreateFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return opener.createFile(name, func(name string) (*os.File, error) { return os.OpenFile(name, flag, perm) })
}

func (opener *Opener) Create(name string) (io.WriteCloser, error) {
	return opener.createFile(name, os.Create)
}

var nops = func(name string) {}

var nopb = func(ok bool) {}

func (opener *Opener) CreateTempFile(tempFileDir, tempFilePattern string) (wc io.WriteCloser, commit func(name string), err error) {
	tmp, err := CreateTempFile(tempFileDir, tempFilePattern)
	if err != nil {
		return nil, nops, err
	}
	wc = tmp
	if !opener.Unbuffered {
		wc = newBufferedWriter(tmp)
	}
	return wc, func(name string) {
		tmp.Destination = name
	}, err
}

func (opener *Opener) CreateTempFileWithDestination(name, tempFileDir, tempFilePattern string) (wc io.WriteCloser, commit func(ok bool), err error) {
	if opener.shouldFallback(name) {
		if opener.Unbuffered {
			return NopWriteCloser(opener.FallbackWriter), nopb, nil
		}
		bw := bufio.NewWriter(opener.FallbackWriter)
		return NewWriteCloser(bw, func(_ io.Writer) error {
			return bw.Flush()
		}), nopb, nil
	}
	wc, commits, err := opener.CreateTempFile(tempFileDir, tempFilePattern)
	return wc, func(ok bool) {
		if ok {
			commits(name)
		} else {
			commits("")
		}
	}, err
}
