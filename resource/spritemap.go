package resource

import (
	"github.com/verokarhu/gorion2/client/gui"
)

type SpriteMap struct {
	cache map[int]*gui.AnimatedSprite
	Tex   *TexMap
}

func (s *SpriteMap) Flush() {
	s.cache = make(map[int]*gui.AnimatedSprite)
}

func (s *SpriteMap) Get(key int) *gui.AnimatedSprite {
	return s.cache[key]
}

func (s *SpriteMap) Put(key int, texname string) error {
	if s.cache == nil {
		s.Flush()
	}

	nf, fd := parseAnimationParams(texname)
	s.cache[key] = gui.NewAnimatedSprite(nf, fd, s.Tex.Get(texname))

	return nil
}
