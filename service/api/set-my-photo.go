package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// It modifies the profile picture of an existing user.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Getting the new profile picture photo
	var newPhoto Image
	err = json.NewDecoder(r.Body).Decode(&newPhoto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    400,
			Message: "The received body is not an image",
		}
		err = json.NewEncoder(w).Encode(badRequest)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "Profile picture setting",
				Message:  "Request encoding for badly formatted image has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
		return
	}

	// Checking if the image is valid
	match, err := regexp.MatchString(`^0b[01]{1,14}$`, newPhoto.Image)

	if err != nil {
		regexError := BackendError{
			Affinity: "Profile picture setting",
			Message:  "The string matching mechanism for picture string has failed",
			OG_error: err,
		}
		fmt.Println(regexError.Error())
		return
	}
	if !match {
		w.WriteHeader(http.StatusBadRequest)
		badPic := Response{
			Code:    400,
			Message: "The profile picture is not valid (it may be too short, or long, or containing not valid characters)",
		}
		err = json.NewEncoder(w).Encode(badPic)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "Profile picture setting",
				Message:  "Request encoding for username not matching with regex response has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
		return
	}

	// Actually setting the picture in the DB
	_, err = rt.db.Update("users", fmt.Sprintf("propic = '%s'", newPhoto.Image), fmt.Sprintf("id = '%s'", token.Identifier))
	if err != nil {
		propicUpdateError := BackendError{
			Affinity: "Profile picture setting",
			Message:  "Updating user with the new profile picture has failed",
			OG_error: err,
		}
		fmt.Println(propicUpdateError.Error())
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
