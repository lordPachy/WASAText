package api

import (
	"fmt"
	"net/http"
	"strconv"
)

/*
This package contains function that get never called directly
but always through other functions. Thus, they get to create
their own errors, and manage the writer.
*/

// It retrieves an array of strings (that should represent a single user) from id. Each string element is a row element in the db.
func UserFromIdRetrieval(id Access_token, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "User retrieval from id"

	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id.Identifier))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking users with the same id failed", err, w, rt)
	}

	// Reading the rows
	users, err := UsersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking users with the same id failed", err, w, rt)
	}

	return users, nil
}

// It retrieves an array of strings (that should represent a single user) from id. Each string element is a row element in the db.
func UserFromUsernameRetrieval(username Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "User retrieval from username"

	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("username = '%s'", username.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking users with the same username failed", err, w, rt)
	}

	// Reading the rows
	other_users, err := UsersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking users with the same username failed", err, w, rt)
	}

	return other_users, nil
}

// It check the existence of a user.
func UserFromUsernameExists(username Username, rt *_router, w http.ResponseWriter) (bool, error) {
	user, err := UserFromUsernameRetrieval(username, rt, w)
	if err != nil {
		return false, err
	}

	if len(user) > 0 {
		return true, nil
	}

	return false, nil
}

// It check the existence of a user.
func UserFromIDExists(token Access_token, rt *_router, w http.ResponseWriter) (bool, error) {
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return false, err
	}

	if len(user) > 0 {
		return true, nil
	}

	return false, nil
}

// It checks for the existence of a message.
func MessageFromIdExistence(id int, rt *_router, w http.ResponseWriter) (bool, error) {
	chat, err := MessageFromIdRetrieval(id, rt, w)
	if err != nil {
		return false, err
	}

	if len(chat) > 0 {
		return true, nil
	}

	return false, nil
}

// It check the existence of a chat.
func ConversationFromIdExistence(id int, rt *_router, w http.ResponseWriter) (bool, error) {
	chat, err := ConversationFromIdRetrieval(id, rt, w)
	if err != nil {
		return false, err
	}

	if len(chat) > 0 {
		return true, nil
	}

	return false, nil
}

// It retrieves the list of private conversations for a user from the database. Each string element is a row element in the privchats table.
func PrivConversationsFromUsernameRetrieval(user Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Private conversations retrieval"

	// SQL query
	rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("member1 = '%s' OR member2 = '%s'", user.Name, user.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking conversations of a user failed", err, w, rt)
	}

	// Reading the rows
	chats, err := PrivchatsRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking conversations of a user failed", err, w, rt)
	}

	return chats, nil

}

// It retrieves the list of users matching a certain prefix.
func UserQuerying(user Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "User querying"

	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("username LIKE '%s%%'", user.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database querying users failed", err, w, rt)
	}

	// Reading the rows
	chats, err := UsersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were querying user failed", err, w, rt)
	}

	return chats, nil

}

// It retrieves the list of group conversations for a user from the database. Each string element is a row element in the groupmembers table.
func GroupConversationsFromUsernameRetrieval(user Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Group conversations retrieval"

	// SQL query
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("member = '%s'", user.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking groups of a user failed", err, w, rt)
	}

	// Reading the rows
	chats, err := GroupMembersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking groups of a user failed", err, w, rt)
	}

	return chats, nil
}

// It retrieves a conversation from the database. Each string element is a row element in the db.
func ConversationFromIdRetrieval(id int, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Single conversation retrieval"

	// Checking that if it is a private conversation
	if id < 5000 {
		// SQL query
		rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("id = %d", id))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking conversations with the same id failed", err, w, rt)
		}

		// Reading the rows
		chats, err := PrivchatsRowReading(rows)

		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking conversations with the same id failed", err, w, rt)
		}

		if len(chats) > 0 {
			return chats, nil
		}
	}

	// It is not a private conversation:
	// Checking if it is a groupchat conversation
	// SQL query
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = %d", id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking conversations with the same id failed", err, w, rt)
	}

	// Reading the rows
	chats, err := GroupMembersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking conversations with the same id failed", err, w, rt)
	}

	return chats, nil
}

// It retrieves a message from the database. Each string element is a row element in the message table.
func MessageFromIdRetrieval(id int, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Single message retrieval"

	// SQL query
	rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = %d", id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking messages with the same id failed", err, w, rt)
	}

	// Reading the rows
	message, err := MessageRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking messages with the same id failed", err, w, rt)
	}

	return message, nil
}

