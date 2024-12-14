package database

// A function that sends a message. This function does not even return error, since error-checking is
// assumed to have been already brought out
func (db *appdbimpl) SendMessage(newMessage Message) {
	var message_query string
	message_query = "('" + string(newMessage.MessageID) + "', '" + newMessage.Username + "', '" + newMessage.Timestamp + "', '" + newMessage.Content + "', '" + newMessage.Photo + "', '" + string(newMessage.Checkmarks) + "', '" + newMessage.ReplyingTo + "');"

	_, _ = db.c.Exec("INSERT INTO messages VALUES" + message_query)
}
