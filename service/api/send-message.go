package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/julienschmidt/httprouter"
)

// It modifies the profile picture of an existing user.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Setting logging information
	affinity := "Message sending"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Getting the new message
	var newMessage RequestMessage
	err = json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    400,
			Message: "The received body is not a message",
		}
		err = json.NewEncoder(w).Encode(badRequest)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: affinity,
				Message:  "Request encoding for badly formatted message response has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
		return
	}

	// Checking if the message is valid
	match, err := checkMessageCorrectness(newMessage, rt)
	if err != nil {
		formatError := BackendError{
			Affinity: affinity,
			Message:  "Message correctness checking has failed",
			OG_error: err,
		}
		fmt.Println(formatError.Error())
		return
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		badPic := Response{
			Code:    400,
			Message: "Message parsed incorrectly or not valid",
		}
		err = json.NewEncoder(w).Encode(badPic)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: affinity,
				Message:  "Request encoding for message not correcly formatted response has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
		return
	}

	// Creating the id
	id, err := MessageIdCreator(rt)
	if err != nil {
		idError := BackendError{
			Affinity: affinity,
			Message:  "Creating the message ID has failed",
			OG_error: err,
		}
		fmt.Println(idError.Error())
		return
	}

	// Getting the timestamp
	timestamp := GetTime()

	// Getting the username
	user, err := UsernameFromID(token.Identifier, rt)
	if err != nil {
		usernameError := BackendError{
			Affinity: affinity,
			Message:  "Retrieving the username has failed",
			OG_error: err,
		}
		fmt.Println(usernameError.Error())
		return
	}

	// Actually writing the message in the DB
	_, err = rt.db.Insert("messages", fmt.Sprintf("('%d', '%s', '%s', '%s', '%s', '%d', '%d')", id, user.Name, timestamp, newMessage.Content, newMessage.Photo, 0, newMessage.ReplyingTo))
	if err != nil {
		insertionError := BackendError{
			Affinity: "Message sending",
			Message:  "Inserting the new message into the database has failed",
			OG_error: err,
		}
		fmt.Println(insertionError.Error())
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}

// package api

/*import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It sends a message to a specified chat.
// Implementation notes:
// 1. Creating the message
// 2. Adding to the message database
// 3. Adding to the chat database
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Authentication
	id, err := rt.authorization(w, r)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
		return
	}

	// Checking whether the conversation exists
	convo := ps.ByName("conversationid")
	convo_check := rt.db.CheckChat(convo, id)

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
*/

// A function that checks the correctness of every field in a sent message
func checkMessageCorrectness(newMessage RequestMessage, rt *_router) (bool, error) {
	var message_validity bool = false
	var replying_to bool
	var err error

	// Checking that text is valid
	if len(newMessage.Content) > 0 {
		message_validity, err = regexp.MatchString(`.{1,300}`, newMessage.Content)
		if err != nil {
			return false, err
		}

		if !message_validity {
			return false, nil
		}
	}

	// Checking that photo is valid
	if len(newMessage.Photo) > 0 {
		message_validity, err = regexp.MatchString(`[-A-Za-z0-9+/=]|=[^=]|={3,16}`, newMessage.Photo)
		if err != nil {
			return false, err
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
		replying_to, err = MessageExists(newMessage.ReplyingTo, rt)
		if err != nil {
			replyingError := BackendError{
				Affinity: "Message sending",
				Message:  "Checking that the message we are replying to's id failed",
				OG_error: err,
			}
			return false, &replyingError
		}
	}

	correctness := replying_to && message_validity

	return correctness, nil
}

func MessageIdCreator(rt *_router) (int, error) {
	var id int

	for {
		id = rand.Intn(10001)
		rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = '%d'", id))
		if err != nil {
			selectionError := BackendError{
				Affinity: "Message sending",
				Message:  "SELECT in the database seeking messages with the same id failed",
				OG_error: err,
			}
			fmt.Println(selectionError.Error())
			return 0, &selectionError
		}

		// Checking that the new id is unique
		other_messages, err := MessageRowReading(rows)

		if err != nil {
			idUniquenessError := BackendError{
				Affinity: "Message sending",
				Message:  "Reading the database rows that were seeking messages with the same id failed",
				OG_error: err,
			}
			return 0, &idUniquenessError
		} else if len(other_messages) == 0 {
			break
		}
	}

	return id, nil
}

func MessageRowReading(res *sql.Rows) ([]string, error) {
	// Retrieving the values from rows
	var answer []string // array of actual values
	var id, sender, created_at, content, photo, checkmarks, replying_to *string
	for {
		if res.Next() { // there is another value to be scanned
			err := res.Scan(&id, &sender, &created_at, &content, &photo, &checkmarks, &replying_to)
			if err == nil {
				tmp := "Null"
				if content == nil {
					content = &tmp
				}
				if photo == nil {
					photo = &tmp
				}
				if replying_to == nil {
					replying_to = &tmp
				}
				answer = append(answer, *id, *sender, *created_at, *content, *photo, *checkmarks, *replying_to)
			} else {
				return nil, err // the scan has had an error
			}
		} else {
			if res.Err() == nil { // there are no more values to scan in the current set
				if res.NextResultSet() { // there are values to be scanned
					continue // in the next set
				} else {
					if res.Err() == nil { // there are no more values, and the scan can end
						break
					} else { // next set scan went unsuccessfully
						return nil, res.Err()
					}
				}
			} else { // next scan went unsuccessfully
				return nil, res.Err()
			}
		}
	}

	return answer, nil
}

func MessageExists(id int, rt *_router) (bool, error) {
	// Querying database rows
	rows, err := rt.db.Select("*", "messages", fmt.Sprintf("id = '%d'", id))
	if err != nil {
		selectionError := BackendError{
			Affinity: "Message sending",
			Message:  "SELECT in the database seeking messages with the same id failed",
			OG_error: err,
		}
		fmt.Println(selectionError.Error())
		return false, &selectionError
	}

	// Checking the queried rows
	other_messages, err := MessageRowReading(rows)

	if err != nil {
		idError := BackendError{
			Affinity: "Message sending",
			Message:  "Reading the database rows that were seeking messages with the same id failed",
			OG_error: err,
		}
		return false, &idError
	} else if len(other_messages) == 0 {
		return false, nil
	}

	if len(other_messages) > 1 {
		idUniquenessError := BackendError{
			Affinity: "Message sending",
			Message:  "The database contains more than 1 message with the same id",
			OG_error: err,
		}
		return false, &idUniquenessError
	}

	return true, nil
}

func GetTime() string {
	currentTime := time.Now()
	return currentTime.Format("2017-07-21T17:32:28Z")
}

func UsernameFromID(id string, rt *_router) (*Username, error) {
	// Querying database rows
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id))
	if err != nil {
		return nil, err
	}

	// Checking the queried rows
	usernames, err := UsersRowReading(rows)

	if err != nil {
		return nil, err
	} else if len(usernames) == 0 {
		return nil, nil
	}

	username := Username{
		Name: usernames[0],
	}

	return &username, nil
}
