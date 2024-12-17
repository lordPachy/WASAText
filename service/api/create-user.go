package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It sets a user's profile picture. It returns error if the picture is wrongly formatted.
func (rt *_router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var newUsername Username

	// Getting the new username and validity check
	err := json.NewDecoder(r.Body).Decode(&newUsername)
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
	err = rt.db.CreateUser(newUsername)

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
	w.WriteHeader(http.StatusCreated)
}
