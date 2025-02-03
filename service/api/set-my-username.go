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
	affinity := "Username setting"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new username
	var newUsername Username
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    400,
			Message: "The received body is not a username",
		}
		err = json.NewEncoder(w).Encode(badRequest)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			_ = createBackendError(affinity, "Request encoding for badly formatted username has failed", err, w)
			return
		}
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newUsername.Name)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for id creation has failed", err, w)
		return
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		badUsername := Response{
			Code:    400,
			Message: "The username is not valid (it may be too short, or long, or containing not valid characters)",
		}
		err = json.NewEncoder(w).Encode(badUsername)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			_ = createBackendError(affinity, "Request encoding for username not matching with regex response has failed", err, w)
			return
		}
		return
	}

	// Uniqueness check
	other_users, err := UsernameRetrieval(newUsername, rt)
	if err != nil {
		_ = createBackendError(affinity, "New username retrieving for uniqueness check has failed", err, w)
		return
	}

	if len(other_users) > 0 {
		w.WriteHeader(http.StatusForbidden)
		forbiddenError := Response{
			Code:    403,
			Message: "The username tried out is already in use",
		}
		err = json.NewEncoder(w).Encode(forbiddenError)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			_ = createBackendError(affinity, "Request encoding for username already in user request has failed", err, w)
			return
		}

		return
	}

	// Actually setting the username in the DB
	_, err = rt.db.Update("users", fmt.Sprintf("username = '%s'", newUsername.Name), fmt.Sprintf("id = '%s'", token.Identifier))
	if err != nil {
		_ = createBackendError(affinity, "Updating user with the new id has failed", err, w)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
