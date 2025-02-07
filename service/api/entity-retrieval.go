package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It checks if a single conversation exists.
func conversationRetrieval(affinity string, w http.ResponseWriter, ps httprouter.Params, rt *_router) (bool, error) {
	// Checking that the conversation actually exists
	convID, err := strconv.Atoi(ps.ByName("conversationid"))
	if err != nil {
		return false, createBackendError(affinity, "Conversation id conversion failed", err, w, rt)
	}

	exists, err := ConversationFromIdExistence(convID, rt, w)
	if err != nil {
		return false, err
	}
	if !exists {
		createFaultyResponse(http.StatusNotFound, "Conversation not found", affinity, "Response message encoding for conversation not found error has failed", w, rt)
		return false, nil
	}

	return true, nil
}
