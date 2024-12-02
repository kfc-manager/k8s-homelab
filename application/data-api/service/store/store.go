package store

import (
	"errors"

	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter"
	"github.com/kfc-manager/k8s-homelab/application/data-api/domain"
)

type Service interface {
	Get(hash string) ([]byte, error)
	Set(hash string, b []byte) error
}

type service struct {
	db     adapter.Database
	bucket adapter.Bucket
	queue  adapter.Queue
}

func New(db adapter.Database, b adapter.Bucket, q adapter.Queue) *service {
	return &service{db: db, bucket: b, queue: q}
}

func (s *service) Get(hash string) ([]byte, error) {
	return s.bucket.Get(hash)
}

func (s *service) Set(hash string, b []byte) error {
	img, err := domain.LoadImage(b)
	if err != nil {
		return err
	}
	if img.Hash != hash {
		return errors.New("hash mismatch")
	}

	err = s.bucket.Put(img.Hash, b)
	if err != nil {
		return err
	}

	err = s.db.InsertImage(img)
	if err != nil {
		return err
	}

	err = s.queue.Send(img.Hash)
	if err != nil {
		return err
	}

	return nil
}
