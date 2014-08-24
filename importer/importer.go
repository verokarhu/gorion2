package importer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/verokarhu/gorion2/lbx"
	li "github.com/verokarhu/gorion2/lbx/image"
	res "github.com/verokarhu/gorion2/resource"
)

func Import(dirname string, r *res.Resource, dumpgif bool) error {
	if err := importVideo(dirname, r); err != nil {
		return err
	}

	if err := importAudio(dirname, r); err != nil {
		return err
	}

	if err := importImages(dirname, r, dumpgif); err != nil {
		return err
	}

	return nil
}

func loadExternalPalettes(dirname string) (pals map[string]color.Palette, err error) {
	pals = make(map[string]color.Palette)

	for _, file := range palette_files {
		data, err := decodeLBX(dirname + "/" + file.Filename + ".lbx")
		if err != nil {
			return nil, err
		}

		pals[fmt.Sprintf("%s%d", file.Filename, file.Index)] = li.ConvertPalette(bytes.NewReader(data[file.Index]), 0, 256)
	}

	p := make(color.Palette, 256)
	for i := 0; i < len(p); i++ {
		p[i] = color.NRGBA{0, 0, 0, 255}
	}

	pals["black"] = p

	p = make(color.Palette, 256)
	for i := 0; i < len(p); i++ {
		p[i] = color.NRGBA{0, 0, 0, 0}
	}

	pals["transparent"] = p

	return
}

func importImages(dirname string, r *res.Resource, dumpgif bool) error {
	pals, err := loadExternalPalettes(dirname)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, file := range image_files {
		data, err := decodeLBX(dirname + "/" + file.Filename + ".lbx")
		if err != nil {
			return err
		}

		for k, v := range data {
			data, err := li.Decode(bytes.NewReader(v))

			if err != nil {
				fmt.Println("skipping", file.Filename, k, ":", err)
			}

			var frames []*li.LbxImage

			if file.Palette != "all" {
				frames = make([]*li.LbxImage, len(data))
			} else {
				frames = make([]*li.LbxImage, len(data)*len(pals))
			}

			name := fmt.Sprintf("%s%d", file.Filename, k)
			for frame, v := range data {
				if file.Palette == "all" {
					i := 0
					for _, pal := range pals {
						imgcopy := v
						imgcopy.Palette = li.MergePalettes(pal, imgcopy.Palette)

						if imgcopy.FillBackground {
							imgcopy.Palette = li.MergePalettes(pals["black"], imgcopy.Palette)
						} else {
							imgcopy.Palette = li.MergePalettes(pals["transparent"], imgcopy.Palette)
						}

						frames[i*len(data)+frame] = imgcopy
						i++
					}
				} else {
					imgcopy := v
					if file.Palette != "none" {
						imgcopy.Palette = li.MergePalettes(pals[file.Palette], imgcopy.Palette)
					}

					if imgcopy.FillBackground {
						imgcopy.Palette = li.MergePalettes(pals["black"], imgcopy.Palette)
					} else {
						imgcopy.Palette = li.MergePalettes(pals["transparent"], imgcopy.Palette)
					}

					frames[frame] = imgcopy
				}
			}

			wg.Add(1)
			if dumpgif {
				go compressGIF(frames, name, r, &wg)
			} else {
				go compressPNG(frames, name, r, &wg)
			}
		}
	}

	wg.Wait()
	return nil
}

func compressPNG(frames []*li.LbxImage, key string, r *res.Resource, wg *sync.WaitGroup) {
	w, h := frames[0].Rect.Dx(), frames[0].Rect.Dy()
	rect := image.Rect(0, 0, w*len(frames), h)
	img := image.NewNRGBA(rect)

	for k, v := range frames {
		target := image.Rect(k*w, 0, k*w+w, h)
		draw.Draw(img, target, v, v.Bounds().Min, draw.Src)
	}

	var b bytes.Buffer

	if err := png.Encode(&b, img); err == nil {
		r.Put(fmt.Sprintf("%s_frames%d.png", key, len(frames)), b.Bytes())
	} else {
		log.Println(err)
	}

	wg.Done()
}

func compressGIF(frames []*li.LbxImage, key string, r *res.Resource, wg *sync.WaitGroup) {
	g := gif.GIF{make([]*image.Paletted, len(frames)), make([]int, len(frames)), 0}

	for k, v := range frames {
		img := image.NewPaletted(v.Rect, v.Palette)
		draw.Draw(img, img.Bounds(), v, v.Bounds().Min, draw.Src)
		g.Image[k] = img
		g.Delay[k] = 10
	}

	var b bytes.Buffer

	if err := gif.EncodeAll(&b, &g); err == nil {
		r.Put(key+".gif", b.Bytes())
	} else {
		log.Println(err)
	}

	wg.Done()
}

func importAudio(dirname string, r *res.Resource) error {
	for _, filename := range audio_files {
		data, err := decodeLBX(dirname + "/" + filename + ".lbx")
		if err != nil {
			return err
		}

		for k, v := range data {
			r.Put(fmt.Sprintf("%s%d.wav", filename, k), v)
		}
	}

	return nil
}

func importVideo(dirname string, r *res.Resource) error {
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

		r.Put(filename+".smk", data)
	}

	return nil
}

func decodeLBX(filename string) ([][]byte, error) {
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
