package util

import (
	"fmt"
	types "green-chat-forum-api/types"
	"regexp"
	"strconv"
	"time"
)

func Contains(arr []*types.User, user types.User) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i].Id == user.Id {
			return true
		}
	}
	return false
}

func SetOnLineStatus(user *types.User, clients map[int]*types.Client) {
	if _, ok := clients[user.Id]; ok {
		user.OnLine = true
	} else {
		user.OnLine = false
	}
}

func GetCurrentMilli() int64 {
	return time.Now().UnixNano() / 1000000
}
func FormatMilli(date int) string {
	t := time.Unix(int64(date)/1000, 0)
	return t.Format("02-Jan-2006 15:04:05")
}

func ErrorHandler(err error) {
	fmt.Println("Error: ", err)
}
func ValidateInput(user *types.User, age_str string) *types.Error {
	if len(user.FirstName) < 2 || len(user.FirstName) > 50 {
		return &types.Error{Type: INVALID_FIRST_NAME, Message: "Error: first name should be between 2 and 50 characters long"}
	}
	if len(user.LastName) < 2 || len(user.LastName) > 50 {
		return &types.Error{Type: INVALID_LAST_NAME, Message: "Error: last name should be between 2 and 50 characters long"}
	}

	i, err := strconv.Atoi(age_str)
	if err != nil {
		return &types.Error{Type: INVALID_AGE, Message: "Error: invalid age"}
	}
	user.Age = i

	if user.Age < 0 || user.Age > 200 {
		return &types.Error{Type: INVALID_AGE, Message: "Error: invalid age"}
	}

	//Validate gender
	if !(user.Gender == "Male" || user.Gender == "Female" || user.Gender == "Other" || user.Gender == "Prefer Not To Say") {
		return &types.Error{Type: INVALID_GENDER, Message: "Error: invalid gender option"}
	}

	if len(user.NickName) < 2 || len(user.NickName) > 50 {
		return &types.Error{Type: INVALID_NICK_NAME, Message: "Error: nick name should be between 2 and 50 characters long"}
	}

	// Check for valid email
	reg := `^[^@\s]+@[^@\s]+.[^@\s]$`
	match, err := regexp.MatchString(reg, user.Email)
	if err != nil || !match {
		return &types.Error{Type: INVALID_EMAIL, Message: "Error: invalid email"}
	}

	// Validate password
	if len(user.Password) < 6 || len(user.Password) > 50 {
		return &types.Error{Type: INVALID_PASSWORD, Message: "Error: password should be between 6 and 50 characters long"}
	}
	// Validate passwords
	if user.Password != user.Password2 {
		return &types.Error{Type: INVALID_PASSWORD_2, Message: "Error: passwords don't match"}
	}

	return nil
}

func RemoveUserInfo(user *types.User) {
	user.FirstName = ""
	user.LastName = ""
	user.Age = 0
	user.Password = ""
	user.Password2 = ""
	user.Email = ""
	user.Gender = ""
}
