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

var transparent = color.NRGBA{0, 0, 0, 0}
var black = color.NRGBA{0, 0, 0, 255}

func (i *LbxImage) At(x, y int) color.Color {
	index := i.Rect.Dx()*y + x

	if i.Visible[index] {
		colno := int(i.Pix[index])
		if colno <= len(i.Palette) {
			if c := i.Palette[colno]; c != nil {
				return c
			}
		}

		return black
	}

	if i.FillBackground {
		return black
	}

	return transparent
}

func (i *LbxImage) Bounds() image.Rectangle { return i.Rect }

func (i *LbxImage) ColorModel() color.Model { return color.NRGBAModel }
