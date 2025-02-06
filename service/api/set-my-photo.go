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

	// Logging information
	const affinity string = "Profile photo setting"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the new profile picture photo
	var newPhoto Image
	err = json.NewDecoder(r.Body).Decode(&newPhoto)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not an image", affinity, "Request encoding for badly formatted image has failed", w, rt)
		return
	}

	// Checking if the image is valid
	match, err := regexp.MatchString(`^0b[01]{1,14}$`, newPhoto.Image)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for picture string has failed", err, w, rt)
		return
	}
	if !match {
		createFaultyResponse(http.StatusBadRequest, "The profile picture is not valid (it may be too short, or long, or containing not valid characters)", affinity, "Request encoding for not valid profile picture response failed", w, rt)
		return
	}

	// Actually setting the picture in the DB
	_, err = rt.db.Update("users", fmt.Sprintf("propic = '%s'", newPhoto.Image), fmt.Sprintf("id = '%s'", token.Identifier))
	if err != nil {
		_ = createBackendError(affinity, "Updating user with the new profile picture has failed", err, w, rt)
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "text-plain")
	w.WriteHeader(http.StatusAccepted)
}
