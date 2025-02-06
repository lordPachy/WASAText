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
