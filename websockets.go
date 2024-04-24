package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	db "green-chat-forum-api/database"
	types "green-chat-forum-api/types"
	util "green-chat-forum-api/util"
)

// type Client struct {
// 	user           *types.User
// 	conn           *websocket.Conn
// 	messageChannel chan []byte
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[int]*types.Client)

func addClient(user types.User, w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorHandler(err)
		return
	}

	client := types.Client{
		User:           &user,
		Conn:           ws,
		MessageChannel: make(chan []byte),
	}
	clients[user.Id] = &client

	go writeMessage(user.Id)
	go readMessages(user.Id)

}

func removeClient(id int) {
	if client, ok := clients[id]; ok {
		client.Conn.Close()
		delete(clients, id)
		fmt.Printf("Deleted: %v\n", id)
	}
}

func readMessages(id int) {
	defer func() {
		removeClient(id)
		broadcastClientsStatus()
	}()
	for {
		_, message, err := clients[id].Conn.ReadMessage()
		if err != nil {
			// Error:  websocket: close 1001 (going away)
			fmt.Println(err, " Connection: ", id)
			return
		}
		fmt.Println("Message ", string(message))
		message = []byte("Using channel. Message: " + string(message))
		for _, conn := range clients {
			conn.MessageChannel <- message
		}
	}
}

func writeMessage(id int) {
	client := clients[id]
	defer func() {
		removeClient(id)
	}()
	for {
		select {
		case message, ok := <-client.MessageChannel:
			if ok {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				if err := client.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func broadcastClientsStatus() {

	for id, client := range clients {

		//Get all users that chatted with current user/id

		chatMates, err := db.GetChatMates(id)

		if err != nil {
			return
		}

		users, err := db.GetUsers()

		if err != nil {
			return
		}

		users = filterUsersByBanned(users)

		//Join chatMates and users
		for _, user := range users {
			if !util.Contains(chatMates, *user) && user.Id != id {
				chatMates = append(chatMates, user)
			}
		}

		//Mark on-line/off-line users
		for _, user := range chatMates {
			//Mark on-line/off-line
			util.SetOnLineStatus(user, clients)
		}

		message := `{"online_users":[`
		for _, user := range chatMates {
			message += fmt.Sprintf(`{"id": %v, "nick_name": "%v", "on_line": "%v"},`, user.Id, user.NickName, user.OnLine)
		}
		message = strings.TrimSuffix(message, ",")
		message += `]}`
		client.MessageChannel <- []byte(message)
	}
}

func notifyClient(id int, message []byte) {
	if client, ok := clients[id]; ok {
		client.MessageChannel <- message
	}
}
