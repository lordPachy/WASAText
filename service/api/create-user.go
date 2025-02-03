package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/lucasjones/reggen"

	"github.com/julienschmidt/httprouter"
)

// It creates a user.
func (rt *_router) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var newUsername Username

	// Logging information
	affinity := "User creation"

	// Getting the new username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a username", affinity, "Encoding of bad username error failed", w)
		return
	}

	// Checking if the username is valid
	match, err := regexp.MatchString(`^\w{3,16}$`, newUsername.Name)

	if err != nil {
		_ = createBackendError(affinity, "The string matching mechanism for id creation has failed", err, w)
		return
	}

	if !match {
		createFaultyResponse(http.StatusBadRequest, "The username is not valid (it may be too short, or long, or containing not valid characters)", affinity, "Request encoding for not regex-matching username has failed", w)
		return
	}

	// Uniqueness check
	other_users, err := UsernameRetrieval(newUsername, rt)
	if err != nil {
		return
	}

	if len(other_users) > 0 {
		createFaultyResponse(http.StatusForbidden, "The username tried out is already in use", affinity, "Request encoding for username already in use has failed", w)
		return
	}

	// Accepted request
	// Creating the id
	id, err := userIdCreator(rt)
	if err != nil {
		idError := BackendError{
			Affinity: "User creation",
			Message:  "Creating the user ID has failed",
			OG_error: err,
		}
		fmt.Println(idError.Error())
	}

	// Inserting the user into the database
	_, err = rt.db.Insert("users", fmt.Sprintf("('%s', '%s', Null)", id, newUsername.Name))
	if err != nil {
		insertionError := BackendError{
			Affinity: "User creation",
			Message:  "Inserting the new user into the database has failed",
			OG_error: err,
		}
		fmt.Println(insertionError.Error())
	}

	// Writing the response in HTTP
	w.WriteHeader(http.StatusCreated)
	newToken := Access_token{
		Identifier: id,
	}
	err = json.NewEncoder(w).Encode(newToken)
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

func userIdCreator(rt *_router) (string, error) {
	var id string

	for {
		id, _ = reggen.Generate("^[A-Za-z]{16}$", 16)
		rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id))
		if err != nil {
			selectionError := BackendError{
				Affinity: "User creation",
				Message:  "SELECT in the database seeking users with the same id failed",
				OG_error: err,
			}
			fmt.Println(selectionError.Error())
			return "", &selectionError
		}

		// Checking that the new id is unique
		other_users, err := UsersRowReading(rows)

		if err != nil {
			idUniquenessError := BackendError{
				Affinity: "User creation",
				Message:  "Reading the database rows that were seeking users with the same id failed",
				OG_error: err,
			}
			return "", &idUniquenessError
		} else if len(other_users) == 0 {
			break
		}
	}

	return id, nil
}
