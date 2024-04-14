package main

import (
	"fmt"
	"net/http"
	database "green-chat-forum-api/database"
)

func main() {

	database.CreateDatabase()

	http.Handle("/", http.FileServer(http.Dir("../")))
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/signout", signoutHandler)
	http.HandleFunc("/newpost", newpostHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/comments", commentsHandler)
	http.HandleFunc("/ws/", websocketHandler)
	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)
}
