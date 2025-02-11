package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It modifies the photo of an existing group.
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Group photo setting"

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
		createFaultyResponse(http.StatusForbidden, "The user does not belong to the said group", affinity, "Request for user not belonging to group encoding failed", w, rt)
		return
	}

	// Getting the new groupphoto
	var newPhoto Image
	err = json.NewDecoder(r.Body).Decode(&newPhoto)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a group photo", affinity, "Request encoding for badly formatted group photo has failed", w, rt)
		return
	}

	// Actually setting the username in the DB
	_, err = rt.db.Update("groupchats", fmt.Sprintf("groupphoto = '%s'", newPhoto.Image), fmt.Sprintf("id = %d", convID))
	if err != nil {
		_ = createBackendError(affinity, "Updating group with the new photo has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
