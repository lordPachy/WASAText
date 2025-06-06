package api

type Username struct {
	Name string `json:"name"`
}

type Image struct {
	Image string `json:"image"`
}

type Preview struct {
	ChatID      ConversationID `json:"chatid"`
	Name        string         `json:"name"`
	Photo       string         `json:"photo"`
	LastMessage Message        `json:"lastmessage"`
}

type ConversationRequest struct {
	IsGroup   bool       `json:"isgroup"`
	Members   []Username `json:"members"`
	GroupName string     `json:"groupname"`
}

type ConversationID struct {
	Id int `json:"id"`
}

type GroupName struct {
	Value string `json:"value"`
}

type AddToGroupRequest struct {
	User  Username       `json:"username"`
	Group ConversationID `json:"group"`
}

type MessageID struct {
	Id int `json:"id"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Username string `json:"username"`
	Propic   string `json:"propic"`
}

type Access_token struct {
	Identifier string `json:"identifier"`
}

type CommentID struct {
	CommentID int `json:"commentid"`
}

type Comment struct {
	CommentID int    `json:"commentid"`
	Sender    string `json:"sender"`
	Reaction  string `json:"reaction"`
}

type CommentRequest struct {
	Reaction string `json:"reaction"`
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
	Photo      string    `json:"photo"`
	Username   string    `json:"username"`
	Checkmarks int       `json:"checkmarks"`
	Comments   []Comment `json:"comments"`
	ReplyingTo int       `json:"replyingto"`
	OG_Sender  string    `json:"og_sender"`
}

type Chat struct {
	ConversationID ConversationID `json:"conversationid"`
	User           User           `json:"user"`
	Messages       []Message      `json:"messages"`
}

type Group struct {
	ConversationID ConversationID `json:"conversationid"`
	Members        []User         `json:"members"`
	Messages       []Message      `json:"messages"`
	Groupname      string         `json:"groupname"`
	Groupphoto     string         `json:"groupphoto"`
}
