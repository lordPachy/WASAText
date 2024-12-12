/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetIdentifier(userID string) (string, error)
	SetName(oldUsername string, newUserName string) error
	ChangeUsername(oldUsername string, newUserName string) error
	ChangeProPic(username string, newPhoto string) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT username FROM users;`).Scan(&tableName)
	// If not, the database is empty, and we need to create the structure
	if err.Error() != "sql: no rows in result set" {
		file, err := os.ReadFile("./service/database/db_init.sql")
		if err != nil {
			return nil, fmt.Errorf("error opening the sql file for database creation")
		}
		requests := strings.Split(string(file), "\n")
		for _, request := range requests {
			fmt.Println(request)
			_, err := db.Exec(request)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
	}

	fmt.Println("Database creation: completed")
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
