package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It sends a message to a specified chat.
// Implementation notes:
// 1. Creating the message
// 2. Adding to the message database
// 3. Adding to the chat database
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Authentication
	username, err := rt.authorization(w, r)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
		return
	}

	// Checking whether the conversation exists
	convo := ps.ByName("conversationid")
	convo_check := rt.db.CheckChat(convo, username.Name)

	if !convo_check {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		notFoundError := Response{
			Code:    404,
			Message: "Conversation not found",
		}
		json.NewEncoder(w).Encode(notFoundError)
		return
	}

	// Reading the message
	var newMessage Message
	err = json.NewDecoder(r.Body).Decode(&newMessage)

	// Verifying the correctness of each and every field
	err_constraints := checkMessageCorrectness(newMessage, rt)
	if err != nil || !(err_constraints) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		forbiddenError := Response{
			Code:    400,
			Message: "The message has one or more non-valid field(s)",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Uploading the message

	// Sending the response
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
}

// A function that checks the correctness of every field in a sent message
func checkMessageCorrectness(newMessage Message, rt *_router) bool {
	id := newMessage.MessageID >= 0 && newMessage.MessageID <= 10000
	timestamp1, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, newMessage.Timestamp)
	timestamp2 := len(newMessage.Timestamp) == 20
	no_content := !(len(newMessage.Content) == 0 && len(newMessage.Photo) == 0)
	var message_validity bool = false
	if len(newMessage.Content) > 0 {
		message_validity, _ = regexp.MatchString(`.{1,300}`, newMessage.Content)
	} else if len(newMessage.Photo) > 0 {
		message_validity, _ = regexp.MatchString(`[-A-Za-z0-9+/=]|=[^=]|={3,16}`, newMessage.Photo)
	}
	user := rt.db.CheckUsername(newMessage.Username)
	checkmarks := newMessage.Checkmarks == 0
	comments := len(newMessage.Comments) == 0
	sentbyme := newMessage.SentByMe
	replyingto := rt.db.CheckMessage(string(newMessage.MessageID))

	correctness := id && timestamp1 && timestamp2 && no_content && message_validity && user && checkmarks && comments && sentbyme && replyingto

	return correctness
}
