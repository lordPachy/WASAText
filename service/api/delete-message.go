package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It deletes a message.
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Message deleting"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Checking that the conversation actually exists
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

	// Checking that the message has been sent from the user
	messInfo, err := MessageFromIdRetrieval(messID, rt, w)
	if err != nil {
		return
	}

	if messInfo[1] != user[1] {
		createFaultyResponse(http.StatusForbidden, "The message has not been sent by this user, and cannot be deleted", affinity, "Response for missed paternity encoding has failed", w, rt)
		return
	}

	// Deleting message from the DB
	query := fmt.Sprintf("id = %d", messID)
	rows, err := rt.db.Delete("messages", query)
	if err != nil || rows.Err() != nil {
		_ = createBackendError(affinity, "Deleting message from the database has failed", err, w, rt)
		return
	}

	// Completing deletion
	_, err = MessageRowReading(rows)
	if err != nil {
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}
