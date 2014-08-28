package image

import (
	"encoding/binary"
	"errors"
	"fmt"
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

// image header flags
const (
	NoCompression = 256 << iota
	Unk
	FillBackground
	FunctionalColor
	InternalPalette
	Junction
)

// Decode converts an lbx image into a Animation using the internal palette (if it exists)
func Decode(r io.ReadSeeker) (anim Animation, err error) {
	sh := subHeader{}
	binary.Read(r, binary.LittleEndian, &sh)

	numentries := int(sh.NumFrames + 1)
	h := header{sh, make([]uint32, numentries)}
	for i := 0; i < numentries; i++ {
		binary.Read(r, binary.LittleEndian, &h.Offsets[i])
	}

	anim = Animation{make([]Image, sh.NumFrames), int(sh.FrameDelay)}
	var p [256]color.Color

	if sh.Flags&InternalPalette != 0 {
		p = decodePalette(r)
	}

	if sh.Flags != 0 {
		if sh.Flags&FillBackground != 0 {
			fmt.Print("fillbackground ")
		}
		if sh.Flags&FunctionalColor != 0 {
			fmt.Print("functionalcolor ")
		}
		if sh.Flags&InternalPalette != 0 {
			fmt.Print("internalpalette ")
		}
		if sh.Flags&Junction != 0 {
			fmt.Print("junction ")
		}
	}

	for i := 0; i < int(h.NumFrames); i++ {
		r.Seek(int64(h.Offsets[i]), 0)

		img := *NewImage(image.Rect(0, 0, int(h.Width), int(h.Height)))
		img.FillBackground = sh.Flags&FillBackground != 0
		img.Mix(p)

		var numPix, yIndent, t uint16
		var xPos, yPos int
		var b byte

		binary.Read(r, binary.LittleEndian, &t)

		// frame always starts with 1
		if t == 1 {

			// frame Y indent
			binary.Read(r, binary.LittleEndian, &t)
			yPos += int(t)

			for {
				binary.Read(r, binary.LittleEndian, &numPix)

				if numPix > 0 {
					binary.Read(r, binary.LittleEndian, &t)
					xPos += int(t)

					for j := 0; j < int(numPix); j++ {
						index := yPos*img.Stride + xPos + j
						binary.Read(r, binary.LittleEndian, &img.Pix[index])
						img.Visible[index] = true
					}

					// if the number of pixels is uneven there is a filler byte at the end
					if numPix%2 != 0 {
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
			return anim, errors.New("frame didn't start with 1")
		}

		anim.Frames[i] = img
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

func decodePalette(r io.Reader) Palette {
	ph := paletteHeader{}
	binary.Read(r, binary.LittleEndian, &ph)

	return ConvertPalette(r, int(ph.Index), int(ph.Numcolors))
}

// ConvertPalette converts an 6-bit lbx palette into a color.Palette
func ConvertPalette(r io.Reader, start int, amount int) (p [256]color.Color) {
	pc := paletteColor{}
	for i := 0; i < amount; i++ {
		binary.Read(r, binary.LittleEndian, &pc)

		if pc.A == 1 {
			p[start+i] = color.NRGBA{4 * pc.R, 4 * pc.G, 4 * pc.B, 0}
		} else {
			p[start+i] = color.NRGBA{4 * pc.R, 4 * pc.G, 4 * pc.B, 255}
		}
	}

	return
}
