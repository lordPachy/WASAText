package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It modifies the username of an existing user.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Username setting"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new username
	var newUsername Username
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a username", affinity, "Request encoding for badly formatted username has failed", w, rt)
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newUsername.Name)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for id creation has failed", err, w, rt)
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "The username is not valid (it may be too short, or long, or containing not valid characters)", affinity, "Request encoding for username not matching with regex response has failed", w, rt)
		return
	}

	// Uniqueness check
	other_users, err := UserFromUsernameRetrieval(newUsername, rt, w)
	if err != nil {
		return
	}

	if len(other_users) > 0 {
		createFaultyResponse(http.StatusForbidden, "The username tried out is already in use", affinity, "Request encoding for username already in user request has failed", w, rt)
		return
	}

	// Actually setting the username in the DB
	_, err = rt.db.Update("users", fmt.Sprintf("username = '%s'", newUsername.Name), fmt.Sprintf("id = '%s'", token.Identifier))
	if err != nil {
		_ = createBackendError(affinity, "Updating user with the new id has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
