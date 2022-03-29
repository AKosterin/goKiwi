package kiwi

import (
	"encoding/binary"
	"io"
	"math"
	"strings"
	"unicode/utf8"
)

type Reader struct {
	r   io.Reader
	buf []byte
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:   r,
		buf: make([]byte, 5),
	}
}

func (r Reader) ReadByte() (byte, error) {
	_, err := io.ReadFull(r.r, r.buf[:1])
	return r.buf[0], err
}
func (r Reader) ReadUint() (uint32, error) {
	v, err := binary.ReadUvarint(r)
	return uint32(v), err
}
func (r Reader) ReadInt() (int32, error) {
	v, err := binary.ReadVarint(r)
	return int32(v), err
}
func (r Reader) ReadFloat() (float32, error) {
	if _, err := io.ReadFull(r.r, r.buf[:1]); err != nil || r.buf[0] == 0 {
		return 0, err
	}
	if _, err := io.ReadFull(r.r, r.buf[1:4]); err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(r.buf[:4])
	bits = (bits << 23) | (bits >> 9)
	return math.Float32frombits(bits), nil
}
func (r Reader) ReadByteArray() ([]byte, error) {
	l, err := r.ReadUint()
	if err != nil {
		return nil, err
	}
	bb := make([]byte, int(l))
	_, err = io.ReadFull(r.r, bb)
	return bb, err
}
func (r Reader) ReadRune() (rn rune, size int, err error) {
	if _, err = io.ReadFull(r.r, r.buf[:1]); err != nil || r.buf[0] < 0xC0 {
		return rune(r.buf[0]), 1, err
	}
	switch {
	case r.buf[0] < 0xE0:
		size = 2
	case r.buf[0] < 0xF0:
		size = 3
	default:
		size = 4
	}
	if _, err = io.ReadFull(r.r, r.buf[1:size]); err != nil {
		return 0, size, err
	}
	rn, size = utf8.DecodeRune(r.buf)
	return rn, size, err
}
func (r Reader) ReadString() (string, error) {
	b := strings.Builder{}
	for {
		if rn, _, err := r.ReadRune(); err != nil || rn == 0 {
			return b.String(), err
		} else {
			b.WriteRune(rn)
		}
	}
}
func (r Reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	return b > 0, err
}
