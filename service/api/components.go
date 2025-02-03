package api

import (
	"encoding/json"
	"net/http"
)

type Username struct {
	Name string `json:"name"`
}

type Image struct {
	Image string `json:"image"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Creates an error message and recovers if another error is produced
func createFaultyResponse(code int, message string, affinity string, failmessage string, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	error := Response{
		Code:    code,
		Message: message,
	}
	err := json.NewEncoder(w).Encode(error)

	// Checking that the bad request encoding has gone through successfully
	if err != nil {
		_ = createBackendError(affinity, failmessage, err, w)
	}
}

type User struct {
	Username string `json:"username"`
	Propic   string `json:"propic"`
}

type Access_token struct {
	Identifier string `json:"identifier"`
}

type Comment struct {
	CommentID int    `json:"commentid"`
	Sender    string `json:"sender"`
	Reaction  string `json:"reaction"`
	SentByMe  bool   `json:"sentbyme"`
}

// Message sent with HTTP request
type RequestMessage struct {
	Content    string `json:"content"`
	Photo      string `json:"photo"`
	ReplyingTo int    `json:"replyingto"`
}

// Message held in the database
type Message struct {
	MessageID  int       `json:"messageid"`
	Timestamp  string    `json:"timestamp"`
	Content    string    `json:"content"`
	Photo      string    `json:"string"`
	Username   string    `json:"username"`
	Checkmarks int       `json:"checkmarks"`
	Comments   []Comment `json:"comments"`
	ReplyingTo int       `json:"replyingto"`
}
