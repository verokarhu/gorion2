package client

import (
	"github.com/verokarhu/gorion2/client/gui"
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

const (
	MM_Background = iota
	MMButton_Continue
	MMButton_Load
	MMButton_NewGame
	MMButton_Multiplayer
	MMButton_HallOfFame
	MMButton_QuitGame
)

func (s *state) mainmenu() {
	s.display = Display_Mainmenu

	s.playMusic(Resname[2], true)

	s.spr.Put(MM_Background, Resname[4])

	s.buttons = gui.Buttons{
		*gui.NewButton(415, 172, s.rw.GetSize(), s.spr.Put(MMButton_Continue, Resname[6]), s.galaxymap),
		*gui.NewButton(415, 285, s.rw.GetSize(), s.spr.Put(MMButton_QuitGame, Resname[7]), s.rw.Close),
	}

	s.controlsDisabled = true
}

func (s *state) runMainmenu() {
	spr := s.spr.Get(MM_Background)
	s.rw.Draw(spr.NextFrame(), sf.DefaultRenderStates())

	if spr.Stopped() {
		s.controlsDisabled = false
	}
}
