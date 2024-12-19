package database

import (
	"database/sql"
	"fmt"
)

// Check user returns the userID associated with the username if this exists; otherwise, it throws an error
func (db *appdbimpl) GetIdentifier(username string) (string, error) {
	var userID string
	err := db.c.QueryRow("SELECT id FROM users WHERE username=?;", username).Scan(&userID)
	// Debugging
	var rows *sql.Rows
	var values string
	rows, _ = db.c.Query("SELECT * FROM users;")
	rows.Scan(values)
	fmt.Println(values)

	return userID, err
}
