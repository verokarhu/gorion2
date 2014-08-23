package gui

import (
	"testing"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

func Test_CalcPosVec(t *testing.T) {
	p := CalcPosVec(sf.Vector2i{1, 1}, sf.Vector2u{4, 2})

	if p.X != 0.25 {
		t.Error("X was ", p.X)
	}

	if p.Y != 0.5 {
		t.Error("Y was ", p.Y)
	}
}

func Test_ButtonMap_Contains(t *testing.T) {
	bm := ButtonMap{sf.Vector2f{0.2, 0.2}, sf.Vector2f{0.3, 0.3}}

	if bm.Contains(sf.Vector2f{0.1, 0.25}) {
		t.Error("returned true")
	}

	if bm.Contains(sf.Vector2f{0.4, 0.25}) {
		t.Error("returned true")
	}

	if bm.Contains(sf.Vector2f{0.25, 0.1}) {
		t.Error("returned true")
	}

	if bm.Contains(sf.Vector2f{0.1, 0.4}) {
		t.Error("returned true")
	}

	if !bm.Contains(sf.Vector2f{0.2, 0.2}) {
		t.Error("returned false")
	}

	if !bm.Contains(sf.Vector2f{0.25, 0.25}) {
		t.Error("returned false")
	}

	if !bm.Contains(sf.Vector2f{0.3, 0.3}) {
		t.Error("returned false")
	}
}
