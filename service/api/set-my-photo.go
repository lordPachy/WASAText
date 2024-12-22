package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It sets a user's profile picture. It returns error if the picture is wrongly formatted.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var newPhoto string

	// Authentication
	id, err := rt.authorization(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UnauthorizedError)
		return
	}

	// Getting the new photo and validity check
	err = json.NewDecoder(r.Body).Decode(&newPhoto)
	err_constraints, _ := regexp.MatchString(`0b[01]{1,14}`, newPhoto)

	if err != nil || !(err_constraints) {
		w.WriteHeader(http.StatusBadRequest)
		forbiddenError := Response{
			Code:    400,
			Message: "The profile picture is not valid",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Loading the photo on the database
	err = rt.db.ChangeProPic(id, newPhoto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		forbiddenError := Response{
			Code:    400,
			Message: "The profile picture is not valid",
		}
		json.NewEncoder(w).Encode(forbiddenError)
		return
	}

	// Accepted request
	var AcceptedResponse = Response{
		Code:    202,
		Message: "Profile picture updated successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AcceptedResponse)
}
