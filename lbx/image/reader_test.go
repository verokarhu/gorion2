package image

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

var (
	lbximg = []byte{
		4, 0, 3, 0, 0, 0, 2, 0,
		10, 0, 0, 0x14, 36, 0, 0, 0,
		68, 0, 0, 0, 88, 0, 0, 0,
		10, 0, 2, 0, 0, 10, 11, 12,
		1, 50, 49, 48, 1, 0, 0, 0,
		3, 0, 1, 0, 11, 11, 10, 0,
		0, 0, 2, 0, 2, 0, 0, 0,
		200, 200, 2, 0, 0, 0, 200, 200,
		0, 0, 0xe8, 0x03, 1, 0, 0, 0,
		8, 0, 0, 0, 200, 200, 200, 200,
		10, 10, 11, 11, 0, 0, 0xe8, 0x03,
	}
	expectedpalette = color.Palette{
		10: color.NRGBA{40, 44, 48, 255},
		11: color.NRGBA{200, 196, 192, 0},
		12: color.NRGBA{2, 1, 2, 255},
	}
	basepalette = color.Palette{
		10:  color.NRGBA{40, 44, 48, 255},
		11:  color.NRGBA{2, 1, 1, 0},
		255: nil,
	}
	overridepalette = color.Palette{
		11: color.NRGBA{200, 196, 192, 0},
		12: color.NRGBA{2, 1, 2, 255},
	}
)

func Test_Decode(t *testing.T) {
	f := bytes.NewReader(lbximg)

	expected := [2]Image{
		Image{
			Pix:    []uint8{0, 11, 11, 10, 0, 0, 0, 0, 200, 200, 200, 200},
			Stride: 4,
			Rect:   image.Rect(0, 0, 4, 3),
		},
		Image{
			Pix:    []uint8{200, 200, 200, 200, 10, 10, 11, 11, 200, 200, 200, 200},
			Stride: 4,
			Rect:   image.Rect(0, 0, 4, 3),
		},
	}

	decoded, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	if c := []int{len(decoded.Frames), 2}; c[0] != c[1] {
		t.Fatal("excepted ", c[1], ", returned ", c[0])
	}

	if c := []int{len(decoded.Frames[0].Palette), 256}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []int{len(decoded.Frames[1].Palette), 256}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []int{decoded.Frames[0].Stride, expected[0].Stride}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []int{decoded.Frames[1].Stride, expected[1].Stride}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := [][]uint8{decoded.Frames[0].Pix, expected[0].Pix}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := [][]uint8{decoded.Frames[1].Pix, expected[1].Pix}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []image.Rectangle{decoded.Frames[0].Rect, expected[0].Rect}; !c[0].Size().Eq(c[1].Size()) {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []image.Rectangle{decoded.Frames[1].Rect, expected[1].Rect}; !c[0].Size().Eq(c[1].Size()) {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}

func Test_ConvertPalette(t *testing.T) {
	f := bytes.NewReader(lbximg[28:36])

	p := ConvertPalette(f, 10, 2)

	if c := []int{len(p), 256}; c[0] != c[1] {
		t.Fatal("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{p[10], expectedpalette[10]}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{p[11], expectedpalette[11]}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}
