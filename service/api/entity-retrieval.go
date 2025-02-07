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

// It checks if a single message exists.
func messageRetrieval(affinity string, w http.ResponseWriter, ps httprouter.Params, rt *_router) (int, error) {
	// Checking that the message actually exists
	messID, err := strconv.Atoi(ps.ByName("messageid"))
	if err != nil {
		return messID, createBackendError(affinity, "Message id conversion failed", err, w, rt)
	}

	exists, err := MessageFromIdExistence(messID, rt, w)
	if err != nil {
		return messID, err
	}
	if !exists {
		createFaultyResponse(http.StatusNotFound, "Message not found", affinity, "Response message encoding for message not found error has failed", w, rt)
		return -1, nil
	}

	return messID, nil
}

// It checks if a single comment exists.
func commentRetrieval(affinity string, w http.ResponseWriter, ps httprouter.Params, rt *_router) ([]string, error) {
	// Checking that the comment actually exists
	commID, err := strconv.Atoi(ps.ByName("commentid"))
	if err != nil {
		return nil, createBackendError(affinity, "Comment id conversion failed", err, w, rt)
	}

	exists, err := CommentFromIdRetrieval(CommentID{commID}, rt, w)
	if err != nil {
		return nil, err
	}
	if len(exists) == 0 {
		createFaultyResponse(http.StatusNotFound, "Comment not found", affinity, "Response message encoding for comment not found error has failed", w, rt)
		return []string{}, nil
	}

	return exists, nil
}
