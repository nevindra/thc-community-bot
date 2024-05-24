package setup

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "thc.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Printf("DB is successfully connected!")

	return db, nil
}
