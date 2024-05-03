package main

import (
	"fmt"
	database "green-chat-forum-api/database"
	"net/http"
	"os"
)

func main() {

	database.OpenDatabase()
	defer database.CloseDatabase()

	//http.Handle("/", http.FileServer(http.Dir("../")))
	http.HandleFunc("/", ping)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/signout", signoutHandler)
	http.HandleFunc("/newpost", newpostHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/comments", commentsHandler)
	http.HandleFunc("/ws/", websocketHandler)
	http.HandleFunc("/admin/signup", adminSignUpHandler)
	http.HandleFunc("/admin/users", adminUsersHandler)
	http.HandleFunc("/admin/users/", adminUsersHandler)
	http.HandleFunc("/admin/posts", adminPostsHandler)
	http.HandleFunc("/admin/posts/", adminPostsHandler)
	http.HandleFunc("/admin/comments", adminCommentsHandler)
	http.HandleFunc("/admin/comments/", adminCommentsHandler)
	http.HandleFunc("/admin/categories",adminCategoriesHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on port :" + port)
	http.ListenAndServe(":"+port, nil)
}
