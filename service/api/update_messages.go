package api

import (
	"fmt"
	"net/http"
)

/*
It uploads a message on a private chat or a group, depending on a flag.
Since this function gets called only in other functions, it is assumed
that the sender actually belongs to the chat/group, that the chat/group
exists, and that the id message is valid and corresponds to a real message.
*/
func update_messages(convID ConversationID, messageID MessageID, w http.ResponseWriter, rt *_router) error {
	// Logging information
	const affinity string = "Message table updating"

	// Actually writing the conversation in the DB
	if convID.Id < 5000 {
		query := fmt.Sprintf("(%d, '%d')", convID.Id, messageID.Id)

		_, err := rt.db.Insert("privmessages", query)
		if err != nil {
			return createBackendError(affinity, "Inserting the new message into the database of messages per chat has failed", err, w, rt)
		}
	} else {
		// Writing in the database of groups
		query := fmt.Sprintf("(%d, '%d')", convID.Id, messageID.Id)

		_, err := rt.db.Insert("groupmessages", query)
		if err != nil {
			return createBackendError(affinity, "Inserting the new message into the database of messages per group has failed", err, w, rt)
		}
	}

	return nil
}
