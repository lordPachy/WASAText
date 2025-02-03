package api

import (
	"net/http"
)

/*
Authenticate the user and return its id. This function manages its own

errors and needs no error handling.
*/
func Authentication(w http.ResponseWriter, r *http.Request, rt *_router) (Access_token, error) {
	// Logging information
	affinity := "authentication"

	// ID retrieval
	id := r.Header.Get("Authentication")
	token := Access_token{
		Identifier: id,
	}

	user, err := IdRetrieval(token, rt, w)
	if err != nil {
		return token, err
	}

	if len(user) == 0 {
		notFoundError := BackendError{
			Affinity: "User not found",
			Message:  "User not found for authentication",
			OG_error: nil,
		}

		createFaultyResponse(http.StatusUnauthorized, "UnauthorizedError: Access token is missing or invalid.", affinity, "Unauthorized error encoding has failed", w)

		return token, &notFoundError
	}

	return token, nil
}
