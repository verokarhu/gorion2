package lbx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type lbxHeader struct {
	NumEntries uint16
	Magic      uint16
	Reserved   uint16
	FileType   uint16
	Offsets    []uint32
}

func processFile(r io.Reader) (map[string][]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte)
	if string(b[0:3]) == "SMK" {
		m["1.smk"] = b
		return m, nil
	}

	lbxHeader := processLbxHeader(b)

	if lbxHeader.Magic == 65197 {
		for i := 0; i < int(lbxHeader.NumEntries); i++ {
			offset := lbxHeader.Offsets[i]
			endoffset := lbxHeader.Offsets[i+1]
			chunk := b[offset:endoffset]
			if string(chunk[0:4]) == "RIFF" {
				k := fmt.Sprintf("%d.wav", i+1)
				m[k] = chunk
			}
		}

		return m, nil
	}

	return nil, nil
}

func processLbxHeader(b []byte) lbxHeader {
	buf := bytes.NewReader(b)
	r := lbxHeader{}

	binary.Read(buf, binary.LittleEndian, &r.NumEntries)
	binary.Read(buf, binary.LittleEndian, &r.Magic)
	binary.Read(buf, binary.LittleEndian, &r.Reserved)
	binary.Read(buf, binary.LittleEndian, &r.FileType)

	r.Offsets = make([]uint32, r.NumEntries+1, r.NumEntries+1)
	for i := 0; i < int(r.NumEntries)+1; i++ {
		binary.Read(buf, binary.LittleEndian, &r.Offsets[i])
	}

	return r
}

// Decode reads the contents of an lbx file, decodes the contained files and returns readers that contain the converted files
func Decode(r io.Reader, name string) (map[string][]byte, error) {
	files, err := processFile(r)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte)
	for k, v := range files {
		m[name+"."+k] = v
	}

	return m, nil
}
