package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kfc-manager/k8s-homelab/application/data-api/domain"
)

type database struct {
	conn *pgxpool.Pool
}

func New(host, port, name, user, pass string) (*database, error) {
	conn, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name),
	)
	if err != nil {
		return nil, err
	}
	return &database{conn: conn}, nil
}

func (db *database) Close() error {
	db.conn.Close()
	return nil
}

func (db *database) InsertImage(img *domain.Image) error {
	_, err := db.conn.Exec(
		context.Background(),
		`INSERT INTO image (hash, size, format, created_at, source, width, height)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		img.Hash,
		img.Size,
		img.Format,
		img.CreatedAt,
		img.Source,
		img.Width,
		img.Height,
	)
	if err != nil {
		return err
	}

	return nil
}
