package lbx

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/verokarhu/gorion2/lbx"
	li "github.com/verokarhu/gorion2/lbx/image"
)

const sheetwidth = 10

func Dump(dirname string, dumpdir string) error {
	if err := importVideo(dirname, dumpdir); err != nil {
		return err
	}

	if err := importAudio(dirname, dumpdir); err != nil {
		return err
	}

	if err := importImages(dirname, dumpdir); err != nil {
		return err
	}

	return nil
}

func loadExternalPalettes(dirname string) (pals map[string]li.Palette, err error) {
	pals = make(map[string]li.Palette)

	for _, file := range palette_files {
		data, err := decodeFile(dirname + "/" + file.Filename + ".lbx")
		if err != nil {
			return nil, err
		}

		pals[fmt.Sprintf("%s%d", file.Filename, file.Index)] = li.ConvertPalette(bytes.NewReader(data[file.Index]), 0, 256)
	}

	return
}

func importImages(dirname string, dumpdir string) error {
	pals, err := loadExternalPalettes(dirname)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, file := range image_files {
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

				framescopy := data.Copy()
				filename = fmt.Sprintf("%s-p_black", name)
				framescopy.SetFillBackground(true)
				compress(framescopy, filename, &wg)

			case "none":
				compress(data, name, &wg)

			case "black":
				data.SetFillBackground(true)
				compress(data, name, &wg)

			default:
				data.Mix(pals[file.Palette])
				compress(data, name, &wg)
			}
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

func compressPNG(frames li.Animation, filename string) {
	filename = fmt.Sprintf("%s_f%d%s", filename, len(frames), ".png")
	w, h := frames[0].Rect.Dx(), frames[0].Rect.Dy()

	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	rows := (len(frames) / sheetwidth)
	cols := sheetwidth
	r := len(frames) % sheetwidth

	if len(frames) < sheetwidth {
		cols = len(frames)
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
			f := frames[sheetwidth*y+x]
			draw.Draw(img, target, &f, f.Bounds().Min, draw.Src)
		}
	}

	for x := 0; x < r; x++ {
		target := image.Rect(x*w, rows*h, x*w+w, rows*h+h)
		f := frames[sheetwidth*rows+x]
		draw.Draw(img, target, &f, f.Bounds().Min, draw.Src)
	}

	if err := png.Encode(f, img); err != nil {
		log.Println(err)
	}
}

func importAudio(dirname string, dumpdir string) error {
	for _, filename := range audio_files {
		data, err := decodeFile(dirname + "/" + filename + ".lbx")
		if err != nil {
			return err
		}

		for k, v := range data {
			if err := ioutil.WriteFile(fmt.Sprintf("%s/%s%d.wav", dumpdir, filename, k), v, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func importVideo(dirname string, dumpdir string) error {
	for _, filename := range video_files {
		f, err := os.Open(dirname + "/" + filename + ".lbx")
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.smk", dumpdir, filename), data, 0644); err != nil {
			return err
		}
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
