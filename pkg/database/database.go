package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kjuulh/rawpotion-go/pkg/config"
	_ "github.com/lib/pq"
)

type Database struct {
	config Config
	Db     *sql.DB
}

func NewDatabase() Database {
	return Database{}
}

func (d *Database) Open() {
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
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	d.Db = db
}

func (d *Database) Close() (err error) {
	if d.Db == nil {
		return
	}

	err = d.Db.Close()
	return
}

func (d *Database) LoadConfigFromFile(path string) {
	cfg, err := config.GetConfigFromFile(path)
	if err != nil {
		log.Fatal("Couldn't read file")
		panic(err)
	}

	d.config = Config{
		Host:     cfg.Database.Host,
		Database: cfg.Database.Database,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Port:     cfg.Database.Port,
	}
}
