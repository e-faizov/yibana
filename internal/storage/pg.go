package storage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
)

func NewPgStore(conn string) (interfaces.Store, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return &pgStore{
		db: db,
	}, nil
}

type pgStore struct {
	db *sql.DB
}

func (p *pgStore) Ping() error {
	return p.db.Ping()
}

func (p *pgStore) SetMetric(metric internal.Metric) error {
	return nil
}

func (p *pgStore) GetMetric(metric internal.Metric) (internal.Metric, bool) {
	return internal.Metric{}, false
}

func (p *pgStore) GetAll() []internal.Metric {
	return nil
}
