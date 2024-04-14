package database

import (
	util "green-chat-forum-api/util"
	types "green-chat-forum-api/types"
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