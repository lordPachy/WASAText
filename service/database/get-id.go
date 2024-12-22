package database

// Check user returns the userID associated with the username if this exists; otherwise, it throws an error
func (db *appdbimpl) GetIdentifier(username string) (string, error) {
	var userID string
	err := db.c.QueryRow("SELECT id FROM users WHERE username=?;", username).Scan(&userID)

	return userID, err
}
