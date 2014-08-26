package main

import (
	"flag"
	"github.com/verokarhu/gorion2/lbx/dump"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var dumpdir = flag.String("dir", "dumpdir", "directory where the dumped files go")
var gamedisc = flag.String("disc", "disc", "path to game disc")
var filename = flag.String("lbx", "", "name of lbx file")
var palette = flag.String("pal", "list", "name of palette to use, list lists the alternatives")
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

	if *filename != "" {
		if err := os.Mkdir(*dumpdir, os.ModeDir); !os.IsExist(err) && err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if *video || (!*audio && !*video && !*image) {
			if err := dumper.DumpVideo(*gamedisc, *dumpdir, *filename); err != nil {
				log.Println(err)
			}
		}

		if *audio || (!*audio && !*video && !*image) {
			if err := dumper.DumpAudio(*gamedisc, *dumpdir, *filename); err != nil {
				log.Println(err)
			}
		}

		if *image || (!*audio && !*video && !*image) {
			i := dumper.ImagePair{*filename, *palette}
			if *palette == "" {
				i.Palette = "all"
			}

			if err := dumper.DumpImage(*gamedisc, *dumpdir, i); err != nil {
				log.Println(err)
			}
		}
	} else {
		flag.PrintDefaults()
	}
}