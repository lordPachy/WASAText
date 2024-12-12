package database

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetName(oldUsername string, newUserName string) error {
	_, err := db.c.Exec("INSERT INTO users (id, name) VALUES (?, ?)", newUserName)
	return err
}
