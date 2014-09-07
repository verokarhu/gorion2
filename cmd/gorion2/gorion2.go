package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/verokarhu/gorion2/client"
	"github.com/verokarhu/gorion2/cmd/gorion2/defs"
	"github.com/verokarhu/gorion2/lbx/dumper"
	res "github.com/verokarhu/gorion2/resource"
)

const resfilename = "moo2_data.blob"
const tempdirectory = "temp_dir_that_gets_deleted"

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var gamedir = flag.String("game", "disc", "path to directory containing moo2 install or the contents of the game disc")
var fullscreen = flag.Bool("fullscreen", false, "fullscreen mode")
var borderless = flag.Bool("borderless", false, "borderless mode")
var skipintro = flag.Bool("skipintro", false, "skip intro")
var width = flag.Uint("width", 640, "window X size")
var height = flag.Uint("height", 480, "window Y size")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Printf(err.Error())
			os.Exit(1)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// todo: load settings from file

	r := res.Resource{}

	if err := r.ReadFile(resfilename); err != nil {
		if !os.IsNotExist(err) {
			log.Println(err.Error())
			os.Exit(1)
		}

		dumplbx()

		if err := r.LoadFiles(client.Resname, tempdirectory); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		if err := os.RemoveAll(tempdirectory); err != nil {
			log.Println(err.Error())
		}

		if err := r.WriteFile(resfilename); err != nil {
			log.Println(err.Error())
		}
	}

	client.Run(client.Params{*width, *height, *fullscreen, *borderless, *skipintro, &r})
}

func dumplbx() {
	if err := os.Mkdir(tempdirectory, os.ModeDir); !os.IsExist(err) && err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, v := range defs.AudioFiles {
		if err := dumper.DumpAudio(*gamedir, tempdirectory, v); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	for _, v := range defs.ImageFiles {
		if err := dumper.DumpImage(*gamedir, tempdirectory, v); err != nil {
			log.Println(err)
		}
	}
}
