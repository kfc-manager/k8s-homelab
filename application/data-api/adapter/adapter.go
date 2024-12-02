package adapter

import "github.com/kfc-manager/k8s-homelab/application/data-api/domain"

type Database interface {
	Close() error
	InsertImage(img *domain.Image) error
}

type Bucket interface {
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
}

type Queue interface {
	Close() error
	Send(msg string) error
}

type Server interface {
	Listen() error
}
