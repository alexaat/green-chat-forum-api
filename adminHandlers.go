package main

import (
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
	}

	if r.Method == "DELETE" {
		fmt.Println("Delete user")

		//Delete user

		//Delete posts

		//Delete comments

		//Delete chat messages
	}

	sendResponse(w, resp)
}

func adminSignUpHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.Response{Payload: nil, Error: nil}

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
	sendResponse(w, resp)
}