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
	buttons          gui.Buttons
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

	s.buttons.ForEach(func(b *gui.Button) {
		if b.Visible {
			s.rw.Draw(b.Sprite.NextFrame(), sf.DefaultRenderStates())
		}
	})
}

func (s *state) handleMouse(ev sf.Event) {
	if s.controlsDisabled {
		return
	}

	mousepos := sf.MouseGetPosition(s.rw)
	s.buttons.ForEach(func(b *gui.Button) {
		b.Update(mousepos)

		if b.MouseOver(mousepos) {
			switch ev.(type) {
			case sf.EventMouseButtonPressed:
				b.ClickFunc()
			}
		}
	})
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

func (s *state) clearButtons() {
	s.buttons = gui.Buttons{}
}
