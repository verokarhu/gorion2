package image

import (
	"image"
	"image/color"
)

type Animation struct {
	Frames     []Image
	FrameDelay int
}

type Palette [256]color.Color

type Image struct {
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

func (i *Image) At(x, y int) color.Color {
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

func (i *Image) Bounds() image.Rectangle { return i.Rect }

func (i *Image) ColorModel() color.Model { return color.NRGBAModel }

func NewImage(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, w*h)
	vis := make([]bool, w*h)
	var pal Palette

	for i := 0; i < 256; i++ {
		pal[i] = color_missing
	}

	return &Image{pix, vis, w, r, pal, false}
}

// Mix grabs colors missing from Image's palette from override
func (i *Image) Mix(override Palette) {
	for k, v := range override {
		if i.Palette[k] == color_missing && v != nil {
			i.Palette[k] = v
		}
	}
}

func (i *Image) Copy() (cop Image) {
	cop = Image{make([]uint8, len(i.Pix)), make([]bool, len(i.Visible)), i.Stride, i.Rect, i.Palette, i.FillBackground}
	for k, v := range i.Pix {
		cop.Pix[k] = v
	}

	for k, v := range i.Visible {
		cop.Visible[k] = v
	}

	return
}

func Blend(base Image, override Image) (result Image) {
	result = base.Copy()
	for k, v := range override.Pix {
		if k < len(override.Visible) && override.Visible[k] {
			result.Pix[k] = v
		}
	}

	return
}

func (anim Animation) Mix(override Palette) {
	for k, _ := range anim.Frames {
		anim.Frames[k].Mix(override)
	}
}

func (anim Animation) SetFillBackground(flag bool) {
	for k, _ := range anim.Frames {
		anim.Frames[k].FillBackground = flag
	}
}

func (anim Animation) Copy() (cop Animation) {
	cop = Animation{make([]Image, len(anim.Frames)), anim.FrameDelay}
	for k, v := range anim.Frames {
		cop.Frames[k] = v
	}

	return
}

func (anim *Animation) BlendFrames() {
	for i := 1; i < len(anim.Frames); i++ {
		anim.Frames[i] = Blend(anim.Frames[i-1], anim.Frames[i])
	}
}
