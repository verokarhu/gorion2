package lbx

import (
	"bytes"
	"testing"
)

var (
	smacker  = []byte("SMK4datadatadata")
	garbage  = []byte("garbage")
	audio    = []byte{2, 0, 173, 254, 0, 0, 0, 0, 20, 0, 0, 0, 24, 0, 0, 0, 28, 0, 0, 0, 0x52, 0x49, 0x46, 0x46, 0x52, 0x49, 0x46, 0x46}
	array    = []byte{1, 0, 173, 254, 0, 0, 0, 0, 16, 0, 0, 0, 26, 0, 0, 0, 3, 0, 2, 0, 97, 97, 98, 98, 99, 99}
	filename = "filename"
)

func TestProcessFile_Smacker(t *testing.T) {
	r := bytes.NewReader(smacker)

	m, err := processFile(r)
	if err != nil {
		t.Error(err)
	}

	if l := len(m); l != 1 {
		t.Error("excepted 1, returned ", l)
	}

	if bytes.Compare(m["1.smk"], smacker) != 0 {
		t.Error("excepted ", smacker, ", returned ", m["1.smk"])
	}
}

func TestProcessFile_Garbage(t *testing.T) {
	r := bytes.NewReader(garbage)

	m, err := processFile(r)
	if err != nil {
		t.Error(err)
	}

	if l := len(m); l != 0 {
		t.Error("excepted 0, returned ", l)
	}
}

func TestProcessFile_Audio(t *testing.T) {
	r := bytes.NewReader(audio)

	m, err := processFile(r)
	if err != nil {
		t.Error(err)
	}

	if l := len(m); l != 2 {
		t.Error("excepted 2, returned ", l)
	}

	if r := "1.wav"; m[r] == nil {
		t.Error("excepted ", r, " to exist")
	}

	if r := "2.wav"; m[r] == nil {
		t.Error("excepted ", r, " to exist")
	}

	if r := m["1.wav"]; bytes.Compare(r, audio[20:24]) != 0 {
		t.Error("excepted ", audio[20:24], ", returned ", r)
	}

	if r := m["2.wav"]; bytes.Compare(r, audio[24:28]) != 0 {
		t.Error("excepted ", audio[24:28], ", returned ", r)
	}
}

func TestProcessFile_Array(t *testing.T) {
	r := bytes.NewReader(array)

	m, err := processFile(r)
	if err != nil {
		t.Error(err)
	}

	if l := len(m); l != 3 {
		t.Error("excepted 3, returned ", l)
	}

	if r := "1.1.blob"; m[r] == nil {
		t.Error("excepted ", r, " to exist")
	}

	if r := "1.2.blob"; m[r] == nil {
		t.Error("excepted ", r, " to exist")
	}

	if r := "1.3.blob"; m[r] == nil {
		t.Error("excepted ", r, " to exist")
	}

	if r := m["1.1.blob"]; bytes.Compare(r, array[20:22]) != 0 {
		t.Error("excepted ", array[20:22], ", returned ", r)
	}

	if r := m["1.2.blob"]; bytes.Compare(r, array[22:24]) != 0 {
		t.Error("excepted ", array[22:24], ", returned ", r)
	}

	if r := m["1.3.blob"]; bytes.Compare(r, array[24:]) != 0 {
		t.Error("excepted ", array[24:], ", returned ", r)
	}
}

func TestProcessHeader(t *testing.T) {
	h := processHeader(audio)

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
		t.Error("expected 3 results, returned", l)
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

func TestProcessArrayHeader_Array(t *testing.T) {
	h, ok := processArrayHeader(array[16:])

	if r := h.NumElements; r != 3 {
		t.Error("expected 3, returned", r)
	}

	if r := h.ElementSize; r != 2 {
		t.Error("expected 2, returned", r)
	}

	if !ok {
		t.Error("expected true, returned false")
	}
}

func TestProcessArrayHeader_Garbage(t *testing.T) {
	_, ok := processArrayHeader(garbage)

	if ok {
		t.Error("expected false, returned true")
	}
}

func TestDecode(t *testing.T) {
	r := bytes.NewReader(smacker)
	m, err := Decode(r, filename)

	if err != nil {
		t.Error(err)
	}

	for k, v := range m {
		if k != filename+".1.smk" {
			t.Error("excepted ", filename, ".1.smk, returned ", k)
		}

		if bytes.Compare(v, smacker) != 0 {
			t.Error("excepted ", smacker, ", returned ", v)
		}
	}
}
