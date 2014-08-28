package gui

import (
	"github.com/verokarhu/gorion2/lbx/dumper/defs"
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

// AnimatedSprite handles animation of sf.Sprites
type AnimatedSprite struct {
	sprites      []*sf.Sprite
	currentframe int
	framedelay   int
}

func NewAnimatedSprite(numframes int, framedelay int, tex *sf.Texture) *AnimatedSprite {
	s := make([]*sf.Sprite, numframes)
	cols := defs.Sheetwidth
	rows := numframes / defs.Sheetwidth
	r := numframes % defs.Sheetwidth

	if numframes < defs.Sheetwidth {
		cols = numframes
	}

	size := tex.GetSize()
	w, h := int(size.X)/cols, int(size.Y)/(rows+1)

	if r == 0 {
		h = int(size.Y) / rows
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			sprite, err := sf.NewSprite(tex)
			if err != nil {
				panic(err)
			}

			sprite.SetTextureRect(sf.IntRect{w * x, h * y, w*x + w, h*y + h})
			s[x+y*defs.Sheetwidth] = sprite
		}
	}

	for x := 0; x < r; x++ {
		sprite, err := sf.NewSprite(tex)
		if err != nil {
			panic(err)
		}

		sprite.SetTextureRect(sf.IntRect{w * x, h * rows, w*x + w, h*rows + h})
		s[x+rows*defs.Sheetwidth] = sprite
	}

	return &AnimatedSprite{s, -1, framedelay}
}

func (s *AnimatedSprite) NextFrame() *sf.Sprite {
	s.currentframe += 1

	if s.currentframe >= len(s.sprites) {
		s.currentframe = 0
	}

	return s.sprites[s.currentframe]
}
