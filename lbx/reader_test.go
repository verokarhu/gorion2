package lbx

import (
	"bytes"
	"testing"
)

var (
	smacker  = []byte("SMK4datadatadata")
	garbage  = []byte("garbage")
	audio    = []byte{2, 0, 173, 254, 0, 0, 0, 0, 20, 0, 0, 0, 24, 0, 0, 0, 28, 0, 0, 0, 0x52, 0x49, 0x46, 0x46, 0x52, 0x49, 0x46, 0x46}
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

func TestProcessLbxHeader(t *testing.T) {
	lbx := processLbxHeader(audio)

	if r := lbx.NumEntries; r != 2 {
		t.Error("expected 2, returned", r)
	}

	if r := lbx.Magic; r != 65197 {
		t.Error("expected 65197, returned", r)
	}

	if r := lbx.Reserved; r != 0 {
		t.Error("expected 0, returned", r)
	}

	if r := lbx.FileType; r != 0 {
		t.Error("expected 0, returned", r)
	}

	if l := len(lbx.Offsets); l != 3 {
		t.Error("expected 3 results, returned", l)
	}

	if r := lbx.Offsets[0]; r != 20 {
		t.Error("expected 20, returned", r)
	}

	if r := lbx.Offsets[1]; r != 24 {
		t.Error("expected 24, returned", r)
	}

	if r := lbx.Offsets[2]; r != 28 {
		t.Error("expected 28, returned", r)
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
