package image

import (
	"image"
	"image/color"
)

type LbxAnimation []LbxImage
type Palette [256]color.Color

type LbxImage struct {
	Pix            []uint8
	Visible        []bool
	Stride         int
	Rect           image.Rectangle
	Palette        Palette
	FillBackground bool
}

var color_missing = color.NRGBA{1, 2, 3, 4}
var color_transparent = color.NRGBA{0, 0, 0, 0}
var color_black = color.NRGBA{0, 0, 0, 255}

func (i *LbxImage) At(x, y int) color.Color {
	index := i.Rect.Dx()*y + x

	if i.Visible[index] {
		if c := i.Palette[int(i.Pix[index])].(color.NRGBA); c != color_missing {
			if c.A == 0 {
				// transparent pixels in the palette aren't actually transparent
				c.A = 255
			}

			return c
		}
	}

	if i.FillBackground {
		return color_black
	}

	return color_transparent
}

func (i *LbxImage) Bounds() image.Rectangle { return i.Rect }

func (i *LbxImage) ColorModel() color.Model { return color.NRGBAModel }

func NewLbxImage(r image.Rectangle) *LbxImage {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, w*h)
	vis := make([]bool, w*h)
	var pal Palette

	for i := 0; i < 256; i++ {
		pal[i] = color_missing
	}

	return &LbxImage{pix, vis, w, r, pal, false}
}

// Mix grabs colors missing from LbxImage's palette from override
func (i *LbxImage) Mix(override Palette) {
	for k, v := range override {
		if i.Palette[k] == color_missing && v != nil {
			i.Palette[k] = v
		}
	}
}

func (anim LbxAnimation) Mix(override Palette) {
	for k, _ := range anim {
		anim[k].Mix(override)
	}
}

func (anim LbxAnimation) SetFillBackground(flag bool) {
	p := anim
	for k, _ := range anim {
		p[k].FillBackground = flag
	}
}

func (anim LbxAnimation) Copy() (cop LbxAnimation) {
	cop = make(LbxAnimation, len(anim))
	for k, v := range anim {
		cop[k] = v
	}

	return
}
