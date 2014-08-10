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
			if string(b[offset:offset+4]) == "RIFF" {
				k := fmt.Sprintf("%d.wav", i+1)

				if int(lbxHeader.NumEntries) == i+1 {
					m[k] = b[offset:]
				} else {
					endoffset := lbxHeader.Offsets[i+1]
					m[k] = b[offset:endoffset]
				}
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

	r.Offsets = make([]uint32, r.NumEntries, r.NumEntries)
	for i := 0; i < int(r.NumEntries); i++ {
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
