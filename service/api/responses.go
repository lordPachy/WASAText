package api

import (
	"encoding/json"
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

func (rt *_router) authorization(w http.ResponseWriter, r *http.Request) (Username, error) {
	var user Username
	id := r.Header.Get("Authorization")

	// If the id is not in the database, throw an unauthorized error; otherwise return the username
	name, err := rt.db.GetIdentifier(id)
	if err != nil {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
		return user, err
	} else {
		user.Name = name
		return user, err
	}
}
