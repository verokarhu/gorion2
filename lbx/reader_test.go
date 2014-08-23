package lbx

import (
	"bytes"
	"testing"
)

var (
	garbage = []byte("garbage")
	lbx     = []byte{
		2, 0, 173, 254, 0, 0, 0, 0,
		20, 0, 0, 0, 24, 0, 0, 0,
		28, 0, 0, 0, 0x52, 0x49, 0x46, 0x46,
		0x52, 0x49, 0x46, 0x46}
)

func Test_Decode_Garbage(t *testing.T) {
	f := bytes.NewReader(garbage)

	m, err := Decode(f)
	if err == nil {
		t.Fatal("expected an error")
	}

	if l := len(m); l != 0 {
		t.Fatal("excepted 0, returned ", l)
	}
}

func Test_Decode_LBX(t *testing.T) {
	f := bytes.NewReader(lbx)

	decoded, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(decoded); l != 2 {
		t.Fatal("excepted 2, returned ", l)
	}

	if c := [][]byte{decoded[0], lbx[20:24]}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := [][]byte{decoded[1], lbx[24:28]}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}

func Test_ProcessHeader(t *testing.T) {
	h := processHeader(lbx)

	if r := h.NumEntries; r != 2 {
		t.Error("expected 2, returned", r)
	}

	if r := h.Magic; r != 65197 {
		t.Error("expected 65197, returned", r)
	}

	if r := h.Reserved; r != 0 {
		t.Error("expected 0, returned", r)
	}

	if r := h.FileType; r != 0 {
		t.Error("expected 0, returned", r)
	}

	if l := len(h.Offsets); l != 3 {
		t.Fatal("expected 3 results, returned", l)
	}

	if r := h.Offsets[0]; r != 20 {
		t.Error("expected 20, returned", r)
	}

	if r := h.Offsets[1]; r != 24 {
		t.Error("expected 24, returned", r)
	}

	if r := h.Offsets[2]; r != 28 {
		t.Error("expected 28, returned", r)
	}
}
