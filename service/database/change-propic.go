package database

// It updates a user's profile picture
func (db *appdbimpl) ChangeProPic(username string, newPhoto string) error {
	_, err := db.c.Exec("UPDATE users SET propic = ? WHERE username =?;", newPhoto, username)
	return err
}
