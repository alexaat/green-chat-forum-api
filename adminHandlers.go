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

	//Get user id from url
	userIdStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/admin/users/"))
	userId, err := strconv.Atoi(userIdStr)

	if err != nil || userId < 1 {
		resp.Error = &types.Error{Type: util.ERROR_PARSING_DATA, Message: fmt.Sprintf("Cannot parse user id: %v. %v", userIdStr, err)}
		sendResponse(w, resp)
	}

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

	if r.Method == "DELETE" {
		fmt.Println("Delete user")

		//Delete user

		//Delete posts

		//Delete comments

		//Delete chat messages
	}

	sendResponse(w, resp)
}
