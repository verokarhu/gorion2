package resource

import (
	"bytes"
	"image/png"
	"log"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
	"github.com/verokarhu/gorion2/third_party/github.com/nfnt/resize"
)

type Resolution struct {
	X uint
	Y uint
}

type TexMap struct {
	cache map[string]*sf.Texture
	R     *Resource
	Res   Resolution
}

func (t *TexMap) Get(key string) *sf.Texture {
	if t.cache == nil {
		t.cache = make(map[string]*sf.Texture)
	}

	if v := t.cache[key]; v != nil {
		return v
	}

	if err := t.loadTexture(key); err != nil {
		log.Println(err)
		return nil
	}

	return t.Get(key)
}

func (t *TexMap) loadTexture(key string) error {
	img, err := png.Decode(bytes.NewReader(t.R.Get(key)))
	if err != nil {
		return err
	}

	resized := resize.Resize(t.Res.X, t.Res.Y, img, resize.Lanczos3)
	var b bytes.Buffer
	png.Encode(&b, resized)

	texture, err := sf.NewTextureFromMemory(b.Bytes(), nil)
	if err != nil {
		return err
	}

	t.cache[key] = texture

	return nil
}
