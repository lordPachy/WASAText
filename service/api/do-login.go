package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It logins an existing user.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var username Username

	// Logging information
	affinity := "Login"

	// Getting the username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a username", affinity, "Request encoding for username not correctly formatted response has failed", w, rt)
		return
	}

	// Id retrieval
	users, err := UserFromUsernameRetrieval(username, rt, w)
	if err != nil {
		return
	}

	if len(users) == 0 {
		createFaultyResponse(http.StatusNotFound, "No user corresponds to the given username", affinity, "Request encoding for user not found in database response has failed", w, rt)
		return
	}

	id := users[0]

	// Writing the response in HTTP
	w.WriteHeader(http.StatusCreated)
	accessToken := Access_token{
		Identifier: id,
	}

	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		_ = createBackendError(affinity, "Encoding the new access token has failed", err, w, rt)
		return
	}
}
