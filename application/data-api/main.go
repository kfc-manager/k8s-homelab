package main

import (
	"fmt"
	"os"

	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter"
	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter/bucket"
	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter/database"
	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter/queue"
	"github.com/kfc-manager/k8s-homelab/application/data-api/adapter/server"
	"github.com/kfc-manager/k8s-homelab/application/data-api/service/store"
)

func main() {
	var db adapter.Database
	var b adapter.Bucket
	var q adapter.Queue
	var err error

	db, err = database.New(
		envOrPanic("DB_HOST"),
		envOrPanic("DB_PORT"),
		envOrPanic("DB_NAME"),
		envOrPanic("DB_USER"),
		envOrPanic("DB_PASS"),
	)
	if err != nil {
		panic(err)
	}

	b, err = bucket.New(
		envOrPanic("BUCKET_PATH"),
	)
	if err != nil {
		panic(err)
	}

	q, err = queue.New(
		envOrPanic("QUEUE_HOST"),
		envOrPanic("QUEUE_PORT"),
		envOrPanic("QUEUE_NAME"),
	)
	if err != nil {
		panic(err)
	}

	s := server.New("80", store.New(db, b, q))
	if err := s.Listen(); err != nil {
		panic(err)
	}
}

func envOrPanic(key string) string {
	val := os.Getenv(key)
	if len(val) <= 0 {
		panic(fmt.Errorf("missing environment variable: %s", key))
	}
	return val
}
