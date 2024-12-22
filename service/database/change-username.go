package database

// It substitutes a username with a different one
func (db *appdbimpl) ChangeUsername(id string, newUserName string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE id =?;", newUserName, id)
	return err
}