// It retrieves the information of a group from the database. Each string element is a row element in the db.
func GroupInfoFromIdRetrieval(convID ConversationID, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Group information retrieval"

	// SQL query
	rows, err := rt.db.Select("*", "groupchats", fmt.Sprintf("id = %d", convID.Id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking group infos failed", err, w, rt)
	}

	// Reading the rows
	info, err := GroupInfoRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking group infos failed", err, w, rt)
	}

	return info, nil
}

// It check the existence of a private chat.
func PrivConversationFromMembersExistence(user1 Username, user2 Username, rt *_router, w http.ResponseWriter) (bool, error) {
	user, err := PrivConversationFromMembersRetrieval(user1, user2, rt, w)
	if err != nil {
		return false, err
	}

	if len(user) > 0 {
		return true, nil
	}

	return false, nil
}

// It retrieves a private conversation from the database, given the two members. Each string element is a row element in the db.
func PrivConversationFromMembersRetrieval(user1 Username, user2 Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Private conversation retrieval, from members"

	// SQL query
	rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("member1 = '%s' AND member2 = '%s'", user1.Name, user2.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking private conversations with the same username failed", err, w, rt)
	}

	// Reading the rows
	chats, err := PrivchatsRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking private conversations with the same username failed", err, w, rt)
	}

	// They might be in the other order around
	if len(chats) == 0 {
		// SQL query
		rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("member1 = '%s' AND member2 = '%s'", user2.Name, user1.Name))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking private conversations with the same username failed", err, w, rt)
		}

		// Reading the rows
		chats, err = PrivchatsRowReading(rows)

		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking private conversations with the same username failed", err, w, rt)
		}

	}

	return chats, nil
}

// It checks if a user (either from username or ID) belongs to a given group. It may also return false if the user or group does not exist.
func UserBelongsToGroup(token Access_token, groupID ConversationID, rt *_router, w http.ResponseWriter) (bool, error) {
	// Logging information
	const affinity string = "User - group relation checking"

	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return false, err
	}

	// Taking the username from the query
	username := user[1]

	// SQL query
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = %d AND member = '%s'", groupID.Id, username))
	if err != nil {
		return false, createBackendError(affinity, "SELECT in the database seeking group-user matching failed", err, w, rt)
	}

	// Reading the rows
	chats, err := GroupMembersRowReading(rows)

	if err != nil {
		return false, createBackendError(affinity, "Reading the database rows that were seeking group-user matching failed", err, w, rt)
	}

	return len(chats) > 0, nil
}

// It returns a list of usernames of a given group. It assumes the group exists.
func UsersInGroup(groupID ConversationID, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "User of a group retrieving"

	// SQL query
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = %d", groupID.Id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking group members failed", err, w, rt)
	}

	// Reading the rows
	users, err := GroupMembersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking group members failed", err, w, rt)
	}

	// Eliminating elements of the array that are not usernames
	var usernames []string
	for i, user := range users {
		if i%2 == 1 {
			usernames = append(usernames, user)
		}
	}

	return usernames, nil
}

// It returns a list of received message ids from a group. It assumes the group exists.
func ReceivedGroupMessages(rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Received group messages retrieving"

	// Retrieving messages
	rows, err := rt.db.Filter("*", "groupmessageschecks", "messageid", "MIN(checkmarks) = 1")
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking received group messages failed", err, w, rt)
	}

	// Reading the rows
	retrievedRows, err := GroupmessageschecksRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking received group messages failed", err, w, rt)
	}

	// Eliminating elements of the array that are not message ids
	var message_ids []string
	for i, el := range retrievedRows {
		if i%4 == 1 {
			message_ids = append(message_ids, el)
		}
	}

	return message_ids, nil
}

// It returns a list of read message ids from a group. It assumes the group exists.
func ReadGroupMessages(rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Read group messages retrieving"

	// Retrieving messages
	rows, err := rt.db.Filter("*", "groupmessageschecks", "messageid", "MIN(checkmarks) = 2")
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking read group messages failed", err, w, rt)
	}

	// Reading the rows
	retrievedRows, err := GroupmessageschecksRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking read group messages failed", err, w, rt)
	}

	// Eliminating elements of the array that are not message ids
	var message_ids []string
	for i, el := range retrievedRows {
		if i%4 == 1 {
			message_ids = append(message_ids, el)
		}
	}

	return message_ids, nil
}

