package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It adds a member to group.
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Group leaving"

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

	// If the current conversation is not a group, it cannot be left
	if convID < 5000 {
		createFaultyResponse(http.StatusUnauthorized, "The current conversation is not a group", affinity, "Response for current conversation isnot a group encoding failed", w, rt)
		return
	}

	// Retrieving the username
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	// Deleting membership from the DB
	query := fmt.Sprintf("id = '%d' AND member = '%s'", convID, user[1])

	rows, err := rt.db.Delete("groupmembers", query)
	if err != nil || rows.Err() != nil {
		_ = createBackendError(affinity, "Deleting group membership from the database has failed", err, w, rt)
		return
	}

	_, err = GroupMembersRowReading(rows)
	if err != nil {
		_ = createBackendError(affinity, "Deleting group membership from the database has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}
