package database

import (
	"encoding/json"

	types "green-chat-forum-api/types"
	util "green-chat-forum-api/util"
)

//      _________posts________________________________________________
//     |  id       |  date     |  user_id  |  content  |  categories  |
//     |  INTEGER  |  INTEGER  |  INTEGER  |  TEXT     |  TEXT        |

func creratePostsTable() error {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts(id INTEGER PRIMARY KEY, date INTEGER NOT NULL, user_id INTEGER NOT NULL, content TEXT NOT NULL, categories TEXT)")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func InsertPost(user *types.User, post *types.Post) error {
	statement, err := db.Prepare("INSERT INTO posts (date, user_id, content, categories) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	date := util.GetCurrentMilli()
	categories, err := json.Marshal(post.Categories)
	if err != nil {
		categories = []byte("[]")
	}
	_, err = statement.Exec(date, user.Id, post.Content, string(categories))
	if err != nil {
		return err
	}
	return nil
}

func GetPosts(user *types.User) (*[]types.Post, error) {
	posts := []types.Post{}

	if user == nil {
		return nil, nil
	}

	sql := `
	SELECT posts.id, date, user_id, users.nick_name, content, categories
	FROM posts
	INNER JOIN users
	ON user_id = users.id
	ORDER BY date DESC`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := types.Post{}
		var categories string
		err = rows.Scan(&(post.Id), &(post.Date), &(post.UserId), &(post.NickName), &(post.Content), &categories)
		if err != nil {
			return nil, err
		}
		var arr []string
		err = json.Unmarshal([]byte(categories), &arr)

		if err == nil {
			post.Categories = arr
		} else {
			post.Categories = []string{}
		}
		numberOfComments, err := GetNumberOfComments(post.Id)
		if err != nil {
			return nil, err
		}
		if numberOfComments == -1 {
			numberOfComments = 0
		}
		post.NumberOfComments = numberOfComments
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func GetPostsByUserId(userId int) (*[]types.Post, error) {
	posts := []types.Post{}

	sql := `
	SELECT posts.id, date, user_id, users.nick_name, content, categories
	FROM posts
	INNER JOIN users
	ON user_id = users.id
	WHERE user_id = ?
	ORDER BY date DESC`
	rows, err := db.Query(sql, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := types.Post{}
		var categories string
		err = rows.Scan(&(post.Id), &(post.Date), &(post.UserId), &(post.NickName), &(post.Content), &categories)
		if err != nil {
			return nil, err
		}
		var arr []string
		err = json.Unmarshal([]byte(categories), &arr)

		if err == nil {
			post.Categories = arr
		} else {
			post.Categories = []string{}
		}
		numberOfComments, err := GetNumberOfComments(post.Id)
		if err != nil {
			return nil, err
		}
		if numberOfComments == -1 {
			numberOfComments = 0
		}
		post.NumberOfComments = numberOfComments
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func GetPost(postId int) (*types.Post, error) {
	post := types.Post{}

	sql := `
	SELECT posts.id, date, user_id, users.nick_name, content, categories
	FROM posts
	INNER JOIN users
	ON user_id = users.id
	WHERE posts.id = ?
	LIMIT 1`
	rows, err := db.Query(sql, postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var categories string
		err = rows.Scan(&(post.Id), &(post.Date), &(post.UserId), &(post.NickName), &(post.Content), &categories)
		if err != nil {
			return nil, err
		}
		var arr []string
		err = json.Unmarshal([]byte(categories), &arr)

		if err == nil {
			post.Categories = arr
		} else {
			post.Categories = []string{}
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	numberOfComments, err := GetNumberOfComments(postId)
	if err != nil {
		return nil, err
	}
	if numberOfComments == -1 {
		numberOfComments = 0
	}
	post.NumberOfComments = numberOfComments

	return &post, nil
}

func DeletePostsByUserId(id int) (*int64, error) {
	statement, err := db.Prepare("DELETE FROM posts WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Exec(id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return &rowsAffected, nil
}