// It returns a list of message ids from a private conversation. It assumes the conversation exists.
func MessageIdsFromPrivateConvo(user Username, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Message ids from chat retrieving"

	// SQL query #1: retrieving chat ids where the user is present
	rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("member1 = '%s' OR member2 = '%s'", user.Name, user.Name))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking chat from username failed", err, w, rt)
	}

	// Reading the rows
	users, err := PrivchatsRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking chat ids from username failed", err, w, rt)
	}

	// Eliminating elements of the array that are not chat ids
	var ids []string
	for i, id := range users {
		if i%3 == 0 {
			ids = append(ids, id)
		}
	}

	// SQL query #2: retrieving messages sent on chat where the user is present
	var message_ids []string
	for _, id := range ids {
		rows, err := rt.db.Select("*", "privmessages", fmt.Sprintf("id = %s", id))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking message ids from chat id failed", err, w, rt)
		}

		// Reading the rows
		queriedrows, err := ChatmessagesRowReading(rows)

		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking message ids from chat id failed", err, w, rt)
		}

		// Eliminating elements of the array that are not chat ids
		for i, el := range queriedrows {
			if i%2 == 1 {
				message_ids = append(message_ids, el)
			}
		}
	}

	return message_ids, nil
}

// It returns a list of message from a private or group conversation. It assumes the conversation exists.
func MessagesFromConvo(convID ConversationID, rt *_router, w http.ResponseWriter) ([]Message, error) {
	// Logging information
	const affinity string = "Messages from chat retrieving"

	// Checking if it is a private chat or a groupchat
	var table string
	if convID.Id < 5000 {
		table = "privmessages"
	} else {
		table = "groupmessages"
	}

	// SQL query #1: retrieving message ids off the chat
	rows, err := rt.db.Select("*", table, fmt.Sprintf("id = %d", convID.Id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking chat messages failed", err, w, rt)
	}

	// Reading the rows
	rawmessages, err := ChatmessagesRowReading(rows)
	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking chat messages failed", err, w, rt)
	}

	// Eliminating elements of the array that are not chat ids
	var messageids []string
	for i, id := range rawmessages {
		if i%2 == 1 {
			messageids = append(messageids, id)
		}
	}

	// SQL query #2: retrieving messages sent on chat
	messages := make([]Message, len(messageids))
	for i, id := range messageids {
		rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = %s", id))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking messages from id failed", err, w, rt)
		}

		// Reading the rows
		queriedrows, err := MessageRowReading(rows)
		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking messages from id failed", err, w, rt)
		}

		// Converting results into the correct formats
		msgid, err := strconv.Atoi(queriedrows[0])
		if err != nil {
			return nil, createBackendError(affinity, "Message id conversion to int failed", err, w, rt)
		}

		comments, err := CommentsFromMessage(MessageID{msgid}, rt, w)
		if err != nil {
			return nil, err
		}

		checkmarks, err := strconv.Atoi(queriedrows[5])
		if err != nil {
			return nil, createBackendError(affinity, "Checkmarks conversion to int failed", err, w, rt)
		}

		var replyingid int
		if queriedrows[6] != nullValue {
			replyingid, err = strconv.Atoi(queriedrows[6])
			if err != nil {
				return nil, createBackendError(affinity, "Message replyed to id conversion to int failed", err, w, rt)
			}
		} else {
			replyingid = -1
		}

		// Packing everything into a message
		tmpMessage := Message{
			MessageID:  msgid,
			Timestamp:  queriedrows[2],
			Content:    queriedrows[3],
			Photo:      queriedrows[4],
			Username:   queriedrows[1],
			Checkmarks: checkmarks,
			Comments:   comments,
			ReplyingTo: replyingid,
			OG_Sender:  queriedrows[7],
		}

		messages[i] = tmpMessage
	}

	return messages, nil
}

