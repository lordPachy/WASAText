package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It modifies the name of an existing group.
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Group name setting"

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
		createFaultyResponse(http.StatusBadRequest, "The current conversation is not a group", affinity, "Response for current conversation is not a group encoding failed", w, rt)
		return
	}

	// Checking if the user is in the group
	belonging, err := UserBelongsToGroup(token, ConversationID{convID}, rt, w)
	if err != nil {
		return
	}

	if !belonging {
		createFaultyResponse(http.StatusUnauthorized, "The user does not belong to the said group", affinity, "Request for user not belonging to group encoding failed", w, rt)
		return
	}

	// Getting the new groupname
	var newName GroupName
	err = json.NewDecoder(r.Body).Decode(&newName)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a group name", affinity, "Request encoding for badly formatted group name has failed", w, rt)
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newName.Value)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for group name checking has failed", err, w, rt)
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "The group name is not valid (it may be too short, or long, or containing not valid characters)", affinity, "Request encoding for group name not matching with regex response has failed", w, rt)
		return
	}

	// Actually setting the username in the DB
	_, err = rt.db.Update("groupchats", fmt.Sprintf("groupname = '%s'", newName.Value), fmt.Sprintf("id = %d", convID))
	if err != nil {
		_ = createBackendError(affinity, "Updating group with the new name has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
