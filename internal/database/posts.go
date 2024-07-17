package database

import "errors"

type Post struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreatePost(userID int, body string) (Post, error) {
	var postID int
	if len(db.unusedIDs) > 0 {
		postID = db.unusedIDs[0]
		db.unusedIDs = db.unusedIDs[1:]
	} else {
		postID = db.postID
		db.postID++
	}
	post := Post{
		ID:       postID,
		AuthorID: userID,
		Body:     body,
	}
	var dbStruct DBStructure
	var err error
	if post.ID != 1 || db.userID != 1 {
		dbStruct, err = db.loadDB()
		if err != nil {
			return post, err
		}
	} else {
		dbStruct = DBStructure{}
	}
	if dbStruct.Posts == nil {
		dbStruct.Posts = make(map[int]Post)
	}
	dbStruct.Posts[postID] = post
	err = db.writeDB(dbStruct)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (db *DB) GetPosts(id int, authorID, asc bool) ([]Post, error) {
	var posts []Post
	dbStruct, err := db.loadDB()
	if err != nil {
		return posts, err
	}
	var key int
	var condition bool
	if asc {
		key = 1
		condition = key <= len(dbStruct.Posts)
	} else {
		key = len(dbStruct.Posts)
		condition = key > 0
	}
	for condition {
		if !authorID {
			posts = append(posts, dbStruct.Posts[key])
			if id == key {
				break
			}
		} else if dbStruct.Posts[key].AuthorID == id {
			posts = append(posts, dbStruct.Posts[key])
		}
		if asc {
			key++
			condition = key <= len(dbStruct.Posts)
		} else {
			key--
			condition = key > 0
		}
	}
	return posts, nil
}

func (db *DB) DeletePost(postID, authorID int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	author := dbStruct.Posts[postID].AuthorID
	if author == authorID {
		delete(dbStruct.Posts, postID)
		db.unusedIDs = append(db.unusedIDs, postID)
	} else {
		return errors.New("invalid authorID")
	}
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}
