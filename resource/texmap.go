package resource

import (
	"log"
	"strconv"
	"strings"

	sf "github.com/verokarhu/gorion2/third_party/bitbucket.org/krepa098/gosfml2"
)

type TexMap struct {
	cache map[string]*sf.Texture
	R     *Resource
	Res   sf.Vector2u
}

func (t *TexMap) Preload() {
	for _, v := range t.R.Keys() {
		if strings.Contains(v, ".png") {
			t.Get(v)
		}
	}
}

func (t *TexMap) Get(key string) *sf.Texture {
	if t.cache == nil {
		t.cache = make(map[string]*sf.Texture)
	}

	if v := t.cache[key]; v != nil {
		return v
	}

	if err := t.loadTexture(key); err != nil {
		log.Println(key, err)
		return nil
	}

	return t.Get(key)
}

func (t *TexMap) loadTexture(key string) error {
	texture, err := sf.NewTextureFromMemory(t.R.Get(key), nil)
	if err != nil {
		return err
	}

	t.cache[key] = texture

	return nil
}

func parseAnimationParams(key string) (numframes int, framedelay int) {
	s := strings.TrimSuffix(key, ".png")
	f_index := strings.LastIndex(s, "_f")
	numframes = 1

	if f_index == -1 {
		return
	}

	d_index := strings.LastIndex(s, "_d")

	if d_index == -1 {
		return
	}

	if i, err := strconv.Atoi(s[f_index+2 : d_index]); err == nil {
		if j, err := strconv.Atoi(s[d_index+2:]); err == nil {
			numframes = i
			framedelay = j
		}
	}

	return
}
