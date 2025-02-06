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

// It creates a private conversation or a group, depending on a flag.
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Conversation creation"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new conversation
	var newConvo ConversationRequest
	err = json.NewDecoder(r.Body).Decode(&newConvo)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a conversation request", affinity, "Request encoding for badly formatted conversation request response has failed", w, rt)
		return
	}

	// Checking if the conversation request is valid
	match, err := checkConversationRequestCorrectness(newConvo, rt, w)
	if err != nil {
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "Conversation request parsed incorrectly or not valid", affinity, "Request encoding for conversation encoding not correcly formatted response has failed", w, rt)
		return
	}

	// Adding sender's own id to the request
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	sender := Username{
		Name: user[1],
	}
	newConvo.Members = append(newConvo.Members, sender)

	// Private conversations are unique: checking if the private convo already exists
	if !newConvo.IsGroup {
		chat, err := PrivConversationFromMembersRetrieval(newConvo.Members[0], newConvo.Members[1], rt, w)
		if err != nil {
			return
		}

		// If it exists, we enforce idempotency by just returning the already existing private chat
		if len(chat) > 0 {
			w.WriteHeader(http.StatusOK)
			id, err := strconv.Atoi(chat[0])

			if err != nil {
				_ = createBackendError(affinity, "Converting private chat id for existing private chat failed", err, w, rt)
				return
			}

			newToken := ConversationID{
				Id: id,
			}
			err = json.NewEncoder(w).Encode(newToken)
			if err != nil {
				_ = createBackendError(affinity, "Encoding the private chat id for an existing private chat has failed", err, w, rt)
				return
			}

			return
		}
	}

	// Creating the id
	id, err := ConversationIdCreator(rt, newConvo.IsGroup, w)
	if err != nil {
		return
	}

	// Actually writing the conversation in the DB
	if !newConvo.IsGroup {
		query := fmt.Sprintf("(%d, '%s', '%s')", id, newConvo.Members[0].Name, newConvo.Members[1].Name)

		_, err = rt.db.Insert("privchats", query)
		if err != nil {
			_ = createBackendError(affinity, "Inserting the new conversation into the database has failed", err, w, rt)
			return
		}
	} else {
		// Writing in the database of groups
		query := fmt.Sprintf("(%d, '%s', NULL)", id, newConvo.GroupName)
		_, err = rt.db.Insert("groupchats", query)
		if err != nil {
			_ = createBackendError(affinity, "Inserting the new group into the database has failed", err, w, rt)
			return
		}

		// Writing in the database of group members
		for _, usr := range newConvo.Members {
			query := fmt.Sprintf("(%d, '%s')", id, usr.Name)
			_, err = rt.db.Insert("groupmembers", query)
			if err != nil {
				_ = createBackendError(affinity, "Inserting the new groupmember into the database has failed", err, w, rt)
				return
			}
		}
	}
	// Writing the response in HTTP
	// Accepted request
	w.WriteHeader(http.StatusCreated)
	newToken := ConversationID{
		Id: id,
	}
	err = json.NewEncoder(w).Encode(newToken)
	if err != nil {
		_ = createBackendError(affinity, "Encoding the chat id for a new chat has failed", err, w, rt)
		return
	}
}

// A function that checks the correctness of every field in a conversation request
func checkConversationRequestCorrectness(newConvo ConversationRequest, rt *_router, w http.ResponseWriter) (bool, error) {
	var user_number bool
	var user_existence bool = false
	var group_name bool = true
	var err error

	// Logging information
	const affinity string = "Conversation creation"

	// Enforcing number of members contraints
	if !newConvo.IsGroup {
		user_number = len(newConvo.Members) == 1
	} else {
		user_number = ((len(newConvo.Members) >= 1) && (len(newConvo.Members) <= 999))
	}

	// Ensuring members actually exist
	for _, usr := range newConvo.Members {
		user_existence, err = UserFromUsernameExists(usr, rt, w)
		if err != nil {
			return false, err
		}
		if !user_existence {
			break
		}
	}

	// Ensuring the groupname is valid (if it is a group)
	if newConvo.IsGroup {
		// Checking if the username is valid
		group_name, err = regexp.MatchString(`^[\w\ ]{3,16}$`, newConvo.GroupName)
		if err != nil {
			return false, createBackendError(affinity, "The string matching mechanism for conversation request correctness has failed", err, w, rt)
		}
	}
	correctness := user_number && user_existence && group_name
	return correctness, nil
}

// It creates a numerical ID for the new conversation.
func ConversationIdCreator(rt *_router, isgroup bool, w http.ResponseWriter) (int, error) {
	var id int

	for {
		id = rand.Intn(5000)

		// Private conversations have ids in [0, 5000); groups have ids in [5000, 10000)
		if isgroup {
			id += 5000
		}

		// Checking that the new id is unique
		other_convos, err := ConversationFromIdRetrieval(id, rt, w)

		if err != nil {
			return -1, err
		} else if len(other_convos) == 0 {
			break
		}
	}

	return id, nil
}
