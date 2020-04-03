package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB       *sqlx.DB
	prepared map[string]*sqlx.Stmt
}

func (db *DB) prepare(key, query string) error {
	stmt, err := db.DB.Preparex(query)
	if err != nil {
		return err
	}

	db.prepared[key] = stmt
	return nil
}

func InitDB() (*DB, error) {
	dbConn, err := sqlx.Connect("postgres", "user=root dbname=corona host='hackthecrisis.c8e6qltbhtn7.eu-north-1.rds.amazonaws.com' password='XKerP8P73XtL9w2'")
	dbb := DB{DB: dbConn, prepared: make(map[string]*sqlx.Stmt)}

	err = prepareQueries(&dbb)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return &dbb, nil
}

func prepareQueries(db *DB) error {
	for _, query := range queries {
		err := db.prepare(query.K, query.V)
		if err != nil {
			return err
		}
	}
	return nil
}
