package kiwi

import (
	"bytes"
	"testing"
)

func TestWriteByte(t *testing.T) {
	testWriteByte(t, 0, []byte{0x00})
	testWriteByte(t, 1, []byte{0x01})
	testWriteByte(t, 254, []byte{0xFE})
	testWriteByte(t, 255, []byte{0xFF})
}

func testWriteByte(t *testing.T, v byte, b []byte) {
	buf := new(bytes.Buffer)
	if err := NewWriter(buf).WriteByte(v); err != nil {
		t.Error(err)
	}
	checkResults(t, buf.Bytes(), b)
}

func TestWriteUint(t *testing.T) {
	testWriteUint(t, 0, []byte{0x00})
	testWriteUint(t, 1, []byte{0x01})
	testWriteUint(t, 2, []byte{0x02})
	testWriteUint(t, 127, []byte{0x7F})
	testWriteUint(t, 128, []byte{0x80, 0x01})
	testWriteUint(t, 129, []byte{0x81, 0x01})
	testWriteUint(t, 255, []byte{0xFF, 0x01})
	testWriteUint(t, 256, []byte{0x80, 0x02})
	testWriteUint(t, 16383, []byte{0xFF, 0x7F})
	testWriteUint(t, 16384, []byte{0x80, 0x80, 0x01})
	testWriteUint(t, 16385, []byte{0x81, 0x80, 0x01})
	testWriteUint(t, 32767, []byte{0xFF, 0xFF, 0x01})
	testWriteUint(t, 32768, []byte{0x80, 0x80, 0x02})
	testWriteUint(t, 2097151, []byte{0xFF, 0xFF, 0x7F})
	testWriteUint(t, 2097152, []byte{0x80, 0x80, 0x80, 0x01})
	testWriteUint(t, 2097153, []byte{0x81, 0x80, 0x80, 0x01})
	testWriteUint(t, 4194303, []byte{0xFF, 0xFF, 0xFF, 0x01})
	testWriteUint(t, 4194304, []byte{0x80, 0x80, 0x80, 0x02})
	testWriteUint(t, 268435455, []byte{0xFF, 0xFF, 0xFF, 0x7F})
	testWriteUint(t, 268435456, []byte{0x80, 0x80, 0x80, 0x80, 0x01})
	testWriteUint(t, 268435457, []byte{0x81, 0x80, 0x80, 0x80, 0x01})
	testWriteUint(t, 536870911, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01})
	testWriteUint(t, 536870912, []byte{0x80, 0x80, 0x80, 0x80, 0x02})
	testWriteUint(t, 4294967294, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0x0F})
	testWriteUint(t, 4294967295, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F})
}

func testWriteUint(t *testing.T, v uint32, b []byte) {
	buf := new(bytes.Buffer)
	if err := NewWriter(buf).WriteUint(v); err != nil {
		t.Error(err)
	}
	checkResults(t, buf.Bytes(), b)
}

func testWriteInt(t *testing.T, v int32, b []byte) {
	buf := new(bytes.Buffer)
	if err := NewWriter(buf).WriteInt(v); err != nil {
		t.Error(err)
	}
	checkResults(t, buf.Bytes(), b)
}

func TestWriteInt(t *testing.T) {
	testWriteInt(t, 0, []byte{0x00})
	testWriteInt(t, 1, []byte{0x01})
	testWriteInt(t, 2, []byte{0x02})
	testWriteInt(t, 127, []byte{0x7F})
	testWriteInt(t, 128, []byte{0x80, 0x01})
	testWriteInt(t, 129, []byte{0x81, 0x01})
	testWriteInt(t, 255, []byte{0xFF, 0x01})
	testWriteInt(t, 256, []byte{0x80, 0x02})
	testWriteInt(t, 16383, []byte{0xFF, 0x7F})
	testWriteInt(t, 16384, []byte{0x80, 0x80, 0x01})
	testWriteInt(t, 16385, []byte{0x81, 0x80, 0x01})
	testWriteInt(t, 32767, []byte{0xFF, 0xFF, 0x01})
	testWriteInt(t, 32768, []byte{0x80, 0x80, 0x02})
	testWriteInt(t, 2097151, []byte{0xFF, 0xFF, 0x7F})
	testWriteInt(t, 2097152, []byte{0x80, 0x80, 0x80, 0x01})
	testWriteInt(t, 2097153, []byte{0x81, 0x80, 0x80, 0x01})
	testWriteInt(t, 4194303, []byte{0xFF, 0xFF, 0xFF, 0x01})
	testWriteInt(t, 4194304, []byte{0x80, 0x80, 0x80, 0x02})
	testWriteInt(t, 268435455, []byte{0xFF, 0xFF, 0xFF, 0x7F})
	testWriteInt(t, 268435456, []byte{0x80, 0x80, 0x80, 0x80, 0x01})
	testWriteInt(t, 268435457, []byte{0x81, 0x80, 0x80, 0x80, 0x01})
	testWriteInt(t, 536870911, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01})
	testWriteInt(t, 536870912, []byte{0x80, 0x80, 0x80, 0x80, 0x02})
	testWriteInt(t, 4294967294, []byte{0xFE, 0xFF, 0xFF, 0xFF, 0x0F})
	testWriteInt(t, 4294967295, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F})
}

func checkResults(t *testing.T, result, control []byte) {
	if len(result) != len(control) {
		t.Errorf("different length:\n"+
			"\tcontrol: [%d]byte{% X}\n"+
			"\tresult : [%d]byte{% X}\n",
			len(control), control, len(result), result)
		return
	}
	for i := 0; i < len(result); i++ {
		if result[i] != control[i] {
			t.Errorf("different in %d element:\n"+
				"\tcontrol: [%d]byte{% X}\n"+
				"\tresult : [%d]byte{% X}\n",
				i, len(control), control, len(result), result)
		}
	}
}
