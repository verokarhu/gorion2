package lbx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

// Decode unarchives an lbx file and returns the bytes of the files within
func Decode(r io.Reader) ([][]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	h := processHeader(b)
	if h.Magic != 65197 {
		return nil, errors.New("not an lbx archive")
	}

	data := make([][]byte, h.NumEntries, h.NumEntries)
	for i := 0; i < int(h.NumEntries); i++ {
		s := h.Offsets[i]
		e := h.Offsets[i+1]
		data[i] = b[s:e]
	}

	return data, nil
}

func processHeader(b []byte) header {
	reader := bytes.NewReader(b)
	sh := subHeader{}

	binary.Read(reader, binary.LittleEndian, &sh)

	numentries := int(sh.NumEntries + 1)
	r := header{sh, make([]uint32, numentries, numentries)}
	for i := 0; i < numentries; i++ {
		binary.Read(reader, binary.LittleEndian, &r.Offsets[i])
	}

	return r
}

type subHeader struct {
	NumEntries uint16
	Magic      uint16
	Reserved   uint16
	FileType   uint16
}

type header struct {
	subHeader
	Offsets []uint32
}
