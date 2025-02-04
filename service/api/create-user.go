package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/lucasjones/reggen"

	"github.com/julienschmidt/httprouter"
)

// It creates a user.
func (rt *_router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var newUsername Username

	// Logging information
	affinity := "User creation"

	// Getting the new username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a username", affinity, "Encoding of bad username error failed", w)
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newUsername.Name)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for id creation has failed", err, w)
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "The username is not valid (it may be too short, or long, or containing not valid characters)", affinity, "Request encoding for not regex-matching username has failed", w)
		return
	}

	// Uniqueness check
	other_users, err := UserFromUsernameRetrieval(newUsername, rt, w)
	if err != nil {
		return
	}

	if len(other_users) > 0 {
		createFaultyResponse(http.StatusForbidden, "The username tried out is already in use", affinity, "Request encoding for username already in use has failed", w)
		return
	}

	// Accepted request
	// Creating the id
	id, err := userIdCreator(rt, w)
	if err != nil {
		return
	}

	// Inserting the user into the database
	_, err = rt.db.Insert("users", fmt.Sprintf("('%s', '%s', Null)", id, newUsername.Name))
	if err != nil {
		_ = createBackendError(affinity, "Inserting the new user into the database has failed", err, w)
		return
	}

	// Writing the response in HTTP
	w.WriteHeader(http.StatusCreated)
	newToken := Access_token{
		Identifier: id,
	}
	err = json.NewEncoder(w).Encode(newToken)
	if err != nil {
		_ = createBackendError(affinity, "Encoding the new access token has failed", err, w)
		return
	}
}

// It creates a userID and returns it a string; otherwise, automatically handles the error
func userIdCreator(rt *_router, w http.ResponseWriter) (string, error) {
	var id string

	// Logging information
	affinity := "User creation"

	for {
		id, _ = reggen.Generate("^[A-Za-z]{16}$", 16)
		rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id))
		if err != nil {
			return "", createBackendError(affinity, "SELECT in the database seeking users with the same id failed", err, w)
		}

		// Checking that the new id is unique
		other_users, err := UsersRowReading(rows)

		if err != nil {
			return "", createBackendError(affinity, "Reading the database rows that were seeking users with the same id failed", err, w)
		} else if len(other_users) == 0 {
			break
		}
	}

	return id, nil
}
