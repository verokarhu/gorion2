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

func (i *LbxImage) At(x int, y int) color.Color {
	index := i.Rect.Dx()*y + x

	if i.Visible[index] {
		return i.Palette[i.Pix[index]]
	}

	return color.NRGBA{0, 0, 0, 0}
}

func (i *LbxImage) Bounds() image.Rectangle { return i.Rect }

func (i *LbxImage) ColorModel() color.Model { return color.NRGBAModel }
