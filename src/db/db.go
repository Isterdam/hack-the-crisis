package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	dbConn, err := sqlx.Connect("postgres", "user="+os.Getenv("DBUSER")+" dbname="+os.Getenv("DBDB")+" host="+os.Getenv("DBHOST")+" "+os.Getenv("DBOPTS")+" password="+os.Getenv("DBPASS"))
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(dbConn.DB, &postgres.Config{})

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)

	if err != nil {
		return nil, err
	}

	err = m.Up()

	if err != nil {
		return nil, err
	}

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
			fmt.Printf("%s", err)
			return err
		}
	}
	return nil
}
