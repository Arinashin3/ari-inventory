package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	driver string
	name   string
	conn   *sql.DB
}

var db *Database

func NewDatabase(driver string, name string) error {
	var err error
	db = &Database{driver: driver, name: name}
	db.conn, err = sql.Open(driver, name)
	defer db.conn.Close()
	if err != nil {
		return err
	}

	schema := `
PRAGMA foreign_keys = ON;
       
CREATE TABLE IF NOT EXISTS job (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	uuid VARCHAR(30) UNIQUE NOT NULL,
    hostname VARCHAR(30) NOT NULL,
    host_group VARCHAR(30) NOT NULL DEFAULT 'none',
	ipaddr VARCHAR(30) NOT NULL,
	port INTEGER NOT NULL,
	job VARCHAR(30) NOT NULL,
    joinAt DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updateAt DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	interval VARCHAR(10) NOT NULL
);

`

	_, err = db.conn.Exec(schema)

	return nil
}
func GetDatabase() *Database {
	return db
}

func (d *Database) Open() error {
	var err error
	d.conn, err = sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	_, err = db.conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	return err
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.conn.Query(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.conn.Exec(query, args...)
}

func (d *Database) Close() {
	d.conn.Close()
}
