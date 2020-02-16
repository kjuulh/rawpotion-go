package tables

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/kjuulh/rawpotion-go/pkg/database"
)

type UsersTable interface {
	Insert(row UsersRow) (newRow UsersRow, err error)
	GetByUsername(username string) (row UsersRow, err error)
	GetAll() (rows []UsersRow, err error)
}

type usersDb struct {
	db *database.Database
}

type UsersRow struct {
	Id       string
	Username string
	Password string
}

func NewUsersTable(db *database.Database) (table usersDb, err error) {
	if db == nil {
		err = errors.New("Cannot create UsersTable without UsersTableConfig")
		return
	}

	table.db = db

	if err = table.createTable(); err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (table *usersDb) createTable() (err error) {
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

func (table *usersDb) Insert(row UsersRow) (newRow UsersRow, err error) {
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

func (table *usersDb) GetByUsername(username string) (row UsersRow, err error) {
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

func (table *usersDb) GetAll() (rows []UsersRow, err error) {
	const qry = `
		SELECT id, username FROM users
	`

	rs, err := table.db.Db.Query(qry)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rs.Close()

	for rs.Next() {
		u := UsersRow{}
		err = rs.Scan(&u.Id, &u.Username)
		if err != nil {

			return
		}
		rows = append(rows, u)

	}
	if !rs.NextResultSet() {
		err = rs.Err()
	}

	return
}

var User UsersTable

func InitUsersTable(db *database.Database) {
	User = &usersDb{db: db}
	_, err := NewUsersTable(db)
	if err != nil {
		panic(err)
	}
}
