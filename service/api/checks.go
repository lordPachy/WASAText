package api

/*import (
	"fmt"
	"reflect"
)

// Check if a user exists given its username, and returns true if it exists
func (db *appdbimpl) CheckUsername(username string) bool {
	var tmp string
	err := db.c.QueryRow("SELECT id FROM users WHERE username=" + username + ";").Scan(&tmp)
	return err == nil
}

// Check if a user exists given its identifier, and returns true if it exists
func (db *appdbimpl) CheckIdentifier(id string) bool {
	var tmp string
	fmt.Println("Identifier checking started...")
	fmt.Println("Authentication is receiving id:")
	fmt.Println(id)
	fmt.Println("of type")
	fmt.Println(reflect.TypeOf(id))
	err := db.c.QueryRow("SELECT id FROM users WHERE id = ?", id).Scan(&tmp)
	fmt.Println("Identifier checking finished...")
	fmt.Println(err.Error())
	return err == nil
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
*/