// It returns a list of message from a private or group conversation, with at most one element, being the last. It assumes the conversation exists.
func LastMessageFromConvo(convID ConversationID, rt *_router, w http.ResponseWriter) ([]Message, error) {
	// Logging information
	const affinity string = "Messages from chat retrieving"

	// Checking if it is a private chat or a groupchat
	var table string
	if convID.Id < 5000 {
		table = "privmessages"
	} else {
		table = "groupmessages"
	}

	// SQL query #1: retrieving message ids off the chat
	rows, err := rt.db.Select("*", table, fmt.Sprintf("id = %d", convID.Id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking last chat messages failed", err, w, rt)
	}

	// Reading the rows
	rawmessages, err := ChatmessagesRowReading(rows)
	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking last chat messages failed", err, w, rt)
	}

	// Eliminating elements of the array that are not chat ids
	var messageids []string
	for i, id := range rawmessages {
		if i%2 == 1 {
			messageids = append(messageids, id)
		}
	}

	// SQL query #2: retrieving last message sent on chat
	messages := make([]Message, 0, 1)
	for _, id := range messageids {
		rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = %s", id))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking last message from id failed", err, w, rt)
		}

		// Reading the rows
		queriedrows, err := MessageRowReading(rows)
		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking last message from id failed", err, w, rt)
		}

		if len(messages) > 0 {
			if messages[0].Timestamp > queriedrows[2] {
				continue
			}
		}

		// Creating the result message
		var emptyComments []Comment

		// Converting results into the correct formats
		msgid, err := strconv.Atoi(queriedrows[0])
		if err != nil {
			return nil, createBackendError(affinity, "Message id conversion to int failed", err, w, rt)
		}
		checkmarks, err := strconv.Atoi(queriedrows[5])
		if err != nil {
			return nil, createBackendError(affinity, "Checkmarks conversion to int failed", err, w, rt)
		}

		var replyingid int
		if queriedrows[6] != nullValue {
			replyingid, err = strconv.Atoi(queriedrows[6])
			if err != nil {
				return nil, createBackendError(affinity, "Message replyed to id conversion to int failed", err, w, rt)
			}
		} else {
			replyingid = -1
		}

		// Packing everything into a message
		tmpMessage := Message{
			MessageID:  msgid,
			Timestamp:  queriedrows[2],
			Content:    queriedrows[3],
			Photo:      queriedrows[4],
			Username:   queriedrows[1],
			Checkmarks: checkmarks,
			Comments:   emptyComments,
			ReplyingTo: replyingid,
			OG_Sender:  queriedrows[7],
		}

		if len(messages) > 0 {
			messages[0] = tmpMessage
		} else {
			messages = append(messages, tmpMessage)
		}
	}

	return messages, nil
}

// It returns a list of comments from a message. It assumes the message exists.
func CommentsFromMessage(messID MessageID, rt *_router, w http.ResponseWriter) ([]Comment, error) {
	// Logging information
	const affinity string = "Comments from message retrieving"

	// SQL query #1: retrieving comment ids off the message
	rows, err := rt.db.Select("*", "messagecomments", fmt.Sprintf("id = %d", messID.Id))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking message comments failed", err, w, rt)
	}

	// Reading the rows
	rawcomments, err := MessageCommentsRowReading(rows)
	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking message comments failed", err, w, rt)
	}

	// Eliminating elements of the array that are not comment ids
	var commentids []string
	for i, id := range rawcomments {
		if i%2 == 1 {
			commentids = append(commentids, id)
		}
	}

	// SQL query #2: retrieving comments
	comments := make([]Comment, len(commentids))
	for i, id := range commentids {
		rows, err := rt.db.Select("*", "comments", fmt.Sprintf("id = %s", id))
		if err != nil {
			return nil, createBackendError(affinity, "SELECT in the database seeking comments from id failed", err, w, rt)
		}

		// Reading the rows
		queriedrows, err := CommentsRowReading(rows)
		if err != nil {
			return nil, createBackendError(affinity, "Reading the database rows that were seeking comments from id failed", err, w, rt)
		}

		// Creating the result message

		// Converting results into the correct formats
		cmtid, err := strconv.Atoi(queriedrows[0])
		if err != nil {
			return nil, createBackendError(affinity, "Comment id conversion to int failed", err, w, rt)
		}

		// Packing everything into a comment
		tmpComment := Comment{
			CommentID: cmtid,
			Sender:    queriedrows[1],
			Reaction:  queriedrows[2],
		}

		comments[i] = tmpComment
	}

	return comments, nil
}

// It retrieves a comment from the database. Each string element is a row element in the comments table.
func CommentFromIdRetrieval(commentID CommentID, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Single comment retrieval"

	// SQL query
	rows, err := rt.db.Select("*", "comments", fmt.Sprintf("id = %d", commentID.CommentID))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking comments with the same id failed", err, w, rt)
	}

	// Reading the rows
	comment, err := CommentsRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking comments with the same id failed", err, w, rt)
	}

	return comment, nil
}
