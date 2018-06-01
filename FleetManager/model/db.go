package model

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDb method unuts postgres connection
func InitDb(dsn string) {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
