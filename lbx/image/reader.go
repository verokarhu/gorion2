package image

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

type subHeader struct {
	Width      uint16
	Height     uint16
	Marker     uint16
	NumFrames  uint16
	FrameDelay uint16
	Flags      uint16
}

type header struct {
	subHeader
	Offsets []uint32
}

// Decode converts an lbx image into a paletted image using the internal palette (if it exists)
func Decode(r io.ReadSeeker) (result []image.Paletted, err error) {
	sh := subHeader{}
	binary.Read(r, binary.LittleEndian, &sh)

	numentries := int(sh.NumFrames + 1)
	h := header{sh, make([]uint32, numentries, numentries)}
	for i := 0; i < numentries; i++ {
		binary.Read(r, binary.LittleEndian, &h.Offsets[i])
	}

	result = make([]image.Paletted, sh.NumFrames, sh.NumFrames)
	var p color.Palette

	if isInternalPalette(sh.Flags) {
		if p, err = DecodePalette(r); err != nil {
			return
		}
	}

	for i := 0; i < int(h.NumFrames); i++ {
		r.Seek(int64(h.Offsets[i]), 0)

		size := int(h.Width) * int(h.Height)
		img := image.Paletted{Stride: 1, Palette: p, Rect: image.Rect(0, 0, int(h.Width), int(h.Height)), Pix: make([]byte, size, size)}
		var numPix, yIndent, t uint16
		var xPos, yPos int

		binary.Read(r, binary.LittleEndian, &t)

		// frame always starts with 1
		if t == 1 {

			// frame Y indent, not using this for anything right now
			binary.Read(r, binary.LittleEndian, &t)

			for {
				binary.Read(r, binary.LittleEndian, &numPix)

				if numPix > 0 {
					binary.Read(r, binary.LittleEndian, &t)
					xPos += int(t)

					for j := 0; j < int(numPix); j++ {
						binary.Read(r, binary.LittleEndian, &img.Pix[yPos*int(h.Width)+xPos+j])
					}

					// if the number of pixels is uneven there is a filler byte at the end
					if numPix%2 != 0 {
						var b byte
						binary.Read(r, binary.LittleEndian, &b)
					}

					xPos += int(numPix)
				} else {
					binary.Read(r, binary.LittleEndian, &yIndent)

					// EOF
					if yIndent == 1000 {
						break
					}

					xPos = 0
					yPos += int(yIndent)
				}
			}
		} else {
			return nil, errors.New("frame didn't start with 1")
		}

		result[i] = img
	}

	return
}

type paletteHeader struct {
	Index     uint16
	Numcolors uint16
}

type paletteColor struct {
	A byte
	R byte
	G byte
	B byte
}

func DecodePalette(r io.Reader) (p color.Palette, err error) {
	ph := paletteHeader{}
	if err = binary.Read(r, binary.LittleEndian, &ph); err != nil {
		return
	}

	p = make(color.Palette, 256, 256)

	pc := paletteColor{}
	for i := 0; i < int(ph.Numcolors); i++ {
		if err = binary.Read(r, binary.LittleEndian, &pc); err != nil {
			return
		}

		p[int(ph.Index)+i] = color.NRGBA{4 * pc.R, 4 * pc.G, 4 * pc.B, 255 * pc.A}
	}

	return
}

func isInternalPalette(f uint16) bool {
	return f >= 4096
}
