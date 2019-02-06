package osplus

import (
	"io"
)

// ComposedReadCloser implements io.ReadCloser.
// It has an underlying Reader and a function that is called when Close method is called.
type ComposedReadCloser struct {
	Reader    io.Reader
	CloseFunc func(r io.Reader) error
}

// NewReadCloser composes a ReadCloser with an underlying Reader and a function that is called when Close method is called.
func NewReadCloser(r io.Reader, closeFunc func(r io.Reader) error) *ComposedReadCloser {
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

func nopReadCloseFunc(_ io.Reader) error {
	return nil
}

// NopReadCloser composes a ReadCloser with an underlying Reader. It does nothing when Close method is called.
func NopReadCloser(r io.Reader) *ComposedReadCloser {
	return NewReadCloser(r, nopReadCloseFunc)
}

// ComposedWriteCloser implements io.WriteCloser.
// It has an underlying Writer and a function that is called when Close method is called.
type ComposedWriteCloser struct {
	Writer    io.Writer
	CloseFunc func(w io.Writer) error
}

// NewWriteCloser composes a WriteCloser with an underlying Writer and a function that is called when Close method is called.
func NewWriteCloser(w io.Writer, closeFunc func(w io.Writer) error) *ComposedWriteCloser {
	return &ComposedWriteCloser{
		Writer:    w,
		CloseFunc: closeFunc,
	}
}

func (cwc *ComposedWriteCloser) Write(p []byte) (int, error) {
	return cwc.Writer.Write(p)
}

func (cwc *ComposedWriteCloser) Close() error {
	return cwc.CloseFunc(cwc.Writer)
}

func nopWriteCloseFunc(_ io.Writer) error {
	return nil
}

// NopWriteCloser composes a WriteCloser with an underlying Writer. It does nothing when Close method is called.
func NopWriteCloser(w io.Writer) *ComposedWriteCloser {
	return NewWriteCloser(w, nopWriteCloseFunc)
}

// ExtendedReadCloser implements io.ReadCloser.
// It has an underlying Reader and a function that is called when Close method is called.
type ExtendedReadCloser struct {
	ReadCloser io.ReadCloser
	CloseFunc  func(rc io.ReadCloser) error
}

// ExtendReadCloser composes a ReadCloser with an underlying Reader and a function that is called when Close method is called.
func ExtendReadCloser(rc io.ReadCloser, closeFunc func(rc io.ReadCloser) error) *ExtendedReadCloser {
	return &ExtendedReadCloser{
		ReadCloser: rc,
		CloseFunc:  closeFunc,
	}
}

func (erc *ExtendedReadCloser) Read(p []byte) (int, error) {
	return erc.ReadCloser.Read(p)
}

func (erc *ExtendedReadCloser) Close() error {
	return erc.CloseFunc(erc.ReadCloser)
}

// ExtendedWriteCloser implements io.WriteCloser.
// It has an underlying Writer and a function that is called when Close method is called.
type ExtendedWriteCloser struct {
	WriteCloser io.WriteCloser
	CloseFunc   func(w io.WriteCloser) error
}

// ExtendWriteCloser composes a WriteCloser with an underlying Writer and a function that is called when Close method is called.
func ExtendWriteCloser(wc io.WriteCloser, closeFunc func(wc io.WriteCloser) error) *ExtendedWriteCloser {
	return &ExtendedWriteCloser{
		WriteCloser: wc,
		CloseFunc:   closeFunc,
	}
}

func (ewc *ExtendedWriteCloser) Write(p []byte) (int, error) {
	return ewc.WriteCloser.Write(p)
}

func (ewc *ExtendedWriteCloser) Close() error {
	return ewc.CloseFunc(ewc.WriteCloser)
}

// ConditionalCloser implements io.Closer.
// When Close method is called, it close the underlying closer only if Enabled is true.
type ConditionalCloser struct {
	Closer  io.Closer
	Enabled bool
	Closed  bool
	Called  bool
}

// NewConditionalCloser creates new ConditionalCloser.
func NewConditionalCloser(c io.Closer, enabled bool) *ConditionalCloser {
	return &ConditionalCloser{
		Closer:  c,
		Enabled: enabled,
	}
}

func (cc *ConditionalCloser) Close() error {
	cc.Called = true
	if cc.Enabled {
		cc.Closed = true
		return cc.Closer.Close()
	}
	return nil
}

type ConditionalReadCloser struct {
	ConditionalCloser
	ReadCloser io.ReadCloser
}

// NewConditionalReadCloser creates new ConditionalReadCloser.
func NewConditionalReadCloser(c io.ReadCloser, enabled bool) *ConditionalReadCloser {
	return &ConditionalReadCloser{
		ConditionalCloser: ConditionalCloser{Closer: c, Enabled: enabled},
		ReadCloser:        c,
	}
}

func (cc *ConditionalReadCloser) Read(p []byte) (int, error) {
	return cc.ReadCloser.Read(p)
}

func (cc *ConditionalReadCloser) Close() error {
	return cc.ConditionalCloser.Close()
}

func (cc *ConditionalReadCloser) SetReadCloser(rc io.ReadCloser) {
	cc.Closer = rc
	cc.ReadCloser = rc
}

type ConditionalWriteCloser struct {
	ConditionalCloser
	WriteCloser io.WriteCloser
}

// NewConditionalWriteCloser creates new ConditionalWriteCloser.
func NewConditionalWriteCloser(c io.WriteCloser, enabled bool) *ConditionalWriteCloser {
	return &ConditionalWriteCloser{
		ConditionalCloser: ConditionalCloser{Closer: c, Enabled: enabled},
		WriteCloser:       c,
	}
}

func (cc *ConditionalWriteCloser) Write(p []byte) (int, error) {
	return cc.WriteCloser.Write(p)
}

func (cc *ConditionalWriteCloser) Close() error {
	return cc.ConditionalCloser.Close()
}

func (cc *ConditionalWriteCloser) SetWriteCloser(wc io.WriteCloser) {
	cc.Closer = wc
	cc.WriteCloser = wc
}
