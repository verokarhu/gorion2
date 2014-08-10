package lbx

import (
	"io"
	"io/ioutil"
)

type decoder struct {
	r   io.Reader
	tmp [4096]byte
}

const (
	filetypeNone = iota
	filetypeSmacker
	filetypeImage
	filetypeSound
)

func (d *decoder) processHeader() (int, error) {
	if _, err := io.ReadFull(d.r, d.tmp[:12]); err != nil {
		return filetypeNone, err
	}

	if string(d.tmp[0:3]) == "SMK" {
		return filetypeSmacker, nil
	}

	return filetypeNone, nil
}

// Decode reads the contents of an lbx file, decodes the contained files and returns readers that contain the converted files
func Decode(r io.Reader, name string) (map[string][]byte, error) {
	d := &decoder{r: r}
	t, err := d.processHeader()

	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte)
	switch t {
	case filetypeSmacker:
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		m[name+".1.smk"] = append(d.tmp[:12], b[:]...)
	}

	return m, nil
}
