package resource

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Resource struct {
	cache map[string][]byte
	mutex sync.Mutex
}

type resfile struct {
	Cache map[string][]byte
}

func (r *Resource) ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&r.cache); err != nil {
		return err
	}

	return nil
}

func (r *Resource) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(r.cache); err != nil {
		return err
	}

	return nil
}

func (r *Resource) Get(key string) []byte {
	return r.cache[key]
}

func (r *Resource) Put(key string, value []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cache == nil {
		r.cache = make(map[string][]byte)
	}

	r.cache[key] = value
}

func (r *Resource) Keys() (keys []string) {
	for k, _ := range r.cache {
		keys = append(keys, k)
	}

	return
}

func (r *Resource) LoadDirectory(dir string) error {
	for _, v := range whitelist {
		b, err := loadFile(fmt.Sprintf("%s/%s", dir, v))
		if err != nil {
			return err
		}

		r.Put(v, b)
	}

	return nil
}

func loadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}
