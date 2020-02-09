package tables

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/kjuulh/rawpotion-go/pkg/database"
)

type UsersTable struct {
	db *database.Database
}

type UsersTableConfig struct {
	Db *database.Database
}

type UsersRow struct {
	Id       string
	Username string
	Password string
}

func NewUsersTable(cfg UsersTableConfig) (table UsersTable, err error) {
	if cfg.Db == nil {
		err = errors.New("Cannot create UsersTable without UsersTableConfig")
		return
	}

	table.db = cfg.Db

	if err = table.createTable(); err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (table *UsersTable) createTable() (err error) {
	const qry = `
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
			CREATE TABLE IF NOT EXISTS users (
				id uuid DEFAULT uuid_generate_v4(),
				username text NOT NULL UNIQUE,
				password text NOT NULL
		)
	`

	if _, err = table.db.Db.Exec(qry); err != nil {
		return
	}

	return
}

func (table *UsersTable) InsertUser(row UsersRow) (newRow UsersRow, err error) {
	if row.Username == "" || row.Password == "" {
		err = errors.New("Can't create user without username and password")
		return
	}

	const qry = `
INSERT INTO users (
			username,
			password
	)
VALUES (
			$1,
			$2
)
RETURNING
	id, username, password
	`

	err = table.db.Db.
		QueryRow(qry, row.Username, row.Password).
		Scan(&newRow.Id, &newRow.Username, &newRow.Password)
	if err != nil {
		return
	}
	return
}

func (table *UsersTable) GetUserByUsername(username string) (row UsersRow, err error) {
	if username == "" {
		err = errors.New("Cannot get empty username")
		return
	}

	const qry = `
	SELECT *
	FROM users
	WHERE username = $1
	`

	err = table.db.Db.QueryRow(qry, username).Scan(&row.Id, &row.Username, &row.Password)
	if err != nil {
		return
	}

	return
}
