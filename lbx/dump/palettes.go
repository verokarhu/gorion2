package dumper

import (
	"bytes"
	"fmt"

	li "github.com/verokarhu/gorion2/lbx/image"
)

type palettePair struct {
	Filename string
	Index    int
}

var PaletteFiles = []palettePair{
	{"fonts", 1},
	{"fonts", 2},
	{"fonts", 3},
	{"fonts", 4},
	{"fonts", 5},
	{"fonts", 6},
	{"fonts", 7},
	{"fonts", 8},
	{"fonts", 9},
	{"fonts", 10},
	{"fonts", 11},
	{"fonts", 12},
	{"fonts", 13},
	{"ifonts", 1},
	{"ifonts", 2},
	{"ifonts", 3},
	{"ifonts", 4},
}

func loadExternalPalettes(dirname string) (pals map[string]li.Palette, err error) {
	pals = make(map[string]li.Palette)

	for _, file := range PaletteFiles {
		data, err := decodeFile(dirname + "/" + file.Filename + ".lbx")
		if err != nil {
			return nil, err
		}

		pals[fmt.Sprintf("%s%d", file.Filename, file.Index)] = li.ConvertPalette(bytes.NewReader(data[file.Index]), 0, 256)
	}

	return
}
