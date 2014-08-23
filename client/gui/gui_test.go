package gui

import (
	"testing"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

func Test_CalcPosVec(t *testing.T) {
	p := CalcPosVec(sf.Vector2u{1, 1}, sf.Vector2u{4, 2})

	if p.X != 0.25 {
		t.Error("X was ", p.X)
	}

	if p.Y != 0.5 {
		t.Error("Y was ", p.Y)
	}
}
