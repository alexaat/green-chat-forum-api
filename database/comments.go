package database

import (
	types "green-chat-forum-api/types"
	util "green-chat-forum-api/util"
)

//       ________comments_____________________________________________
//      |  id       |  date     |  user_id   |  post_id   |  content  |
//      |  INTEGER  |  INTEGER  |  INTEGER   |  INTEGER   |  TEXT     |

// Create comments table
func crerateCommentsTable() error {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments(id INTEGER PRIMARY KEY, date INTEGER NOT NULL, user_id INTEGER NOT NULL, post_id INTEGER NOT NULL, content TEXT NOT NULL)")
	if err != nil {
		return err
	}
	defer statement.Close()
	statement.Exec()
	return nil
}

func SaveComment(userId int, postId int, comment string) error {
	statement, err := db.Prepare("INSERT INTO comments (date, user_id, post_id, content) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	date := util.GetCurrentMilli()
	_, err = statement.Exec(date, userId, postId, comment)
	if err != nil {
		return err
	}
	return nil
}

func GetComments(postId int) ([]*types.Comment, error) {
	comments := []*types.Comment{}
	sql := `
	SELECT comments.id, comments.date, comments.user_id, users.nick_name, comments.post_id, comments.content 
	FROM comments
	INNER JOIN users
	ON comments.user_id = users.id	
	WHERE post_id = ?
	ORDER BY comments.date DESC	
	`
	rows, err := db.Query(sql, postId)
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		comment := types.Comment{}
		err = rows.Scan(&(comment.Id), &(comment.Date), &(comment.UserId), &(comment.UserNickName), &(comment.PostId), &(comment.Content))
		if err != nil {
			return comments, err
		}
		comments = append(comments, &comment)
	}
	err = rows.Err()
	if err != nil {
		return comments, err
	}
	return comments, nil
}

func GetNumberOfComments(postId int) (int, error) {

	sql := `
	SELECT content
	FROM comments
	WHERE post_id = ? AND content IS NOT NULL
	`
	rows, err := db.Query(sql, postId)
	if err != nil {
		return -1, err
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if rows.Err() != nil {
		return -1, err
	}
	return counter, nil
}

func DeleteCommentsByUserId(id int) (*int64, error) {
	statement, err := db.Prepare("DELETE FROM comments WHERE user_id = ?")
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

func DeleteCommentById(id int) (*int64, error) {
	statement, err := db.Prepare("DELETE FROM comments WHERE id = ?")
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

func GetAllComments() (*[]types.Comment, error) {
	comments := []types.Comment{}
	sql := `
	SELECT comments.id, comments.date, comments.user_id, users.nick_name, comments.post_id, comments.content 
	FROM comments
	INNER JOIN users
	ON comments.user_id = users.id	
	ORDER BY comments.date DESC	
	`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		comment := types.Comment{}
		err = rows.Scan(&(comment.Id), &(comment.Date), &(comment.UserId), &(comment.UserNickName), &(comment.PostId), &(comment.Content))
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &comments, nil
}

func GetCommentById(id int) (*types.Comment, error) {
	comment := types.Comment{}
	sql := `
	SELECT comments.id, comments.date, comments.user_id, users.nick_name, comments.post_id, comments.content 
	FROM comments
	INNER JOIN users
	ON comments.user_id = users.id	
	WHERE comments.id = ?
	LIMIT 1	
	`
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&(comment.Id), &(comment.Date), &(comment.UserId), &(comment.UserNickName), &(comment.PostId), &(comment.Content))
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
