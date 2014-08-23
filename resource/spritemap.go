package resource

import (
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

type SpriteMap struct {
	cache map[int]*sf.Sprite
}

func (s *SpriteMap) Flush() {
	s.cache = make(map[int]*sf.Sprite)
}

func (s *SpriteMap) Get(key int) *sf.Sprite {
	return s.cache[key]
}

func (s *SpriteMap) Put(key int, tex *sf.Texture) error {
	spr, err := sf.NewSprite(tex)
	if err != nil {
		return err
	}

	if s.cache == nil {
		s.Flush()
	}

	s.cache[key] = spr

	return nil
}
