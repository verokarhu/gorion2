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

func parseNumframes(key string) int {
	s := strings.TrimSuffix(key, ".png")
	index := strings.LastIndex(s, "_f")

	if index != -1 {
		if i, err := strconv.Atoi(s[index+2:]); err == nil {
			return i
		}
	}

	return 1
}
