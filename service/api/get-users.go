package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It retrieves a list of users, possibly filtered (as in a search).
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Getting list of users"

	// Authentication
	_, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the user name
	var userQuery Username
	err = json.NewDecoder(r.Body).Decode(&userQuery)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a user query", affinity, "Request encoding for badly formatted user query has failed", w, rt)
		return
	}

	// Retrieving users
	var users []Username
	usersraw, err := UserQuerying(userQuery, rt, w)
	if err != nil {
		return
	}

	i := 0
	for i < len(usersraw) {
		users = append(users, Username{usersraw[i+1]})
		i += 3
	}

	// Writing the response in HTTP
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		_ = createBackendError(affinity, "Encoding users has failed", err, w, rt)
		return
	}
}
