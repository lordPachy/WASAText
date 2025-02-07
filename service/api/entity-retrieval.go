package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It checks if a single conversation exists.
func conversationRetrieval(affinity string, w http.ResponseWriter, ps httprouter.Params, rt *_router) (int, error) {
	// Checking that the conversation actually exists
	convID, err := strconv.Atoi(ps.ByName("conversationid"))
	if err != nil {
		return convID, createBackendError(affinity, "Conversation id conversion failed", err, w, rt)
	}

	exists, err := ConversationFromIdExistence(convID, rt, w)
	if err != nil {
		return convID, err
	}
	if !exists {
		createFaultyResponse(http.StatusNotFound, "Conversation not found", affinity, "Response message encoding for conversation not found error has failed", w, rt)
		return -1, nil
	}

	return convID, nil
}
