package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It modifies the profile picture of an existing user.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Message sending"

	// Checking that the conversation actually exists
	convID, err := strconv.Atoi(ps.ByName("conversationid"))
	if err != nil {
		_ = createBackendError(affinity, "Conversation retrieval has failed", err, w)
		return
	}

	convs, err := ConversationFromIdRetrieval(convID, rt, w)
	if err != nil {
		return
	}

	if len(convs) == 0 {

	}

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new message
	var newMessage RequestMessage
	err = json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a message", affinity, "Request encoding for badly formatted message response has failed", w)
		return
	}

	// Checking if the message is valid
	match, err := checkMessageCorrectness(newMessage, rt, w)
	if err != nil {
		return
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		createFaultyResponse(http.StatusBadRequest, "Message parsed incorrectly or not valid", affinity, "Request encoding for message not correcly formatted response has failed", w)
		return
	}

	// Creating the id
	id, err := MessageIdCreator(rt, w)
	if err != nil {
		return
	}

	// Getting the timestamp
	timestamp := GetTime()

	// Getting the username
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	// Getting the message id of who we are replying to
	replyingTo := fmt.Sprintf("'%s'", strconv.Itoa(newMessage.ReplyingTo))
	if replyingTo == "'-1'" {
		replyingTo = "NULL"
	}

	// Actually writing the message in the DB
	query := fmt.Sprintf("(%d, '%s', '%s', '%s', '%s', %d, %s)", id, user[1], timestamp, newMessage.Content, newMessage.Photo, 0, replyingTo)

	_, err = rt.db.Insert("messages", query)
	if err != nil {
		_ = createBackendError(affinity, "Inserting the new message into the database has failed", err, w)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}

// A function that checks the correctness of every field in a sent message
func checkMessageCorrectness(newMessage RequestMessage, rt *_router, w http.ResponseWriter) (bool, error) {
	var message_validity bool = false
	var replying_to bool
	var err error

	// Logging information
	const affinity string = "Message sending"

	// Checking that text is valid
	if len(newMessage.Content) > 0 {
		message_validity, err = regexp.MatchString(`.{1,300}`, newMessage.Content)
		if err != nil {
			return false, createBackendError(affinity, "Matching message content with appropriate regex failed", err, w)
		}

		if !message_validity {
			return false, nil
		}
	}

	// Checking that photo is valid
	if len(newMessage.Photo) > 0 {
		message_validity, err = regexp.MatchString(`[-A-Za-z0-9+/=]|=[^=]|={3,16}`, newMessage.Photo)
		if err != nil {
			return false, createBackendError(affinity, "Matching message photo with appropriate regex failed", err, w)
		}

		if !message_validity {
			return false, nil
		}
	}

	// Replying to a message with id -1 corresponds to
	// responding to no one, since messages' id go up from 0
	if newMessage.ReplyingTo == -1 {
		replying_to = true
	} else {
		replying_to, err = MessageFromIdExists(newMessage.ReplyingTo, rt, w)
		if err != nil {
			return false, createBackendError(affinity, "Checking that the message we are replying to's id failed", err, w)
		}
	}

	correctness := replying_to && message_validity

	return correctness, nil
}

// It creates a numerical ID for the new message.
func MessageIdCreator(rt *_router, w http.ResponseWriter) (int, error) {
	var id int

	// Logging information
	const affinity string = "Message sending"

	for {
		id = rand.Intn(10001)
		rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = '%d'", id))
		if err != nil {
			return 0, createBackendError(affinity, "SELECT in the database seeking messages with the same id failed", err, w)
		}

		// Checking that the new id is unique
		other_messages, err := MessageRowReading(rows)

		if err != nil {
			return 0, createBackendError(affinity, "Reading the database rows that were seeking messages with the same id failed", err, w)
		} else if len(other_messages) == 0 {
			break
		}
	}

	return id, nil
}
