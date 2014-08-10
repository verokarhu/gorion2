package lbx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type header struct {
	NumEntries uint16
	Magic      uint16
	Reserved   uint16
	FileType   uint16
	Offsets    []uint32
}

type arrayHeader struct {
	NumElements uint16
	ElementSize uint16
}

type animationHeader struct {
	Width     uint16
	Height    uint16
	Unk1      uint16
	NumFrames uint16
	Unk2      uint32
	Offsets   []uint32
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

	h := processHeader(b)

	if h.Magic == 65197 {
		for i := 0; i < int(h.NumEntries); i++ {
			offset := h.Offsets[i]
			endoffset := h.Offsets[i+1]
			chunk := b[offset:endoffset]

			if string(chunk[0:4]) == "RIFF" {
				k := fmt.Sprintf("%d.wav", i+1)
				m[k] = chunk
			} else if ah, ok := processArrayHeader(chunk); ok {
				for j := 0; j < int(ah.NumElements); j++ {
					s := 4 + j*int(ah.ElementSize)
					e := s + int(ah.ElementSize)
					k := fmt.Sprintf("%d.%d.blob", i+1, j+1)
					m[k] = bytes.Trim(chunk[s:e], "\x00")
				}
			}
		}

		return m, nil
	}

	return nil, nil
}

func processHeader(b []byte) header {
	buf := bytes.NewReader(b)
	r := header{}

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

func processArrayHeader(b []byte) (arrayHeader, bool) {
	buf := bytes.NewReader(b)
	r := arrayHeader{}

	binary.Read(buf, binary.LittleEndian, &r.NumElements)
	binary.Read(buf, binary.LittleEndian, &r.ElementSize)

	if len(b) == int((r.NumElements*r.ElementSize + 4)) {
		return r, true
	}

	return r, false
}

func processAnimationHeader(b []byte) animationHeader {
	buf := bytes.NewReader(b)
	r := animationHeader{}

	binary.Read(buf, binary.LittleEndian, &r.Width)
	binary.Read(buf, binary.LittleEndian, &r.Height)
	binary.Read(buf, binary.LittleEndian, &r.Unk1)
	binary.Read(buf, binary.LittleEndian, &r.NumFrames)
	binary.Read(buf, binary.LittleEndian, &r.Unk2)

	r.Offsets = make([]uint32, r.NumFrames+1, r.NumFrames+1)
	for i := 0; i < int(r.NumFrames)+1; i++ {
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
