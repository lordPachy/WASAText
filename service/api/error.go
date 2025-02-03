package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BackendError struct {
	Affinity       string
	Message        string
	OG_error       error
	ResponseWriter http.ResponseWriter
}

func (e BackendError) Error() string {
	title := "ERROR!!\n"
	ast := "***************************************************************\n"
	val := ast + title + "Affinity: " + e.Affinity + "\nMessage: " + e.Message + "\n"
	if e.OG_error != nil {
		val += "Original error: " + e.OG_error.Error() + "\n"
	}

	if e.ResponseWriter != nil {
		w := e.ResponseWriter
		w.WriteHeader(http.StatusBadRequest)
		badRequest := Response{
			Code:    500,
			Message: "Internal server error",
		}
		_ = json.NewEncoder(w).Encode(badRequest)
	}

	return val + ast
}

func createBackendError(affinity string, message string, og_error error, responsewriter http.ResponseWriter) *BackendError {
	backenderror := BackendError{
		Affinity:       affinity,
		Message:        message,
		OG_error:       og_error,
		ResponseWriter: responsewriter,
	}

	fmt.Println(backenderror.Error())
	return &backenderror
}
