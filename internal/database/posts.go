package database

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreatePost(body string) (Post, error) {
	post := Post{
		ID:   db.postID,
		Body: body,
	}
	var dbStruct DBStructure
	var err error
	if post.ID != 1 {
		dbStruct, err = db.loadDB()
		if err != nil {
			return post, err
		}
	} else {
		dbStruct = DBStructure{
			Posts: make(map[int]Post),
		}
	}
	dbStruct.Posts[database.postID] = post
	err = db.writeDB(dbStruct)
	if err != nil {
		return post, err
	}
	db.postID++
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
