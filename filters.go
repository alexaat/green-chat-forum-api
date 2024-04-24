package main

import (
	db "green-chat-forum-api/database"
	types "green-chat-forum-api/types"
)

func filterPostsByBannedUser(posts []types.Post) (*[]types.Post, error) {
	filtered := []types.Post{}
	for _, post := range posts {
		user, err := db.GetUserById(post.UserId)
		if err != nil {
			return nil, err
		}
		if user.Status != "banned" {
			filtered = append(filtered, post)
		}
	}
	return &filtered, nil
}

func filterCommentsByUserId(comments []*types.Comment) ([]*types.Comment, error) {
	filtered := []*types.Comment{}
	for _, comment := range comments {
		userId := comment.UserId
		user, err := db.GetUserById(userId)

		if err != nil {
			return nil, err
		}

		if user.Status != "banned" {
			filtered = append(filtered, comment)
		}
	}
	return filtered, nil
}

func filterUsersByBanned(users []*types.User) []*types.User {
	filtered := []*types.User{}
	for _, user := range users {
		if user.Status != "banned" {
			filtered = append(filtered, user)
		}
	}
	return filtered
}
