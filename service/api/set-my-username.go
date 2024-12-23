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

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		fmt.Println(err.Error())
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
		json.NewEncoder(w).Encode(badRequest)
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newUsername.Name)

	if err != nil {
		regexError := BackendError{
			Affinity: "User creation",
			Message:  "The string matching mechanism for id creation has failed",
			OG_error: err,
		}
		fmt.Println(regexError.Error())
		return
	}
	if !match {
		w.WriteHeader(http.StatusBadRequest)
		badUsername := Response{
			Code:    400,
			Message: "The username is not valid (it may be too short, or long, or containing not valid characters)",
		}
		json.NewEncoder(w).Encode(badUsername)
		return
	}

	// Uniqueness check
	other_users, err := UsernameRetrieval(newUsername, rt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(other_users) > 0 {
		w.WriteHeader(http.StatusForbidden)
		forbiddenError := Response{
			Code:    403,
			Message: "The username tried out is already in use",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Actually setting the username in the DB
	_, err = rt.db.Update("users", fmt.Sprintf("username = '%s'", newUsername.Name), fmt.Sprintf("id = '%s'", token.Identifier))
	if err != nil {
		usernameUpdateError
	}

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
