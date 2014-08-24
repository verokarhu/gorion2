package image

import (
	"image"
	"image/color"
)

type LbxImage struct {
	Pix            []uint8
	Visible        []bool
	Stride         int
	Rect           image.Rectangle
	Palette        color.Palette
	FillBackground bool
}

var unknown = color.NRGBA{1, 2, 3, 4}
var transparent = color.NRGBA{0, 0, 0, 0}
var black = color.NRGBA{0, 0, 0, 255}

func (i *LbxImage) At(x, y int) color.Color {
	index := i.Rect.Dx()*y + x

	if i.Visible[index] {
		if c := i.Palette[int(i.Pix[index])]; c != unknown {
			return c
		}
	}

	if i.FillBackground {
		return black
	}

	return transparent
}

func (i *LbxImage) Bounds() image.Rectangle { return i.Rect }

func (i *LbxImage) ColorModel() color.Model { return color.NRGBAModel }

func NewLbxImage(r image.Rectangle) *LbxImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, w*h)
	vis := make([]bool, w*h)
	pal := make(color.Palette, 256)

	for k, _ := range pal {
		pal[k] = unknown
	}

	return &LbxImage{pix, vis, w, r, pal, false}
}
