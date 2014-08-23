package importer

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/verokarhu/gorion2/lbx"
	"github.com/verokarhu/gorion2/lbx/image"
	res "github.com/verokarhu/gorion2/resource"
)

func Import(dirname string, r *res.Resource) error {
	if err := importVideo(dirname, r); err != nil {
		return err
	}

	if err := importAudio(dirname, r); err != nil {
		return err
	}

	if err := importImages(dirname, r); err != nil {
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

		pals[fmt.Sprintf("%s%d", file.Filename, file.Index)] = image.ConvertPalette(bytes.NewReader(data[file.Index]), 0, 256)
	}

	p := make(color.Palette, 256, 256)

	for i := 0; i < len(p); i++ {
		p[i] = color.NRGBA{150, 0, 0, 128}
	}

	pals["none"] = p

	return
}

func importImages(dirname string, r *res.Resource) error {
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
			data, err := image.Decode(bytes.NewReader(v))

			if err != nil {
				fmt.Println("skipping", file.Filename, k, ":", err)
			}

			for frame, v := range data {
				if file.Palette == "all" {
					for palname, pal := range pals {
						imgcopy := v
						imgcopy.Palette = image.MergePalettes(pal, imgcopy.Palette)

						wg.Add(1)
						go compressPNG(imgcopy, fmt.Sprintf("%s%d-pal-%s.png", file.Filename, k, palname), r, &wg)
					}
					break
				} else {
					imgcopy := v
					imgcopy.Palette = image.MergePalettes(pals[file.Palette], imgcopy.Palette)

					wg.Add(1)
					if len(data) != 1 {
						go compressPNG(imgcopy, fmt.Sprintf("%s%d-%d.png", file.Filename, k, frame+1), r, &wg)
					} else {
						go compressPNG(imgcopy, fmt.Sprintf("%s%d.png", file.Filename, k), r, &wg)
					}
				}
			}
		}
	}

	wg.Wait()
	return nil
}

func compressPNG(img image.LbxImage, key string, r *res.Resource, wg *sync.WaitGroup) {
	var b bytes.Buffer

	if err := png.Encode(&b, &img); err == nil {
		r.Put(key, b.Bytes())
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
