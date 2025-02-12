package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.PUT("/session", rt.createUser)
	rt.router.POST("/session", rt.doLogin)

	// Settings methods
	rt.router.PUT("/settings/username", rt.setMyUserName)
	rt.router.PUT("/settings/profilepicture", rt.setMyPhoto)

	// Conversation methods
	rt.router.GET("/users", rt.getUsers)
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.PUT("/conversations", rt.createConversation)
	rt.router.GET("/conversations/:conversationid", rt.getConversation)
	rt.router.POST("/conversations/:conversationid", rt.sendMessage)
	rt.router.PUT("/groups", rt.addToGroup)
	rt.router.DELETE("/conversations/:conversationid", rt.leaveGroup)

	rt.router.PUT("/conversations/:conversationid/settings/groupname", rt.setGroupName)
	rt.router.PUT("/conversations/:conversationid/settings/grouphoto", rt.setGroupPhoto)

	// Message methods
	rt.router.PUT("/conversations/:conversationid/messages/:messageid", rt.commentMessage)
	rt.router.POST("/conversations/:conversationid/messages/:messageid", rt.forwardMessage)
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid", rt.deleteMessage)
	rt.router.DELETE("/conversations/:conversationid/messages/:messageid/comments/:commentid", rt.uncommentMessage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
