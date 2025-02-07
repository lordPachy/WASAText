package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It deletes a comment.
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Comment deleting"

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

	// Retrieving the comment
	comment, err := commentRetrieval(affinity, w, ps, rt)
	if err != nil || len(comment) == 0 {
		return
	}

	// Checking that the comment has been posted from the user
	if comment[1] != user[1] {
		createFaultyResponse(http.StatusForbidden, "The comment has not been sent by this user, and cannot be deleted", affinity, "Response for missed paternity encoding has failed", w, rt)
		return
	}

	// Deleting comment from the DB
	query := fmt.Sprintf("id = %s", comment[0])

	rows, err := rt.db.Delete("comments", query)
	if err != nil || rows.Err() != nil {
		_ = createBackendError(affinity, "Deleting comment from the database has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}
