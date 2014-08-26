package main

import (
	"flag"
	"fmt"
	"github.com/verokarhu/gorion2/lbx/dump"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var dumpdir = flag.String("dir", "dumpdir", "directory where the dumped files go")
var gamedir = flag.String("game", "disc", "path to directory containing moo2 install or the contents of the game disc")
var filename = flag.String("lbx", "", "name of lbx file")
var palette = flag.String("pal", "all", "name of palette to use, list gives the alternatives")
var audio = flag.Bool("a", false, "assume audio content")
var video = flag.Bool("v", false, "assume video content")
var image = flag.Bool("i", false, "assume image content")

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

	if *palette == "list" {
		fmt.Println("available palettes:")
		for _, v := range dumper.PaletteFiles {
			fmt.Println(fmt.Sprintf("%s%d", v.Filename, v.Index))
		}
		os.Exit(0)
	}

	if *filename != "" {
		if err := os.Mkdir(*dumpdir, os.ModeDir); !os.IsExist(err) && err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if *video || (!*audio && !*video && !*image) {
			if err := dumper.DumpVideo(*gamedir, *dumpdir, *filename); err != nil {
				log.Println(err)
			}
		}

		if *audio || (!*audio && !*video && !*image) {
			if err := dumper.DumpAudio(*gamedir, *dumpdir, *filename); err != nil {
				log.Println(err)
			}
		}

		if *image || (!*audio && !*video && !*image) {
			i := dumper.ImagePair{*filename, *palette}
			if *palette == "" {
				i.Palette = "all"
			}

			if err := dumper.DumpImage(*gamedir, *dumpdir, i); err != nil {
				log.Println(err)
			}
		}
	} else {
		flag.PrintDefaults()
	}
}
