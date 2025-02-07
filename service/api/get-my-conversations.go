package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// It retrieves a list of conversations, be it a private chat or a groupchat, for a user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Logging information
	const affinity string = "Getting list of conversations"

	// Authentication
	token, err := Authentication(w, r, rt)
	if err != nil {
		return
	}

	// Getting the username
	userraw, err := UserFromIdRetrieval(token, rt, w)
	if err != nil {
		return
	}

	user := Username{
		Name: userraw[1],
	}

	// PRIVATE CONVERSATIONS
	var privchats []ChatPreview
	var convID ConversationID
	chats, err := PrivConversationsFromUsernameRetrieval(user, rt, w)
	if err != nil {
		return
	}

	i := 0
	for i < len(chats) {
		rawid, err := strconv.Atoi(chats[i])
		if err != nil {
			_ = createBackendError(affinity, "Converting private chat id to int has failed", err, w, rt)
		}
		convID = ConversationID{rawid}

		// Getting last message
		lastmessage, err := LastMessageFromConvo(convID, rt, w)
		if err != nil {
			return
		}
		// Since conversations are created before the first message is sent, there might be no sent message for the user
		if len(lastmessage) == 0 {
			var emptyComments []Comment
			lastmessage[0] = Message{
				MessageID:  0,
				Timestamp:  "",
				Content:    "",
				Photo:      "",
				Username:   "",
				Checkmarks: -1,
				Comments:   emptyComments,
				ReplyingTo: -1,
			}
		}

		// Adding the actual recipient, being careful it is not the requesting user itself
		var recipient string
		if chats[i+1] == user.Name {
			recipient = chats[i+2]
		} else {
			recipient = chats[i+1]
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

		chatpreview := ChatPreview{
			ChatID:      convID,
			User:        user,
			LastMessage: lastmessage[0],
		}
		privchats = append(privchats, chatpreview)

		i += 3
	}

	// GROUPS
	var groupchats []GroupPreview
	chats, err = GroupConversationsFromUsernameRetrieval(user, rt, w)
	if err != nil {
		return
	}

	i = 0
	for i < len(chats) {
		rawid, err := strconv.Atoi(chats[i])
		if err != nil {
			_ = createBackendError(affinity, "Converting groupchat id to int has failed", err, w, rt)
		}
		convID = ConversationID{rawid}

		// Getting last message
		lastmessage, err := LastMessageFromConvo(convID, rt, w)
		if err != nil {
			return
		}
		// Since conversations are created before the first message is sent, there might be no sent message for the user
		if len(lastmessage) == 0 {
			var emptyComments []Comment
			lastmessage[0] = Message{
				MessageID:  0,
				Timestamp:  "",
				Content:    "",
				Photo:      "",
				Username:   "",
				Checkmarks: -1,
				Comments:   emptyComments,
				ReplyingTo: -1,
			}
		}

		// Getting group info
		info, err := GroupInfoFromIdRetrieval(convID, rt, w)
		if err != nil {
			return
		}

		groupPreview := GroupPreview{
			ChatID:      convID,
			Groupname:   info[1],
			Groupphoto:  info[2],
			LastMessage: lastmessage[0],
		}

		groupchats = append(groupchats, groupPreview)
		i += 2
	}

	// Putting together the two responses
	conversations := Conversations{
		Privchats: privchats,
		Groups:    groupchats,
	}

	// Writing the response in HTTP
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		_ = createBackendError(affinity, "Encoding conversations has failed", err, w, rt)
		return
	}

	// Updating received messages
	err = receivedCheckmarksUpdate(user, w, rt)
	if err != nil {
		return
	}
}
