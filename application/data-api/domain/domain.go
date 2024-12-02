package domain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"image"
	_ "image/jpeg"
	_ "image/png"
)

type Image struct {
	Hash      string
	Size      int
	Format    string
	CreatedAt time.Time
	Source    string
	Width     int
	Height    int
}

func LoadImage(b []byte) (*Image, error) {
	img, form, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	h, err := encryptSHA256(b)
	if err != nil {
		return nil, err
	}

	result := &Image{
		Hash:      h,
		Format:    form,
		CreatedAt: time.Now(),
		Size:      len(b),
	}
	result.Width = img.Bounds().Dx()
	result.Height = img.Bounds().Dy()

	return result, nil
}

func encryptSHA256(b []byte) (string, error) {
	hash := sha256.New()
	_, err := hash.Write(b)
	if err != nil {
		return "", err
	}
	s := hash.Sum(nil)
	return hex.EncodeToString(s), nil
}
