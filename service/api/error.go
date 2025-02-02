package api

import (
	"net/http"
)

type BackendError struct {
	Affinity       string
	Message        string
	OG_error       error
	ResponseWriter http.ResponseWriter
}

func (e *BackendError) Error() string {
	title := "ERROR!!\n"
	ast := "***************************************************************\n"
	val := ast + title + "Affinity: " + e.Affinity + "\nMessage: " + e.Message + "\n"
	if e.OG_error != nil {
		val += "Original error: " + e.OG_error.Error() + "\n"
	}
	/* What was this?
	if e.ResponseWriter != nil {

	}
	*/
	return val + ast
}
