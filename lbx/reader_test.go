package lbx

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	smacker  = []byte("SMK4datadatadata")
	garbage  = []byte("garbage")
	filename = "filename"
)

func TestProcessHeader_Smacker(t *testing.T) {
	r := bytes.NewReader(smacker)
	d := decoder{r: r}

	filetype, err := d.processHeader()
	if err != nil {
		t.Error(err)
	}

	if filetype != filetypeSmacker {
		t.Error("wrong type")
	}
}

func TestProcessHeader_Garbage(t *testing.T) {
	r := bytes.NewReader(garbage)
	d := decoder{r: r}

	if filetype, _ := d.processHeader(); filetype != filetypeNone {
		t.Error("wrong type")
	}
}

func TestDecode_Smacker(t *testing.T) {
	r := bytes.NewReader(smacker)
	m, err := Decode(r, filename)

	if err != nil {
		t.Error(err)
	}

	if l := len(m); l != 1 {
		t.Error("expected 1 result, returned", l)
	}

	b, err := ioutil.ReadAll(m[filename+".1.smk"])
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(b, smacker) != 0 {
		t.Error("excepted ", smacker, ", returned ", m[filename+".1.smk"])
	}
}
