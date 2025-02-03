package api

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
