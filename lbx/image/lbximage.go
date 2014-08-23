package image

import (
	"image"
	"image/color"
)

type LbxImage struct {
	Pix     []uint8
	Visible []bool
	Stride  int
	Rect    image.Rectangle
	Palette color.Palette
}

var transparent = color.NRGBA{0, 0, 0, 0}

func (i *LbxImage) At(x int, y int) color.Color {
	index := i.Rect.Dx()*y + x

	if i.Visible[index] {
		if c := i.Palette[i.Pix[index]]; c != nil {
			return c
		}
	}

	return transparent
}

func (i *LbxImage) Bounds() image.Rectangle { return i.Rect }

func (i *LbxImage) ColorModel() color.Model { return color.NRGBAModel }
