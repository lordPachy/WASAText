package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It adds a member to group.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "New group membership"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new adding request
	var newMembership AddToGroupRequest
	err = json.NewDecoder(r.Body).Decode(&newMembership)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a group membership request", affinity, "Request encoding for badly formatted group membership request response has failed", w, rt)
		return
	}

	// Checking if the conversation request is valid
	match := checkMembershipRequestCorrectness(newMembership, token, rt, w)
	if !match {
		return
	}

	// Actually adding the membership in the DB
	query := fmt.Sprintf("(%d, '%s')", newMembership.Group.Id, newMembership.User.Name)

	_, err = rt.db.Insert("groupmembers", query)
	if err != nil {
		_ = createBackendError(affinity, "Inserting the new group membership into the database has failed", err, w, rt)
		return
	}
	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusOK)
}

// A function that checks the correctness of every field in a new group membership request.
func checkMembershipRequestCorrectness(newMembership AddToGroupRequest, adder Access_token, rt *_router, w http.ResponseWriter) bool {
	var err error

	// Logging information
	const affinity string = "New group membership"

	// Checking that the member to be added exists
	user_existence, err := UserFromUsernameExists(newMembership.User, rt, w)
	if err != nil {
		return false
	}

	if !user_existence {
		createFaultyResponse(http.StatusNotFound, "The user to be added does not exist", affinity, "Response encoding for user not found error has failed", w, rt)
		return false
	}

	// Checking that the group to be added to exists
	group_existence, err := ConversationFromIdExistence(newMembership.Group.Id, rt, w)
	if err != nil {
		return false
	}

	if !group_existence {
		createFaultyResponse(http.StatusNotFound, "The group to be added to does not exist", affinity, "Response encoding for conversation not found error has failed", w, rt)
		return false
	}

	// Checking that the adder belongs to the group they are adding to
	adder_belonging, err := UserBelongsToGroup(adder, newMembership.Group, rt, w)
	if err != nil {
		return false
	}

	if !adder_belonging {
		createFaultyResponse(http.StatusNotFound, "The user adding does not belong to the group", affinity, "Response encoding for user not having adding privileges has failed", w, rt)
		return false
	}

	return true
}
