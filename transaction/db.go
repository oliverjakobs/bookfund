package transaction

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDB(name string) error {
	var err error
	db, err = sql.Open("sqlite3", name)

	return err
}

func CloseDB() {
	db.Close()
}
