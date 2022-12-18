package version1

import (
	"database/sql"
	"exercise/geektime/homework4/version1/model"
)

type DBOption func(*DB)

type DB struct {
	r  model.Regist
	db *sql.DB
	//valCreator valuer.Creator
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:  model.Newregistry(),
		db: db,
		//valCreator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
