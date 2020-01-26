package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Database interface {
	OpenConnection()
}

type database struct{}

func NewDatabase() Database {
	return &database{}
}

func (d *database) OpenConnection() {
	connStr := "user=docker password=docker host=postgres dbname=docker sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
