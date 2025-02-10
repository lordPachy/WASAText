package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
)

// It retrieves a single conversation, be it a private chat or a groupchat, for a user.
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Getting conversation"

	// Checking that the conversation actually exists
	convIDValue, err := conversationRetrieval(affinity, w, ps, rt)
	if convIDValue == -1 || err != nil {
		return
	}
	convID := ConversationID{convIDValue}

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting messages
	messages, err := MessagesFromConvo(convID, rt, w)
	if err != nil {
		return
	}

	// Sorting messages
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp > messages[j].Timestamp
	})

	// Retrieving the sender's username
	senderUsername, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	// Responding for a private chat
	if convID.Id < 5000 {
		// Getting recipient
		recipients, err := ConversationFromIdRetrieval(convID.Id, rt, w)
		if err != nil {
			return
		}

		// Adding the actual recipient, being careful it is not the requesting user itself
		var recipient string
		if recipients[1] == senderUsername[1] {
			recipient = recipients[2]
		} else {
			recipient = recipients[1]
		}

		// Retrieving recipient's information
		userraw, err := UserFromUsernameRetrieval(Username{recipient}, rt, w)
		if err != nil {
			return
		}
		user := User{
			Username: userraw[1],
			Propic:   userraw[2],
		}

		response := Chat{
			ConversationID: convID,
			User:           user,
			Messages:       messages,
		}

		// Writing the response in HTTP
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			_ = createBackendError(affinity, "Encoding the private chat has failed", err, w, rt)
			return
		}
	} else {
		// Responding for groups

		// Getting recipients
		recipients, err := UsersInGroup(convID, rt, w)
		if err != nil {
			return
		}

		var groupusers []User
		for _, recipient := range recipients {
			// Retrieving recipients' information
			userraw, err := UserFromUsernameRetrieval(Username{recipient}, rt, w)
			if err != nil {
				return
			}
			user := User{
				Username: userraw[1],
				Propic:   userraw[2],
			}

			groupusers = append(groupusers, user)
		}

		// Getting group info
		info, err := GroupInfoFromIdRetrieval(convID, rt, w)
		if err != nil {
			return
		}

		response := Group{
			ConversationID: convID,
			Members:        groupusers,
			Messages:       messages,
			Groupname:      info[1],
			Groupphoto:     info[2],
		}

		// Writing the response in HTTP
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			_ = createBackendError(affinity, "Encoding the group has failed", err, w, rt)
			return
		}
	}

	// Updating read messages
	err = readCheckmarksUpdate(Username{senderUsername[1]}, convID, w, rt)
	if err != nil {
		return
	}
}
