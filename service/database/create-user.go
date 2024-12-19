package database

import (
	"github.com/lucasjones/reggen"
)

// It creates a new username, by assigning to it a random id
func (db *appdbimpl) CreateUser(newUsername string) error {
	var id string
	condition := true
	for condition {
		id, _ = reggen.Generate("^[0-9A-Za-z-$!&/(),.]{4,16}$", 16)
		var ret string
		err := db.c.QueryRow("SELECT id FROM users WHERE id = ?;", id).Scan(&ret)
		// If the previous query results in an error,
		// it means that an user with such ID already exists
		if err != nil {
			condition = false
		}
	}
	user := "('" + id + "', '" + newUsername + "')"
	_, err := db.c.Exec("INSERT INTO users (id, username) VALUES " + user + ";")
	return err
}
