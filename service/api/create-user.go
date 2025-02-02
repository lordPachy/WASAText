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

	// Getting the new username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    400,
			Message: "The received body is not a username",
		}
		err = json.NewEncoder(w).Encode(badRequest)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "User creation",
				Message:  "Request encoding for bad username has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
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
		err = json.NewEncoder(w).Encode(badUsername)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "User creation",
				Message:  "Request encoding for not regex-matching username has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
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
		err = json.NewEncoder(w).Encode(forbiddenError)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "User creation",
				Message:  "Request encoding for username already in use has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return
		}
		return
	}

	// Accepted request
	// Creating the id
	id, err := IdCreator(rt)
	if err != nil {
		fmt.Println(err.Error())
		return
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

func UsernameRetrieval(username Username, rt *_router) ([]string, error) {
	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("username = '%s'", username.Name))
	if err != nil {
		selectionError := BackendError{
			Affinity: "User retrieval",
			Message:  "SELECT in the database seeking users with the same username failed",
			OG_error: err,
		}
		return nil, &selectionError
	}

	// Reading the rows
	other_users, err := UsersRowReading(rows)

	if err != nil {
		uniquenessError := BackendError{
			Affinity: "User retrieval",
			Message:  "Reading the database rows that were seeking users with the same username failed",
			OG_error: err,
		}
		fmt.Println(uniquenessError.Error())
		return nil, &uniquenessError
	}

	return other_users, nil
}

func IdCreator(rt *_router) (string, error) {
	var id string

	for {
		id, _ = reggen.Generate("^[0-9A-Za-z-$&/(),.]{4,16}$", 16)
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
