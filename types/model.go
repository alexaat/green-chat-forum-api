package types

import "github.com/gorilla/websocket"

type User struct {
	Id        int         `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Age       int         `json:"age,string"`
	Gender    string      `json:"gender"`
	NickName  string      `json:"nick_name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Password2 string      `json:"password2"`
	SessionId string      `json:"session_id"`
	OnLine    bool        `json:"on_line"`
	Status    interface{} `json:"status"`
}

type Post struct {
	Id               int      `json:"id"`
	Date             int      `json:"date"`
	UserId           int      `json:"user_id"`
	NickName         string   `json:"nick_name"`
	Content          string   `json:"content"`
	Categories       []string `json:"categories"`
	NumberOfComments int      `json:"number_of_comments"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Payload interface{} `json:"payload"`
	Error   *Error      `json:"error"`
}

type Data struct {
	User  *User   `json:"user"`
	Posts *[]Post `json:"posts"`
}

type Chat struct {
	UserId     int        `json:"user_id"`
	ChatMateId int        `json:"chat_mate_id"`
	Messages   *[]Message `json:"messages"`
	Error      *Error     `json:"error"`
}

type Connection struct {
	User *User
	Conn *websocket.Conn
}

type MessageWrapper struct {
	Message Message `json:"message"`
}

type Message struct {
	Id           int    `json:"id"`
	FromId       int    `json:"from_id"`
	FromNickName string `json:"from_nick_name"`
	ToId         int    `json:"to_id"`
	Content      string `json:"content"`
	Date         int64  `json:"date"`
}

type Comment struct {
	Id           int    `json:"id"`
	Date         int    `json:"date"`
	UserId       int    `json:"user_id"`
	UserNickName string `json:"user_nick_name"`
	PostId       int    `json:"post_id"`
	Content      string `json:"content"`
}

type CommentsPageObject struct {
	User     *User      `json:"user"`
	Post     *Post      `json:"post"`
	Comments []*Comment `json:"comments"`
}

type NewPostPageObject struct {
	User       *User  `json:"user"`
	Categories string `json:"categories"`
}

type Client struct {
	User           *User
	Conn           *websocket.Conn
	MessageChannel chan []byte
}

type RowsAffected struct {
	RowsAffected int64 `json:"rows-affected"`
}
