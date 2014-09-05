package image

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

var testimage = NewImage(image.Rect(0, 0, 2, 4))

func Test_Image_At(t *testing.T) {
	testimage.Pix = []uint8{210, 11, 11, 12, 150, 200, 0, 250}
	testimage.Visible = []bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 7: false}
	testimage.Palette[11] = color.NRGBA{200, 196, 192, 255}
	testimage.Palette[12] = color.NRGBA{2, 1, 2, 0}

	if c := []color.Color{testimage.At(0, 0), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 0), testimage.Palette[11]}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 1), color.NRGBA{2, 1, 2, 255}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(0, 2), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 2), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(0, 3), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.At(1, 3), color.NRGBA{0, 0, 0, 0}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	testimage.FillBackground = true
	if c := []color.Color{testimage.At(0, 0), color.NRGBA{0, 0, 0, 255}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

}

func Test_Image_Mix(t *testing.T) {
	pal := Palette{50: color.NRGBA{200, 200, 100, 100}, 150: color.NRGBA{10, 10, 20, 20}}
	testimage.Mix(pal)

	if c := []color.Color{testimage.Palette[50], color.NRGBA{200, 200, 100, 100}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := []color.Color{testimage.Palette[150], color.NRGBA{10, 10, 20, 20}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	pal = Palette{50: color.NRGBA{10, 10, 20, 20}}
	testimage.Mix(pal)

	if c := []color.Color{testimage.Palette[50], color.NRGBA{200, 200, 100, 100}}; c[0] != c[1] {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}

func Test_Image_Blend(t *testing.T) {
	base, override := NewImage(image.Rect(0, 0, 5, 1)), NewImage(image.Rect(0, 0, 5, 1))
	base.Pix = []uint8{1, 1, 1, 1, 1}
	override.Pix = []uint8{5, 5, 5, 5, 5}
	override.Visible = []bool{2: true}
	expected := []uint8{1, 1, 5, 1, 1}

	result := Blend(*base, *override)

	if c := [][]uint8{result.Pix, expected}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}

func Test_Animation_BlendFrames(t *testing.T) {
	f0, f1, f2 := NewImage(image.Rect(0, 0, 5, 1)), NewImage(image.Rect(0, 0, 5, 1)), NewImage(image.Rect(0, 0, 5, 1))
	f0.Pix = []uint8{1, 1, 1, 1, 1}
	f1.Pix = []uint8{5, 5, 5, 5, 5}
	f2.Pix = f1.Pix
	f1.Visible = []bool{0: true}
	f2.Visible = []bool{4: true}
	anim := Animation{Frames: []Image{*f0, *f1, *f2}}
	expected := [][]uint8{f0.Pix, []uint8{5, 1, 1, 1, 1}, []uint8{5, 1, 1, 1, 5}}

	anim.BlendFrames()

	if c := [][]uint8{anim.Frames[0].Pix, expected[0]}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := [][]uint8{anim.Frames[1].Pix, expected[1]}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}

	if c := [][]uint8{anim.Frames[2].Pix, expected[2]}; bytes.Compare(c[0], c[1]) != 0 {
		t.Error("excepted ", c[1], ", returned ", c[0])
	}
}

func Benchmark_Image_At(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < testimage.Rect.Dx(); x++ {
			for y := 0; y < testimage.Rect.Dy(); y++ {
				testimage.At(x, y)
			}
		}
	}
}
