package bucket

import (
	"errors"
	"os"
	"path/filepath"
)

type bucket struct {
	path string
}

func New(path string) (*bucket, error) {
	if len(path) < 1 {
		return nil, errors.New("path can't be empty")
	}
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	return &bucket{path: path}, nil
}

func (b *bucket) Get(key string) ([]byte, error) {
	data, err := os.ReadFile(b.path + "/" + key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *bucket) Put(key string, data []byte) error {
	dir := filepath.Dir(key)
	if err := os.MkdirAll(b.path+"/"+dir, 0777); err != nil {
		return err
	}
	return os.WriteFile(b.path+"/"+key, data, 0777)
}
