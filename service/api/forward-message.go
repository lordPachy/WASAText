package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It modifies the profile picture of an existing user.
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Message forwarding"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Checking that the current conversation actually exists
	convID, err := conversationRetrieval(affinity, w, ps, rt)
	if err != nil || convID == -1 {
		return
	}

	// Checking that the message actually exists
	messID, err := messageRetrieval(affinity, w, ps, rt)
	if err != nil || messID == -1 {
		return
	}

	// Retrieving the username
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	// Getting message info
	messInfo, err := MessageFromIdRetrieval(messID, rt, w)
	if err != nil {
		return
	}

	// Getting the new conversation id
	var newConversation ConversationID
	err = json.NewDecoder(r.Body).Decode(&newConversation)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a conversation id", affinity, "Request encoding for badly formatted message response has failed", w, rt)
		return
	}

	// Checking that the new conversation exists
	exists, err := ConversationFromIdExistence(newConversation.Id, rt, w)
	if err != nil {
		return
	}
	if !exists {
		createFaultyResponse(http.StatusNotFound, "The conversation forwarding to does not exist", affinity, "Response encoding for not existing new conversation has failed", w, rt)
		return
	}

	// Creating the new id
	id, err := MessageIdCreator(rt, w)
	if err != nil {
		return
	}

	// Getting the new timestamp
	timestamp := GetTime()

	// Actually writing the message in the DB
	query := fmt.Sprintf("(%d, '%s', '%s', '%s', '%s', %d, %s, %s)", id, user[1], timestamp, messInfo[3], messInfo[4], 0, nullValue, messInfo[1])

	_, err = rt.db.Insert("messages", query)
	if err != nil {
		_ = createBackendError(affinity, "Inserting the new message into the database has failed", err, w, rt)
		return
	}

	// Writing the message into the message recording tables
	err = update_messages(ConversationID{Id: newConversation.Id}, MessageID{Id: id}, w, rt)
	if err != nil {
		return
	}

	if newConversation.Id >= 5000 {
		// If it is a group message, we must update the checkmarks table
		err = groupMessageCheckmarksUpdate(ConversationID{newConversation.Id}, MessageID{id}, Username{user[1]}, w, rt)
		if err != nil {
			return
		}
	}
	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}
