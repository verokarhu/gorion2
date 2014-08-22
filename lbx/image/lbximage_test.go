package image

import (
	"image"
	"image/color"
	"testing"
)

var (
	testimage = LbxImage{
		Pix:     []uint8{0, 11, 11, 12, 0, 0, 0, 0},
		Visible: []bool{1: true, 2: true, 3: true, 7: false},
		Stride:  1,
		Rect:    image.Rect(0, 0, 2, 4),
		Palette: color.Palette{
			11: color.NRGBA{200, 196, 192, 255},
			12: color.NRGBA{2, 1, 2, 0},
		},
	}
)

func TestLbxImage_At(t *testing.T) {
	if c := []color.Color{testimage.At(1, 0), testimage.Palette[11]}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 1), testimage.Palette[12]}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 3), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}
