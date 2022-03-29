package kiwi

import (
	"encoding/binary"
	"io"
	"math"
)

type Writer struct {
	w    io.Writer
	buf  []byte
	size int
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:   w,
		buf: make([]byte, 5),
	}
}

func (w *Writer) WriteByte(v byte) error {
	w.buf[0], w.size = v, 1
	return w.Flush()
}
func (w *Writer) WriteUint(v uint32) error {
	w.size = binary.PutUvarint(w.buf, uint64(v))
	return w.Flush()
}
func (w *Writer) WriteInt(v int32) error {
	w.size = binary.PutVarint(w.buf, int64(v))
	return w.Flush()
}
func (w *Writer) WriteFloat(v float32) error {
	bits := math.Float32bits(v)
	bits = (bits >> 23) | (bits << 9)
	if (bits & 255) == 0 {
		return w.WriteByte(0)
	}
	binary.LittleEndian.PutUint32(w.buf, bits)
	w.size = 4
	return w.Flush()
}
func (w *Writer) WriteByteArray(v []byte) error {
	if err := w.WriteUint(uint32(len(v))); err != nil {
		return err
	}
	return writeFull(w.w, v)
}
func (w *Writer) WriteString(v string) error {
	return writeFull(w.w, append([]byte(v), 0))
}
func (w *Writer) Flush() error {
	err := writeFull(w.w, w.buf[:w.size])
	w.size = 0
	return err
}
func (w *Writer) WriteBool(v bool) error {
	if v {
		return w.WriteByte(1)
	}
	return w.WriteByte(0)
}
func writeFull(w io.Writer, buf []byte) (err error) {
	for n, nn := 0, 0; n < len(buf) && err == nil; nn = 0 {
		nn, err = w.Write(buf[n:])
		n += nn
	}
	return err
}
