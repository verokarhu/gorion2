package dumper

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/verokarhu/gorion2/lbx"
	"github.com/verokarhu/gorion2/lbx/dumper/defs"
	li "github.com/verokarhu/gorion2/lbx/image"
)

type ImagePair struct {
	Filename string
	Palette  string
}

func DumpImage(dirname string, targetdir string, files ...ImagePair) error {
	for _, file := range files {
		if err := importImage(dirname, targetdir, file); err != nil {
			return err
		}
	}

	return nil
}

func DumpAudio(dirname string, targetdir string, filenames ...string) error {
	for _, filename := range filenames {
		if err := importAudio(dirname, targetdir, filename); err != nil {
			return err
		}
	}

	return nil
}

func DumpVideo(dirname string, targetdir string, filenames ...string) error {
	for _, filename := range filenames {
		if err := importVideo(dirname, targetdir, filename); err != nil {
			return err
		}
	}

	return nil
}

func importImage(dirname string, dumpdir string, file ImagePair) error {
	pals, err := loadExternalPalettes(dirname)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	data, err := decodeFile(dirname + "/" + file.Filename + ".lbx")
	if err != nil {
		return err
	}

	for k, v := range data {
		fmt.Printf("decoding %s-%d:  ", file.Filename, k)

		data, err := li.Decode(bytes.NewReader(v))
		if err != nil {
			fmt.Println("skipping", file.Filename, k, ":", err)
		}

		fmt.Println()

		name := fmt.Sprintf("%s/%s_%d", dumpdir, file.Filename, k)
		switch file.Palette {
		case "all":
			for palname, pal := range pals {
				filename := fmt.Sprintf("%s-p_%s", name, palname)
				framescopy := data.Copy()
				framescopy.Mix(pal)
				compress(framescopy, filename, &wg)
			}

			filename := fmt.Sprintf("%s-p_none", name)
			compress(data.Copy(), filename, &wg)

		case "none":
			compress(data, name, &wg)

		default:
			data.Mix(pals[file.Palette])
			compress(data, name, &wg)
		}
	}

	wg.Wait()
	return nil
}

func compress(frames li.Animation, filename string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		compressPNG(frames, filename)
		wg.Done()
	}()
}

func compressPNG(anim li.Animation, filename string) {
	filename = fmt.Sprintf("%s_f%d_d%d%s", filename, len(anim.Frames), anim.FrameDelay, ".png")
	w, h := anim.Frames[0].Rect.Dx(), anim.Frames[0].Rect.Dy()

	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	rows := (len(anim.Frames) / defs.Sheetwidth)
	cols := defs.Sheetwidth
	r := len(anim.Frames) % defs.Sheetwidth

	if len(anim.Frames) < defs.Sheetwidth {
		cols = len(anim.Frames)
	}

	var img *image.NRGBA
	if r > 0 {
		img = image.NewNRGBA(image.Rect(0, 0, cols*w, rows*h+h))
	} else {
		img = image.NewNRGBA(image.Rect(0, 0, cols*w, rows*h))
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			target := image.Rect(x*w, y*h, x*w+w, y*h+h)
			f := anim.Frames[defs.Sheetwidth*y+x]
			draw.Draw(img, target, &f, f.Bounds().Min, draw.Src)
		}
	}

	for x := 0; x < r; x++ {
		target := image.Rect(x*w, rows*h, x*w+w, rows*h+h)
		f := anim.Frames[defs.Sheetwidth*rows+x]
		draw.Draw(img, target, &f, f.Bounds().Min, draw.Src)
	}

	if err := png.Encode(f, img); err != nil {
		log.Println(err)
	}
}

func importAudio(dirname string, targetdir string, filename string) error {
	data, err := decodeFile(dirname + "/" + filename + ".lbx")
	if err != nil {
		return err
	}

	for k, v := range data {
		if string(v[:4]) == "RIFF" {
			if err := ioutil.WriteFile(fmt.Sprintf("%s/%s%d.wav", targetdir, filename, k), v, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func importVideo(dirname string, targetdir string, filename string) error {
	f, err := os.Open(dirname + "/" + filename + ".lbx")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if string(data[:3]) != "SMK" {
		return errors.New("not a smacker video file")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.smk", targetdir, filename), data, 0644); err != nil {
		return err
	}

	return nil
}

func decodeFile(filename string) ([][]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := lbx.Decode(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
