package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database interface {
	OpenConnection()
}

type database struct {
	config Config
}

func NewDatabase(config Config) Database {
	return &database{config}
}

func (d *database) OpenConnection() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		d.config.User,
		d.config.Password,
		d.config.Host,
		d.config.Port,
		d.config.Database,
	)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
