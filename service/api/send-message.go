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

	// Reading the message
	var newMessage Message
	err = json.NewDecoder(r.Body).Decode(&newMessage)

	// Verifying the correctness of each and every field

}

func checkMessageCorrectness(newMessage Message) {
	id_check := newMessage.MessageID >= 0 && newMessage.MessageID <= 10000
	timestamp_check1, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, newMessage.Timestamp)
	timestamp_check2 := len(newMessage.Timestamp) == 20
	no_content = !(len(newMessage.Content) == 0 && len(newMessage.Photo) == 0)

}
