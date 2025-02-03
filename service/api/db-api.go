package api

import (
	"fmt"
	"net/http"
)

/*
This package contains function that get never called directly
but always through other functions. Thus, they get to create
their own errors, and manage the writer.
*/

// Retrieve array of strings (that should represent a single user) from id
func IdRetrieval(id Access_token, rt *_router, w http.ResponseWriter) ([]string, error) {
	// Logging information
	affinity := "User retrieval"
	// SQL query
	rows, err := rt.db.Select("*", "users", fmt.Sprintf("id = '%s'", id.Identifier))
	if err != nil {
		return nil, createBackendError(affinity, "SELECT in the database seeking users with the same id failed", err, w)
	}

	// Reading the rows
	users, err := UsersRowReading(rows)

	if err != nil {
		return nil, createBackendError(affinity, "Reading the database rows that were seeking users with the same id failed", err, w)
	}

	return users, nil
}

// Retrieve array of strings (that should represent a single user) from username
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
