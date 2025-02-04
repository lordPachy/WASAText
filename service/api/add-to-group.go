package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It adds a member to group.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Conversation creation"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new adding request
	var newMembership AddToGroupRequest
	err = json.NewDecoder(r.Body).Decode(newMembership)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a group membership request", affinity, "Request encoding for badly formatted group membership request response has failed", w)
		return
	}

	// Checking if the conversation request is valid
	match, err := checkConversationRequestCorrectness(newConvo, rt, w)
	if err != nil {
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "Conversation request parsed incorrectly or not valid", affinity, "Request encoding for conversation encoding not correcly formatted response has failed", w)
		return
	}

	// Checking that the conversation actually exists
	convID, err := strconv.Atoi(ps.ByName("conversationid"))
	if err != nil {
		_ = createBackendError(affinity, "Conversation retrieval has failed", err, w)
		return
	}

	exists, err := ConversationFromIdExistence(convID, rt, w)
	if err != nil {
		return
	}
	if !exists {
		createFaultyResponse(http.StatusNotFound, "Conversation not found", affinity, "Response message encoding for conversation not found error has failed", w)
		return
	}

	// Checking if the conversation is a group
	if convID < 5000 {
		createFaultyResponse()
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
				_ = createBackendError(affinity, "Converting private chat id for existing private chat failed", err, w)
				return
			}

			newToken := ConversationID{
				Id: id,
			}
			err = json.NewEncoder(w).Encode(newToken)
			if err != nil {
				_ = createBackendError(affinity, "Encoding the private chat id for an existing private chat has failed", err, w)
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
			_ = createBackendError(affinity, "Inserting the new conversation into the database has failed", err, w)
			return
		}
	} else {
		// Writing in the database of groups
		query := fmt.Sprintf("(%d, '%s', NULL)", id, newConvo.GroupName)
		_, err = rt.db.Insert("groupchats", query)
		if err != nil {
			_ = createBackendError(affinity, "Inserting the new group into the database has failed", err, w)
			return
		}

		// Writing in the database of group members
		for _, usr := range newConvo.Members {
			query := fmt.Sprintf("(%d, '%s')", id, usr.Name)
			_, err = rt.db.Insert("groupmembers", query)
			if err != nil {
				_ = createBackendError(affinity, "Inserting the new groupmember into the database has failed", err, w)
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
		_ = createBackendError(affinity, "Encoding the chat id for a new chat has failed", err, w)
		return
	}
}

// A function that checks the correctness of every field in a new group membership request
func checkMembershipRequestCorrectness(newMembership AddToGroupRequest, adder Access_token, rt *_router, w http.ResponseWriter) (bool, error) {
	var group_existence bool = false
	var adder_belonging bool = true
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
			return false, createBackendError(affinity, "The string matching mechanism for conversation request correctness has failed", err, w)
		}
	}
	correctness := user_number && user_existence && group_name
	return correctness, nil
}
