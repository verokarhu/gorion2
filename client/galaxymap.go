package client

import (
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

const (
	GM_Sprite = iota
	GM_Sprite2
)

const (
	GMButton_Turn = iota
)

func (s *state) galaxymap() {
	s.display = Display_GalaxyMap

	s.playMusic(Resname[3], false)

	s.spr.Put(GM_Sprite, Resname[5]).SetLoop(true)
	s.spr.Put(GM_Sprite2, Resname[1]).SetLoop(true)

	s.clearButtons()
}

func (s *state) runGalaxyMap() {
	s.rw.Draw(s.spr.Get(GM_Sprite).NextFrame(), sf.DefaultRenderStates())
	s.rw.Draw(s.spr.Get(GM_Sprite2).NextFrame(), sf.DefaultRenderStates())
}
