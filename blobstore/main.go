package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	v, err := newVolume("./")
	if err != nil {
		panic(err)
	}
	panic(listen(v))
}

type volume struct {
	path string
}

func newVolume(path string) (*volume, error) {
	if len(path) < 1 {
		return nil, errors.New("path can't be empty")
	}
	if path[len(path)-1:] == "/" {
		path = path[:len(path)-1]
	}
	return &volume{path: path}, nil
}

func (v *volume) object(key string) ([]byte, error) {
	b, err := os.ReadFile(v.path + "/" + key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (v *volume) setObject(key string, b []byte) error {
	dir := filepath.Dir(key)
	if err := os.MkdirAll(v.path+"/"+dir, 0777); err != nil {
		return err
	}
	return os.WriteFile(v.path+"/"+key, b, 0777)
}

func (v *volume) delObject(key string) error {
	err := os.Remove(v.path + "/" + key)
	if err != nil {
		return err
	}
	return nil
}

func listen(v *volume) error {
	r := http.NewServeMux()
	r.HandleFunc("/{hash}", v.handler)
	return http.ListenAndServe(":80", r)
}

func (v *volume) handler(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")

	if r.Method == http.MethodGet {
		if len(hash) < 1 {
			err := errors.New("hash must be defined")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, err := v.object(hash)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(b)
		return
	}

	if r.Method == http.MethodPost {
		if r.Body == nil {
			err := errors.New("missing request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(hash) < 1 {
			err := errors.New("hash must be defined")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("couldn't load request body: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sum, err := encryptSHA256(b)
		if err != nil {
			err = fmt.Errorf("couldn't calculate checksum: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if hash != sum {
			err = fmt.Errorf("checksum mismatch")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = v.setObject(sum, b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Created")
		return
	}

	if r.Method == http.MethodDelete {
		if len(hash) < 1 {
			err := errors.New("hash must be defined")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := v.delObject(hash)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No Content")
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func encryptSHA256(b []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
