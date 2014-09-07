package client

import (
	"log"

	"github.com/verokarhu/gorion2/client/gui"
	res "github.com/verokarhu/gorion2/resource"
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

const (
	None = iota
	Display_Intro
	Display_AnimatedMainmenu
	Display_Mainmenu
	Display_GalaxyMap
)

type state struct {
	display          int
	music            *sf.Music
	Resources        *res.Resource
	rw               *sf.RenderWindow
	tex              res.TexMap
	spr              res.SpriteMap
	buttons          []gui.ButtonMap
	controlsDisabled bool
}

func (s *state) run() {
	s.rw.Clear(sf.ColorBlack())

	switch s.display {
	case Display_Mainmenu:
		s.runMainmenu()
	case Display_GalaxyMap:
		s.runGalaxyMap()
	}

	for _, v := range s.buttons {
		if v.Visible {
			s.rw.Draw(v.Sprite.NextFrame(), sf.DefaultRenderStates())
		}
	}
}

func (s *state) handleMouse(ev sf.Event) {
	if s.controlsDisabled {
		return
	}

	for k, v := range s.buttons {
		if s.buttons[k].MouseOver(sf.MouseGetPosition(s.rw)) {
			switch ev.(type) {
			case sf.EventMouseButtonPressed:
				v.ClickFunc()
			}
		}
	}
}

func (s *state) playMusic(name string, loop bool) {
	if s.music != nil {
		s.music.Stop()
	}

	music, err := sf.NewMusicFromMemory(s.Resources.Get(name))
	if err != nil {
		log.Println(err)
		return
	}

	music.SetLoop(loop)
	music.Play()
	s.music = music
}
