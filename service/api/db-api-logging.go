package api

import (
	"fmt"
	"net/http"
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
func MessageFromIdExists(id int, rt *_router, w http.ResponseWriter) (bool, error) {
	// Logging information
	const affinity string = "Message existence checking"

	// Querying database rows
	rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = '%d'", id))
	if err != nil {
		return false, createBackendError(affinity, "SELECT in the database seeking messages with the same id failed", err, w, rt)
	}

	// Checking the queried rows
	other_messages, err := MessageRowReading(rows)

	if err != nil {
		return false, createBackendError(affinity, "Reading the database rows that were seeking messages with the same id failed", err, w, rt)
	} else if len(other_messages) == 0 {
		return false, nil
	}

	return true, nil
}

// // It check the existence of a chat.
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

// It retrieves a conversation from the database. Each string element is a row element in the db.
func ConversationFromIdRetrieval(id int, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	const affinity string = "Single conversation retrieval"

	// Checking that if it is a private conversation
	if id < 5000 {
		// SQL query
		rows, err := rt.db.Select("*", "privchats", fmt.Sprintf("id = '%d'", id))
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
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = '%d'", id))
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
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = '%d' AND member = '%s'", groupID.Id, username))
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
	rows, err := rt.db.Select("*", "groupmembers", fmt.Sprintf("id = '%d'", groupID.Id))
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
