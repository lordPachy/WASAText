package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var user Username
	var id Access_token

	// Getting the username and validity check
	err := json.NewDecoder(r.Body).Decode(&user)
	err_constraints, _ := regexp.MatchString(`\w{3,16}`, user.Name)

	if err != nil || !(err_constraints) {
		w.WriteHeader(http.StatusBadRequest)
		forbiddenError := Response{
			Code:    400,
			Message: "The username is not valid",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Identifier retrieval
	id.Identifier, err = rt.db.GetIdentifier(user.Name)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(id)
	}
}
