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

func (db *DB) GetPosts() ([]Post, error) {
	var posts []Post
	dbStruct, err := db.loadDB()
	if err != nil {
		return posts, err
	}
	for key := 1; key < len(dbStruct.Posts)+1; key++ {
		posts = append(posts, dbStruct.Posts[key])
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
