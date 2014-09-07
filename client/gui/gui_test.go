package gui

import (
	"image"
	"testing"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

func Test_scalePosition(t *testing.T) {
	p := scalePosition(sf.Vector2i{100, 150}, sf.Vector2u{2 * 640, 2 * 480})

	if p.X != 200 {
		t.Error("X was ", p.X)
	}

	if p.Y != 300 {
		t.Error("Y was ", p.Y)
	}
}

func Test_ButtonMap_Contains(t *testing.T) {
	bm := ButtonMap{image.Rect(2, 2, 8, 8), nil, false, nil}

	if bm.MouseOver(sf.Vector2i{1, 3}) {
		t.Error("returned true")
	}

	if bm.MouseOver(sf.Vector2i{9, 3}) {
		t.Error("returned true")
	}

	if bm.MouseOver(sf.Vector2i{3, 1}) {
		t.Error("returned true")
	}

	if bm.MouseOver(sf.Vector2i{1, 3}) {
		t.Error("returned true")
	}

	if bm.Visible {
		t.Error("returned true")
	}

	if !bm.MouseOver(sf.Vector2i{2, 2}) {
		t.Error("returned false")
	}

	if !bm.MouseOver(sf.Vector2i{3, 3}) {
		t.Error("returned false")
	}

	if !bm.MouseOver(sf.Vector2i{8, 8}) {
		t.Error("returned false")
	}

	if !bm.Visible {
		t.Error("returned true")
	}
}
