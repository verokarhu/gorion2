package image

import (
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

func Benchmark_Image_At(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for x := 0; x < testimage.Rect.Dx(); x++ {
			for y := 0; y < testimage.Rect.Dy(); y++ {
				testimage.At(x, y)
			}
		}
	}
}
