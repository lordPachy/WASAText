package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var UnauthorizedError = Response{
	Code:    401,
	Message: "UnauthorizedError: Access token is missing or invalid.",
}

var ConversationNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation id given was not found.",
}

var ConversationOrMessageNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation/Message id given was not found.",
}

var ConversationOrMessageOrCommentNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation/Message/Comment id given was not found.",
}

func (rt *_router) authorization(w http.ResponseWriter, r *http.Request) (string, error) {
	id := r.Header.Get("Authentication")

	// If the id is not in the database, throw an unauthorized error; otherwise return the username
	check := rt.db.CheckIdentifier(id)
	if !check {
		fmt.Println("I'm in an error state!")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
		return "", errors.New("Unauthorized")
	}

	return id, nil
}
