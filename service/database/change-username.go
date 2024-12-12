package database

// It substitutes a username with a different one
func (db *appdbimpl) ChangeUsername(oldUsername string, newUserName string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE username =?;", newUserName, oldUsername)
	return err
}
