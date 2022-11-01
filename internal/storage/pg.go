package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
)

func createMetricsTable(ctx context.Context, db *sql.DB) error {
	sql := `create table metrics
(
	mid serial
		constraint metrics_pk
			primary key,
	id text not null,
	mtid int not null,
	value double precision,
	delta bigint,
	hash text
)`
	_, err := db.ExecContext(ctx, sql)
	return err
}

func createMTypesTable(ctx context.Context, db *sql.DB) error {
	sqlCreate := `create table metric_types
	(
		mtid serial,
		mtname text
	)`
	sqlIndex := `create unique index metric_types_mtid_uindex
	on metric_types (mtid)`
	sqlAlter := `alter table metric_types
	add constraint metric_types_pk
		primary key (mtid)`

	_, err := db.ExecContext(ctx, sqlCreate)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, sqlIndex)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, sqlAlter)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("insert into metric_types (mtname) values ('%s')", internal.GaugeType))
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, fmt.Sprintf("insert into metric_types (mtname) values ('%s')", internal.CounterType))
	if err != nil {
		return err
	}
	return nil
}

func clearTable(db *sql.DB) {
	db.Exec("drop table metric_types")
	db.Exec("drop table metrics")
}

func tableExist(ctx context.Context, db *sql.DB, tb string) bool {
	sql := `SELECT EXISTS (
	   SELECT FROM information_schema.tables
	   WHERE  table_schema = 'public'
	   AND    table_name   = $1
	   )
`

	var exists bool
	row := db.QueryRowContext(ctx, sql, tb)
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func initTables(ctx context.Context, db *sql.DB) error {
	//clearTable(db)
	var err error
	exist := tableExist(ctx, db, "metric_types")
	if !exist {
		err = createMTypesTable(ctx, db)
		if err != nil {
			return fmt.Errorf("error create metric_types: %v", err)
		}
	}
	exist = tableExist(ctx, db, "metrics")
	if !exist {
		err = createMetricsTable(ctx, db)
		if err != nil {
			return fmt.Errorf("error create metrics: %v", err)
		}
	}
	return nil
}

func NewPgStore(conn string) (interfaces.Store, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = initTables(ctx, db)
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

func (p *pgStore) SetMetric(ctx context.Context, metric internal.Metric) error {
	sql := "insert into metrics (id, mtid, value, delta, hash) values ($1, (select mtid from metric_types where mtname=$2), $3, $4, $5)"
	_, err := p.db.ExecContext(ctx, sql, metric.ID, metric.MType, metric.Value, metric.Delta, metric.Hash)
	return err
}

func (p *pgStore) GetMetric(ctx context.Context, metric internal.Metric) (internal.Metric, bool, error) {
	sql := `select t1.id, t2.mtname, t1.value, t1.delta, t1.hash from metrics t1
join metric_types t2
on t2.mtid = t1.mtid
where t1.id = $1 and t2.mtname = $2`

	row := p.db.QueryRowContext(ctx, sql, metric.ID, metric.MType)
	var ret internal.Metric
	var exist bool

	err := row.Scan(&ret.ID, &ret.MType, &ret.Value, &ret.Delta, &ret.Hash)
	if err != nil {
		return internal.Metric{}, false, err
	} else {
		exist = true
	}

	return ret, exist, nil
}

func (p *pgStore) GetAll(ctx context.Context) ([]internal.Metric, error) {
	sql := `select t1.id, t2.mtid, t1.value, t1.delta, t1.hash from metrics t1
join metric_types t2
on t2.mtid = t1.mtid`

	var ret []internal.Metric

	rows, err := p.db.QueryContext(ctx, sql)
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp internal.Metric
		err = rows.Scan(&tmp.ID, &tmp.MType, &tmp.Value, &tmp.Delta, &tmp.Hash)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
