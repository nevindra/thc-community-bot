package setup

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func DB() (*sql.DB, error) {
	dbpath := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dbpath)
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
