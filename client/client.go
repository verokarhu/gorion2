package client

import (
	"runtime"
	"time"

	res "github.com/verokarhu/gorion2/resource"
	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

const title = "gorion2 0.0.2"

type Params struct {
	Res_X      uint
	Res_Y      uint
	Fullscreen bool
	Borderless bool
	SkipIntro  bool
	Resources  *res.Resource
}

func Run(p Params) {
	runtime.LockOSThread()
	style := sf.StyleClose

	if p.Fullscreen {
		style = sf.StyleFullscreen
	}

	if p.Borderless {
		style = sf.StyleNone
	}

	texmap := res.TexMap{R: p.Resources, Res: sf.Vector2u{p.Res_X, p.Res_Y}}
	texmap.Preload()

	renderWindow := sf.NewRenderWindow(sf.VideoMode{p.Res_X, p.Res_Y, 32}, title, style, sf.DefaultContextSettings())
	renderWindow.SetFramerateLimit(0)
	renderWindow.SetVSyncEnabled(true)

	ticker := time.NewTicker(time.Second / 10)
	s := state{Resources: p.Resources, rw: renderWindow, tex: texmap, spr: res.SpriteMap{Tex: &texmap}}
	s.mainmenu()

	for renderWindow.IsOpen() {
		select {
		case <-ticker.C:
			for event := renderWindow.PollEvent(); event != nil; event = renderWindow.PollEvent() {
				switch event.(type) {
				case sf.EventClosed:
					renderWindow.Close()
				case sf.EventMouseButtonPressed, sf.EventMouseMoved, sf.EventMouseWheelMoved:
					s.handleMouse(event)
				}
			}

			s.run()
			renderWindow.Display()
		}
	}
}
