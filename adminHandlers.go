package main

import (
	"encoding/json"
	"fmt"
	db "green-chat-forum-api/database"
	types "green-chat-forum-api/types"
	util "green-chat-forum-api/util"
	"net/http"
	"strconv"
	"strings"
)

func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

	//Verify admin
	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]
	user, err := db.GetAdminBySessionId(session_id)
	if err != nil || user == nil {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: cannot verify admin"}
		sendResponse(w, resp)
		return
	}

	//Get user id from url
	userId := 0
	if strings.Contains(r.URL.Path, "/admin/users/") {
		userIdStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/admin/users/"))

		if userIdStr != "" {
			userId, err = strconv.Atoi(userIdStr)
			if err != nil || userId < 1 {
				resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Cannot parse user id: %v. %v", userIdStr, err)}
				sendResponse(w, resp)
			}
		}
	}

	if r.Method == "GET" {
		if userId == 0 {
			//No user id. Get All users
			users, err := db.GetUsers()
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = users
		} else {
			user, err := db.GetUserById(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = user
		}

	} else if r.Method == "DELETE" {

		if userId > 0 {

			//Get Posts
			posts, err := db.GetPostsByUserId(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot read database. %v", err)}
				sendResponse(w, resp)
				return
			}

			//Delete user
			num1, err := db.DeleteUserById(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
				sendResponse(w, resp)
				return
			}

			//Delete posts
			num2, err := db.DeletePostsByUserId(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
				sendResponse(w, resp)
				return
			}

			//Delete comments
			num3, err := db.DeleteCommentsByUserId(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
				sendResponse(w, resp)
				return
			}

			var num5 int64 = 0

			for _, post := range *posts {
				n, err := db.DeleteCommentByPostId(post.Id)
				if err != nil {
					resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
					sendResponse(w, resp)
					return
				}
				num5 += *n
			}

			//Delete chat messages
			num4, err := db.DeleteMessagesByUserId(userId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
				sendResponse(w, resp)
				return
			}

			resp.Payload = types.RowsAffected{
				RowsAffected: *num1 + *num2 + *num3 + *num4 + num5,
			}
		}

	} else if r.Method == "PATCH" {
		if userId > 0 {
			ban := r.FormValue("ban")
			value := ""
			if ban == "true" {
				value = "banned"
			}
			num, err := db.UpdateUserStatus(userId, value)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
				sendResponse(w, resp)
				return
			}

			broadcastClientsStatus()

			resp.Payload = types.RowsAffected{RowsAffected: *num}

		}

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}

func adminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

	if r.Method == "POST" {
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		//Validate
		admin, err := db.GetAdminByEmailAndPassword(email, password)

		if err != nil || admin == nil {
			resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: fmt.Sprintf("Cannot get admin from database. %v", err)}
			sendResponse(w, resp)
			return
		}

		(*admin).SessionId = generateSessionId()

		err = db.UpdateAdminSessionId(admin)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
			sendResponse(w, resp)
			return
		}

		resp.Payload = admin

	} else if r.Method == "DELETE" {
		//Verify admin
		keys, ok := r.URL.Query()["session_id"]
		if !ok || len(keys[0]) < 1 {
			resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
			sendResponse(w, resp)
			return
		}
		session_id := keys[0]
		admin, err := db.GetAdminBySessionId(session_id)
		if err != nil || admin == nil {
			resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: cannot verify admin"}
			sendResponse(w, resp)
			return
		}

		(*admin).SessionId = ""

		err = db.UpdateAdminSessionId(admin)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Cannot update database. %v", err)}
			sendResponse(w, resp)
			return
		}

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}
	sendResponse(w, resp)
}

func adminPostsHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

	//Verify admin
	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]
	user, err := db.GetAdminBySessionId(session_id)
	if err != nil || user == nil {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: cannot verify admin"}
		sendResponse(w, resp)
		return
	}

	//Get post id from url
	postId := 0
	if strings.Contains(r.URL.Path, "/admin/posts/") {
		postIdStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/admin/posts/"))
		if postIdStr != "" {
			postId, err = strconv.Atoi(postIdStr)
			if err != nil || postId < 1 {
				resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Cannot parse user id: %v. %v", postIdStr, err)}
				sendResponse(w, resp)
			}
		}
	}

	if r.Method == "GET" {
		if postId == 0 {
			//No post id. Get All posts
			posts, err := db.GetAllPosts()
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = posts
		} else {
			post, err := db.GetPost(postId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = post
		}

	} else if r.Method == "DELETE" {
		if postId > 0 {
			num, err := db.DeletePost(postId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = types.RowsAffected{
				RowsAffected: *num,
			}
		}

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}

func adminCommentsHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

	//Verify admin
	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]
	user, err := db.GetAdminBySessionId(session_id)
	if err != nil || user == nil {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: cannot verify admin"}
		sendResponse(w, resp)
		return
	}

	//Get comment id from url
	commentId := 0
	if strings.Contains(r.URL.Path, "/admin/comments/") {
		commentsIdStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/admin/comments/"))
		if commentsIdStr != "" {
			commentId, err = strconv.Atoi(commentsIdStr)
			if err != nil || commentId < 1 {
				resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Cannot parse user id: %v. %v", commentsIdStr, err)}
				sendResponse(w, resp)
			}
		}
	}

	if r.Method == "GET" {
		if commentId == 0 {
			//Get All Comments
			comments, err := db.GetAllComments()
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = comments
		} else {
			comments, err := db.GetCommentById(commentId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = comments

		}

	} else if r.Method == "DELETE" {
		if commentId > 0 {
			num, err := db.DeleteCommentById(commentId)
			if err != nil {
				resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
				sendResponse(w, resp)
				return
			}
			resp.Payload = types.RowsAffected{
				RowsAffected: *num,
			}
		}

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}

func adminCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

	//Verify admin
	keys, ok := r.URL.Query()["session_id"]
	if !ok || len(keys[0]) < 1 {
		resp.Error = &types.Error{Type: util.MISSING_PARAM, Message: "Error: missing request parameter: session_id"}
		sendResponse(w, resp)
		return
	}
	session_id := keys[0]
	user, err := db.GetAdminBySessionId(session_id)
	if err != nil || user == nil {
		resp.Error = &types.Error{Type: util.AUTHORIZATION, Message: "Error: cannot verify admin"}
		sendResponse(w, resp)
		return
	}

	if r.Method == "GET" {
		//Get All Categories
		categories, err := db.GetCategories()
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		resp.Payload = categories

	} else if r.Method == "POST" {
		//Categories as json ["apple", "pear", "lemon"]
		categories := r.FormValue("categories")
		var arr []string
		err = json.Unmarshal([]byte(categories), &arr)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		//1. Remove categories
		_, err := db.DeleteCategories()
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		//2. Post Categories
		_, err = db.InsertCategories(arr)
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		//3. Get Updated Categories
		arr, err = db.GetCategories()
		if err != nil {
			resp.Error = &types.Error{Type: util.ERROR_ACCESSING_DATABASE, Message: fmt.Sprintf("Error: %v", err)}
			sendResponse(w, resp)
			return
		}
		resp.Payload = arr

	} else {
		resp.Error = &types.Error{Type: util.WRONG_METHOD, Message: "Error: wrong http method"}
	}

	sendResponse(w, resp)
}
