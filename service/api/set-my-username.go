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

// Set a username and return the new username
func (rt *_router) setMyUserName2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var newUsername Username

	// Authentication
	id, err := rt.authorization(w, r)
	fmt.Println("Authorization: successful. Id is:")
	fmt.Println(id)
	if err != nil {
		return
	}

	// Getting the new username and validity check
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	err_constraints, _ := regexp.MatchString(`\w{3,16}`, newUsername.Name)

	if err != nil || !(err_constraints) {
		w.WriteHeader(http.StatusBadRequest)
		forbiddenError := Response{
			Code:    400,
			Message: "The username is not valid",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Setting the new username and uniqueness check
	err = rt.db.ChangeUsername(id, newUsername.Name)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		forbiddenError := Response{
			Code:    403,
			Message: "The username tried out is already in use",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusNoContent)
}
