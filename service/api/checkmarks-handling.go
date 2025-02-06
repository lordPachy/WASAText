package api

import (
	"fmt"
	"net/http"
)

// It updates checkmarks after sending a message in a group.
func groupMessageCheckmarksUpdate(groupID ConversationID, messageID MessageID, user Username, w http.ResponseWriter, rt *_router) error {
	// Logging information
	affinity := "Updating checkmarks for group messages"

	// Since this function get only called inside other http requests, we are sure
	// that groupID and messageID are valid. Thus, no check is going to be performed
	// on those.

	// Retrieving group members
	usernames, err := UsersInGroup(groupID, rt, w)
	if err != nil {
		return err
	}

	// Actually adding the membership in the DB
	for _, el := range usernames {
		if el == user.Name {
			continue
		}
		query := fmt.Sprintf("(%d, %d, '%s', %d)", groupID.Id, messageID.Id, el, 0)

		_, err = rt.db.Insert("groupmessageschecks", query)
		if err != nil {
			_ = createBackendError(affinity, "Inserting group message checkmarks into the database has failed", err, w, rt)
			return err
		}
	}

	return nil
}

// It updates checkmarks after updating all conversations
func receivedCheckmarksUpdate(user Username, w http.ResponseWriter, rt *_router) error {
	// Logging information
	affinity := "Updating checkmarks with received"

	// Since this function get only called inside other http requests, we are sure
	// that username is valid. Thus, no check is going to be performed
	// on that.

	// Group messages update
	query := fmt.Sprintf("member = '%s' AND checkmarks = 0", user.Name)

	_, err := rt.db.Update("groupmessageschecks", query, "checkmarks = 1")
	if err != nil {
		_ = createBackendError(affinity, "Updating group message checkmarks into the database has failed", err, w, rt)
		return err
	}

	// Rendering group messages as received if they are received by everyone
	recv_messages, err := ReceivedGroupMessages(rt, w)
	if err != nil {
		return err
	}

	for _, mess := range recv_messages {
		query = fmt.Sprintf("id = %s AND checkmarks = 0", mess)

		_, err := rt.db.Update("messages", query, "checkmarks = 1")
		if err != nil {
			_ = createBackendError(affinity, "Updating group messages checkmarks into the database has failed", err, w, rt)
			return err
		}
	}

	// Private messages update
	messageids, err := MessagesFromPrivateConvo(user, rt, w)
	if err != nil {
		return err
	}

	for _, mess := range messageids {
		// We need to update messages received by the user, not the ones sent
		query = fmt.Sprintf("sender NOT IN ('%s') AND id = %s AND checkmarks = 0", user.Name, mess)

		_, err := rt.db.Update("messages", query, "checkmarks = 1")
		if err != nil {
			_ = createBackendError(affinity, "Updating private message checkmarks into the database has failed", err, w, rt)
			return err
		}
	}

	return nil
}

// It updates checkmarks after updating all conversations
func readCheckmarksUpdate(user Username, convID ConversationID, w http.ResponseWriter, rt *_router) error {
	// Logging information
	affinity := "Updating checkmarks with read"

	// Since this function get only called inside other http requests, we are sure
	// that username is valid. Thus, no check is going to be performed
	// on that.

	if convID.Id >= 5000 {
		// Group messages update
		query := fmt.Sprintf("groupid = %d AND member = '%s' AND checkmarks < 2", convID.Id, user.Name)

		_, err := rt.db.Update("groupmessageschecks", query, "checkmarks = 2")
		if err != nil {
			_ = createBackendError(affinity, "Updating group message checkmarks into the database has failed", err, w, rt)
			return err
		}

		// Rendering group messages as received if they are received by everyone
		recv_messages, err := ReadGroupMessages(rt, w)
		if err != nil {
			return err
		}

		for _, mess := range recv_messages {
			query = fmt.Sprintf("id = %s AND checkmarks < 2", mess)

			_, err := rt.db.Update("messages", query, "checkmarks = 2")
			if err != nil {
				_ = createBackendError(affinity, "Updating group messages checkmarks into the database has failed", err, w, rt)
				return err
			}
		}

	} else {
		// Private messages update
		messageids, err := MessagesFromPrivateConvo(user, rt, w)
		if err != nil {
			return err
		}

		for _, mess := range messageids {
			// We need to update messages received by the user, not the ones sent
			query := fmt.Sprintf("sender NOT IN ('%s') AND id = %s AND checkmarks < 2", user.Name, mess)

			_, err := rt.db.Update("messages", query, "checkmarks = 2")
			if err != nil {
				_ = createBackendError(affinity, "Updating private message checkmarks into the database has failed", err, w, rt)
				return err
			}
		}
	}
	return nil
}
