package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It logins an existing user.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var username Username

	// Getting the username
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    400,
			Message: "The received body is not a username",
		}
		json.NewEncoder(w).Encode(badRequest)
		return
	}

	// Id retrieval
	users, err := UsernameRetrieval(username, rt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		notExisting := Response{
			Code:    404,
			Message: "No user corresponds to the given username",
		}
		json.NewEncoder(w).Encode(notExisting)
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
		encodingError := BackendError{
			Affinity: "User creation",
			Message:  "Encoding the new access token has failed",
			OG_error: err,
		}
		fmt.Println(encodingError.Error())
		return
	}
}
