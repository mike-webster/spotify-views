package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	Exec(context.Context, string) (sql.Result, error)
}

type LiveDB struct {
	db *sqlx.DB
}

func GetLiveDB(conn string) (*LiveDB, error) {
	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		return nil, err
	}
	return &LiveDB{db: db}, nil
}

func (db *LiveDB) Exec(ctx context.Context, sql string) (sql.Result, error) {
	conn, err := db.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	res, err := conn.ExecContext(ctx, sql)
	return res, err
}

type TestDB struct {
	shouldExecErr bool
	execResult    *sql.Result
}

func (db *TestDB) Exec(ctx context.Context, sql string) (sql.Result, error) {
	if db.shouldExecErr {
		return nil, errors.New("test error")
	}

	res := db.execResult
	if res == nil {
		return nil, errors.New("no result from sql")
	}

	return *res, nil
}
