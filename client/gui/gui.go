package gui

import (
	"image"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

type ButtonMap struct {
	Rect         image.Rectangle
	SpriteNumber int
	Visible      bool
	ClickFunc    func()
}

// CalcPosVec calculates a resolution independent vector of the x,y position
func CalcPosVec(pos sf.Vector2i, res sf.Vector2u) (v sf.Vector2f) {
	v.X = float32(pos.X) / float32(res.X)
	v.Y = float32(pos.Y) / float32(res.Y)

	return
}

func (b *ButtonMap) MouseOver(pos sf.Vector2i) bool {
	b.Visible = false

	if pos.X >= b.Rect.Min.X && pos.X <= b.Rect.Max.X && pos.Y >= b.Rect.Min.Y && pos.Y <= b.Rect.Max.Y {
		b.Visible = true
	}

	return b.Visible
}

func NewButtonMap(x, y, width, height, spriteno int, clickevent func()) *ButtonMap {
	return &ButtonMap{image.Rect(x, y, x+width, y+height), spriteno, false, clickevent}
}
