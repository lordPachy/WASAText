package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Authentication(w http.ResponseWriter, r *http.Request, rt *_router) (Access_token, error) {
	id := r.Header.Get("Authentication")
	token := Access_token{
		Identifier: id,
	}

	user, err := IdRetrieval(token, rt)
	if err != nil {
		return token, err
	}

	if len(user) == 0 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		unauthorizedError := Response{
			Code:    401,
			Message: "UnauthorizedError: Access token is missing or invalid.",
		}

		notFoundError := BackendError{
			Affinity: "User not found",
			Message:  "User not found for authentication",
			OG_error: nil,
		}
		err = json.NewEncoder(w).Encode(unauthorizedError)

		// Checking that the bad request encoding has gone through successfully
		if err != nil {
			encodingError := BackendError{
				Affinity: "Authentication",
				Message:  "Request encoding for unauthorized error has failed",
				OG_error: err,
			}
			fmt.Println(encodingError.Error())
			return token, &notFoundError
		}
		return token, &notFoundError
	}
	return token, nil
}

func IdRetrieval(id Access_token, rt *_router) ([]string, error) {
	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id.Identifier))
	if err != nil {
		selectionError := BackendError{
			Affinity: "User retrieval",
			Message:  "SELECT in the database seeking users with the same id failed",
			OG_error: err,
		}
		return nil, &selectionError
	}

	// Reading the rows
	users, err := UsersRowReading(rows)

	if err != nil {
		uniquenessError := BackendError{
			Affinity: "User retrieval",
			Message:  "Reading the database rows that were seeking users with the same id failed",
			OG_error: err,
		}
		fmt.Println(uniquenessError.Error())
		return nil, &uniquenessError
	}

	return users, nil
}
