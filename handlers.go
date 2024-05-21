package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	crypto "green-chat-forum-api/crypto"
	db "green-chat-forum-api/database"
	types "green-chat-forum-api/types"
	util "green-chat-forum-api/util"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Green Chat Forum API")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	session_id := strings.TrimPrefix(r.URL.Path, "/ws/")

	user, err := db.GetUserBySessionId(session_id)
	if err != nil {
		return
	}
	if user == nil {
		return
	}

	addClient(*user, w, r)
	broadcastClientsStatus()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	if r.Method == "GET" {

		keys, ok := r.URL.Query()["session_id"]
		if !ok || len(keys[0]) < 1 {
			resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
			sendResponse(w, resp)
			return
		}
		session_id := keys[0]

		user, err := db.GetUserBySessionId(session_id)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

		if user != nil && user.Status != nil && user.Status.(string) == "banned" {
			resp.Payload = nil
			sendResponse(w, resp)
			return
		}

		if user != nil {
			posts, err := db.GetPosts(user)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}

			posts, err = filterPostsByBannedUser(*posts)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}

			util.RemoveUserInfo(user)
			data := types.Data{Posts: posts, User: user}
			resp.Payload = data

		}
	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	user_name := r.FormValue("user_name")
	password := r.FormValue("password")

	user, e := db.GetUserByEmailOrNickNameAndPassword(types.User{NickName: user_name, Password: password})

	if e != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: unable access database: %v", e)}
	} else {
		if user == nil {
			resp.Error = &types.Error{Type: util.NO_USER_FOUND, Message: "Error: no such user"}
		} else {

			user.SessionId = generateSessionId()
			err := db.UpdateSessionId(user)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: unable access database: %v", err)}
				sendResponse(w, resp)
				return
			}

			posts, err := db.GetPosts(user)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: unable access database: %v", err)}
				sendResponse(w, resp)
				return
			}

			util.RemoveUserInfo(user)
			data := types.Data{Posts: posts, User: user}
			resp.Payload = data
		}
	}
	sendResponse(w, resp)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	data := types.Data{User: &types.User{}, Posts: nil}

	data.User.FirstName = strings.TrimSpace(r.FormValue("first_name"))
	data.User.LastName = strings.TrimSpace(r.FormValue("last_name"))
	data.User.NickName = strings.TrimSpace(r.FormValue("nick_name"))
	data.User.Email = strings.TrimSpace(r.FormValue("email"))
	data.User.Gender = strings.TrimSpace(r.FormValue("gender"))
	data.User.Password = r.FormValue("password")
	data.User.Password2 = r.FormValue("password2")

	age_str := strings.TrimSpace(r.FormValue("age"))
	resp.Error = util.ValidateInput(data.User, age_str)

	if resp.Error == nil {
		// Try to insert User
		data.User.SessionId = generateSessionId()
		data.User.Password = crypto.Encrypt(data.User.Password)
		data.User.Password2 = ""
		id, err := db.SaveUser(data.User)
		if err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed: users.nick_name") {
				resp.Error = &types.Error{Type: util.INVALID_NICK_NAME, Message: "Error: nick name is already in use"}
				resp.Payload = nil
			} else if strings.HasPrefix(err.Error(), "UNIQUE constraint failed: users.email") {
				resp.Error = &types.Error{Type: util.INVALID_EMAIL, Message: "Error: email is already in use"}
				resp.Payload = nil
			} else {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				resp.Payload = nil
			}
		} else {
			data.User.Id = int(id)
		}

	}
	if resp.Error == nil {

		posts, err := db.GetPosts(data.User)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: unable access database: %v", err)}
			resp.Payload = nil
			sendResponse(w, resp)
			return
		}
		data.Posts = posts
		util.RemoveUserInfo(data.User)
		resp.Payload = data
	}
	sendResponse(w, resp)
}

func signoutHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	session_id := r.FormValue("session_id")

	u, err := db.GetUserBySessionId(session_id)
	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	err = db.ResetSessionId(session_id)

	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	if u != nil {
		removeClient(u.Id)
		broadcastClientsStatus()
	}
	sendResponse(w, resp)
}

func newpostHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]

	//Verify User
	user, err := db.GetUserBySessionId(session_id)
	if user == nil {
		resp.Error = &types.Error{Type: util.NO_USER_FOUND, Message: fmt.Sprintf("Error: unable to authorize user: %v", err)}
		sendResponse(w, resp)
		return
	}

	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	if user.Status == "banned" {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: Banned by administrator"}
		sendResponse(w, resp)
		return
	}
	util.RemoveUserInfo(user)

	if r.Method == "GET" {

		//Get Categories
		categories, err := db.GetCategories()

		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

		j, err := json.Marshal(categories)

		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: unable to parse data %v", err)}
			sendResponse(w, resp)
			return
		}

		npop := types.NewPostPageObject{}
		npop.User = user
		npop.Categories = string(j)

		resp.Payload = npop
	} else if r.Method == "POST" {
		//session_id := r.FormValue("session_id")
		content := strings.TrimSpace(r.FormValue("content"))
		categories := r.FormValue("categories")

		//0. Validate content
		if len(content) == 0 {
			resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Empty post is not allowed"}
			sendResponse(w, resp)
			return
		}

		if len(content) > 10000 {
			resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Post is too large"}
			sendResponse(w, resp)
			return
		}

		// 2. Insert Post
		var arr []string
		err = json.Unmarshal([]byte(categories), &arr)

		if err != nil {
			return
		}

		post := types.Post{
			UserId:     user.Id,
			Content:    content,
			Categories: arr,
		}
		err = db.InsertPost(user, &post)

		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	if r.Method != "POST" {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: fmt.Sprintf("Error: %v", "Wrong method used")}
		sendResponse(w, resp)
		return
	}

	session_id := r.FormValue("session_id")
	to_id := r.FormValue("to_id")
	message := strings.TrimSpace(r.FormValue("message"))

	//Verify input
	if len(message) == 0 {
		resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Empty message is not allowed"}
		sendResponse(w, resp)
		return
	}

	if len(message) > 1000 {
		resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Message is too large"}
		sendResponse(w, resp)
		return
	}

	user, err := db.GetUserBySessionId(session_id)
	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}
	if user == nil {
		sendResponse(w, resp)
		return
	}

	to_id_int, err := strconv.Atoi(to_id)
	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	util.RemoveUserInfo(user)

	resp.Payload = user

	m := types.Message{
		FromId:       user.Id,
		FromNickName: user.NickName,
		ToId:         to_id_int,
		Content:      message,
		Date:         util.GetCurrentMilli(),
	}

	err = db.InsertMessage(m)
	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	mw := types.MessageWrapper{Message: m}

	b, err := json.Marshal(mw)

	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	//Notify both sender and receiver
	notifyClient(m.FromId, b)
	notifyClient(m.ToId, b)

	sendResponse(w, resp)
}

func commentsHandler(w http.ResponseWriter, r *http.Request) {

	resp := types.Response{Payload: nil, Error: nil}

	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]

	//Verify session_id
	user, err := db.GetUserBySessionId(session_id)
	if user == nil {
		resp.Error = &types.Error{Type: util.NO_USER_FOUND, Message: fmt.Sprintf("Error: unable to authorize user: %v", err)}
		sendResponse(w, resp)
		return
	}

	if err != nil {
		resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		sendResponse(w, resp)
		return
	}

	if user.Status == "banned" {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: Banned by administrator"}
		sendResponse(w, resp)
		return
	}

	util.RemoveUserInfo(user)

	if r.Method == "GET" {

		keys, ok = r.URL.Query()["post_id"]
		if !ok || len(keys[0]) < 1 {
			resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: post_id"}
			sendResponse(w, resp)
			return
		}
		post_id := keys[0]

		// Get Post by post_id
		postId, err := strconv.Atoi(post_id)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		post, err := db.GetPost(postId)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

		//4. Get Comments
		comments, err := db.GetComments(postId)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

		//Filter by banned user
		comments, err = filterCommentsByUserId(comments)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

		cpo := types.CommentsPageObject{}
		cpo.User = user
		cpo.Post = post
		cpo.Comments = comments
		resp.Payload = cpo

	} else if r.Method == "POST" {
		//session_id := r.FormValue("session_id")

		post_id := r.FormValue("post_id")
		comment := strings.TrimSpace(r.FormValue("comment"))

		//Verify comment
		if len(comment) == 0 {
			resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Empty comment is not allowed"}
			sendResponse(w, resp)
			return
		}

		if len(comment) > 10000 {
			resp.Error = &types.Error{Type: util.INVALID_INPUT, Message: "Comment is too large"}
			sendResponse(w, resp)
			return
		}

		postId, err := strconv.Atoi(post_id)

		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: unable to parse data %v", err)}
			sendResponse(w, resp)
			return
		}

		err = db.SaveComment(user.Id, postId, comment)

		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}
	sendResponse(w, resp)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	chat := types.Chat{UserId: -1, ChatMateId: -1, Messages: nil, Error: nil}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != "POST" {
		chat.Error = &types.Error{Type: util.WRONG_METHOD, Message: fmt.Sprintf("Error: %v", "Wrong method used")}

		json.NewEncoder(w).Encode(chat)
		return
	}

	session_id := r.FormValue("session_id")

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		chat.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
		json.NewEncoder(w).Encode(chat)
		return
	}
	if page <= 0 {
		page = 1
	}

	chat_mate_id, err := strconv.Atoi(r.FormValue("chat_mate_id"))
	if err != nil {
		chat.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
		json.NewEncoder(w).Encode(chat)
		return
	}

	chat.ChatMateId = chat_mate_id

	user, err := db.GetUserBySessionId(session_id)
	if err != nil {
		chat.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		json.NewEncoder(w).Encode(chat)
		return
	}

	chat.UserId = user.Id

	messages, err := db.GetChat(user.Id, chat_mate_id, page)

	if err != nil {
		chat.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
		json.NewEncoder(w).Encode(chat)
		return
	}

	chat.Messages = messages

	json.NewEncoder(w).Encode(chat)
}

func errorHandler(err error) {
	fmt.Println("Error: ", err)
}

func sendResponse(w http.ResponseWriter, resp types.Response) {
	//host := "http://localhost:8000"
	host := "http://alexaat.com"
	w.Header().Set("Access-Control-Allow-Origin", host)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	json.NewEncoder(w).Encode(resp)
}
