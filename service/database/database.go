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
	"strings"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Insert(table string, values string) (sql.Result, error)
	Update(table string, update string, condition string) (sql.Result, error)
	Select(columns string, table string, conditions string) (*sql.Rows, error)
	Filter(columns string, table string, group_by string, conditions string) (*sql.Rows, error)
	Delete(table string, conditions string) (*sql.Rows, error)
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
	_ = db.QueryRow(`SELECT username FROM users;`).Scan(&tableName)
	// If not, the database is empty, and we need to create the structure
	// if err.Error() != "sql: no rows in result set" {
	if true {
		file := "-- Cleaning up the database\nPRAGMA writable_schema = 1;\nDELETE FROM sqlite_master;\nPRAGMA writable_schema = 0;\nVACUUM;\nPRAGMA integrity_check;\n\n\n-- Creating the actual schema\nCREATE TABLE users (id TEXT PRIMARY KEY, username TEXT UNIQUE NOT NULL, propic TEXT);\nCREATE TABLE privchats (id INTEGER PRIMARY KEY, member1 TEXT NOT NULL, member2 TEXT NOT NULL, FOREIGN KEY (member1) REFERENCES users(username) ON UPDATE CASCADE, FOREIGN KEY (member2) REFERENCES users(username) ON UPDATE CASCADE);\nCREATE TABLE privmessages (id INTEGER, messageID INTEGER, PRIMARY KEY (id, messageID), FOREIGN KEY (id) REFERENCES privchats(id), FOREIGN KEY (messageID) REFERENCES messages(id) ON DELETE CASCADE);\nCREATE TABLE groupchats (id INTEGER PRIMARY KEY, groupname TEXT NOT NULL, groupphoto TEXT);\nCREATE TABLE groupmembers (id INTEGER, member TEXT, PRIMARY KEY (id, member), FOREIGN KEY (member) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE, FOREIGN KEY (id) REFERENCES groupchats(id) ON DELETE CASCADE);\nCREATE TABLE groupmessages (id INTEGER, messageID INTEGER, PRIMARY KEY (id, messageID), FOREIGN KEY (id) REFERENCES groupchats(id) ON DELETE CASCADE, FOREIGN KEY (messageID) REFERENCES messages(id) ON DELETE CASCADE);\nCREATE TABLE messages (id INTEGER PRIMARY KEY, sender TEXT, created_at TIMESTAMP NOT NULL, content TEXT, photo TEXT, checkmarks INTEGER NOT NULL, replying_to INTEGER, og_sender TEXT, FOREIGN KEY (sender) REFERENCES users(username) ON UPDATE CASCADE, FOREIGN KEY (replying_to) REFERENCES messages(id) ON DELETE SET NULL, FOREIGN KEY (og_sender) REFERENCES users(username) ON UPDATE CASCADE);\nCREATE TABLE messagecomments (id INTEGER, commentID INTEGER, PRIMARY KEY(id, commentID), FOREIGN KEY (id) REFERENCES messages(id) ON DELETE CASCADE, FOREIGN KEY (commentID) REFERENCES comments(id) ON DELETE CASCADE);\nCREATE TABLE comments (id INTEGER PRIMARY KEY, sender TEXT, reaction TEXT NOT NULL, FOREIGN KEY (sender) REFERENCES users(username) ON UPDATE CASCADE);\nCREATE TABLE groupmessageschecks (groupID INTEGER, messageID INTEGER, member TEXT, checkmarks INTEGER NOT NULL, PRIMARY KEY (groupID, messageID, member), FOREIGN KEY (groupID, member) REFERENCES groupmembers(id, member) ON UPDATE CASCADE ON DELETE CASCADE, FOREIGN KEY (groupID, messageID) REFERENCES groupmessages(id, messageID) ON DELETE CASCADE);"
		requests := strings.Split(string(file), "\n")
		for _, request := range requests {
			// fmt.Println(request)
			_, err := db.Exec(request)
			if err != nil {
				// fmt.Println(err.Error())
				continue
			}
		}
	}

	// fmt.Println("Database creation: completed")

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
