package api

var ConversationNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation id given was not found.",
}

var ConversationOrMessageNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation/Message id given was not found.",
}

var ConversationOrMessageOrCommentNotFound = Response{
	Code:    404,
	Message: "Not found: Conversation/Message/Comment id given was not found.",
}
