package api

import (
	"encoding/json"
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

func createBackendError(affinity string, message string, og_error error, responsewriter http.ResponseWriter, rt *_router) *BackendError {
	backenderror := BackendError{
		Affinity:       affinity,
		Message:        message,
		OG_error:       og_error,
		ResponseWriter: responsewriter,
	}

	rt.baseLogger.Println(backenderror.Error())
	return &backenderror
}

// Creates an error message and recovers if another error is produced
func createFaultyResponse(code int, message string, affinity string, failmessage string, w http.ResponseWriter, rt *_router) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	error := Response{
		Code:    code,
		Message: message,
	}
	err := json.NewEncoder(w).Encode(error)

	// Checking that the bad request encoding has gone through successfully
	if err != nil {
		_ = createBackendError(affinity, failmessage, err, w, rt)
	}
}
