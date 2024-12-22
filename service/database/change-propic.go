package database

// It updates a user's profile picture
func (db *appdbimpl) ChangeProPic(id string, newPhoto string) error {
	_, err := db.c.Exec("UPDATE users SET propic = ? WHERE id =?;", newPhoto, id)
	return err
}
