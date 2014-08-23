package resource

import (
	"encoding/gob"
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

	m := make(map[string][]byte)
	for _, v := range whitelist {
		m[v] = r.Get(v)
	}

	enc := gob.NewEncoder(f)
	if err := enc.Encode(m); err != nil {
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

func (r *Resource) DumpAsFiles(dir string) error {
	if err := os.Mkdir(dir, os.ModeDir); !os.IsExist(err) && err != nil {
		return err
	}

	for k, v := range r.cache {
		if err := ioutil.WriteFile(dir+"/"+k, v, 0644); err != nil {
			return err
		}
	}

	return nil
}
