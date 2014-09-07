package gui

import (
	"image"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

type Button struct {
	Rect      image.Rectangle
	Sprite    *AnimatedSprite
	Visible   bool
	ClickFunc func()
}

type Buttons []Button

func scalePosition(pos sf.Vector2i, res sf.Vector2u) (v sf.Vector2f) {
	v.X = float32(pos.X) * float32(res.X) / 640.0
	v.Y = float32(pos.Y) * float32(res.Y) / 480

	return
}

func (b *Button) MouseOver(pos sf.Vector2i) bool {
	b.Visible = false

	if pos.X >= b.Rect.Min.X && pos.X <= b.Rect.Max.X && pos.Y >= b.Rect.Min.Y && pos.Y <= b.Rect.Max.Y {
		b.Visible = true
	}

	return b.Visible
}

func NewButton(xpos, ypos int, res sf.Vector2u, spr *AnimatedSprite, clickevent func()) *Button {
	w, h := spr.GetSize().X, spr.GetSize().Y
	spr.SetPosition(scalePosition(sf.Vector2i{xpos, ypos}, res))

	return &Button{image.Rect(xpos, ypos, xpos+w, ypos+h), spr, false, clickevent}
}
