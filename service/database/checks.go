package database

// Check if a user exists, and returns a boolean accordingly
func (db *appdbimpl) CheckUsername(username string) bool {
	var tmp string
	err := db.c.QueryRow("SELECT id FROM users WHERE username=" + username + ";").Scan(&tmp)
	return err != nil
}

// Check if a message exists, and returns a boolean accordingly
func (db *appdbimpl) CheckMessage(messageid string) bool {
	var tmp string
	err := db.c.QueryRow("SELECT id FROM messages WHERE id=" + messageid + ";").Scan(&tmp)
	return err != nil
}

// Check if a conversation exists, and returns a boolean accordingly
func (db *appdbimpl) CheckChat(conversationid string, username string) bool {
	tmp := ""

	// Case 1: the conversation is private
	err := db.c.QueryRow("SELECT id FROM privchats WHERE id=" + conversationid + ";").Scan(&tmp)

	// Case 2: it is a group
	if err != nil {
		rows, err_query := db.c.Query("SELECT id FROM groupmembers WHERE member=" + username + ";")
		if err_query != nil {
			return false
		}

		defer rows.Close()

		for rows.Next() {
			_ = rows.Scan(&tmp)
		}
	} else {
		return true
	}

	return tmp != ""
}
