package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// It comments an existing message.
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Message sending"

	// Checking that the conversation actually exists
	convID, err := conversationRetrieval(affinity, w, ps, rt)
	if err != nil || convID == -1 {
		return
	}

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Checking that the message actually exists
	messID, err := messageRetrieval(affinity, w, ps, rt)
	if err != nil || messID == -1 {
		return
	}

	// Retrieving the username
	user, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	// Getting the new comment
	var newComment CommentRequest
	err = json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		createFaultyResponse(http.StatusBadRequest, "The received body is not a comment", affinity, "Request encoding for badly formatted comment response has failed", w, rt)
		return
	}

	// Checking if the comment is valid
	match, err := checkCommentCorrectness(newComment)
	if err != nil {
		return
	}

	if !match {
		w.WriteHeader(http.StatusBadRequest)
		createFaultyResponse(http.StatusBadRequest, "Comment parsed incorrectly or not valid", affinity, "Request encoding for comment not correcly formatted response has failed", w, rt)
		return
	}

	// Creating the id
	id, err := CommentIdCreator(rt, w)
	if err != nil {
		return
	}

	// Actually writing the message in the DB
	query := fmt.Sprintf("(%d, '%s', '%s')", id, user[1], newComment.Reaction)

	_, err = rt.db.Insert("comments", query)
	if err != nil {
		_ = createBackendError(affinity, "Inserting the new comment into the database has failed", err, w, rt)
		return
	}

	// Writing the comment into the comment recording tables
	err = update_comments(CommentID{id}, MessageID{Id: messID}, w, rt)
	if err != nil {
		return
	}

	// Writing the response in HTTP
	// Accepted request
	w.Header().Set("content-type", "application-json")
	w.WriteHeader(http.StatusCreated)
	responseComment := Comment{
		CommentID: id,
		Sender:    user[1],
		Reaction:  newComment.Reaction,
		SentByMe:  true,
	}
	err = json.NewEncoder(w).Encode(responseComment)
	if err != nil {
		_ = createBackendError(affinity, "Encoding the new comment has failed", err, w, rt)
		return
	}
}

// A function that checks the correctness of every field in a sent comment
func checkCommentCorrectness(newComment CommentRequest) (bool, error) {
	// Checking that the comment is valid
	for _, el := range []string{"laugh", "sad", "thumbs_up", "surprised", "love", "pray"} {
		if el == newComment.Reaction {
			return true, nil
		}
	}

	return false, nil
}

// It creates a numerical ID for the new comment.
func CommentIdCreator(rt *_router, w http.ResponseWriter) (int, error) {
	var id int

	// Logging information
	const affinity string = "Comment posting"

	for {
		id = rand.Intn(10001)
		rows, err := rt.db.Select("*", "comments", fmt.Sprintf("id = %d", id))
		if err != nil {
			return 0, createBackendError(affinity, "SELECT in the database seeking comments with the same id failed", err, w, rt)
		}

		// Checking that the new id is unique
		other_comments, err := CommentsRowReading(rows)

		if err != nil {
			return 0, createBackendError(affinity, "Reading the database rows that were seeking comments with the same id failed", err, w, rt)
		} else if len(other_comments) == 0 {
			break
		}
	}

	return id, nil
}
