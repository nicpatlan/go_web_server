package database

import (
	"encoding/json"
	"os"
	"sync"
)

var database DB
var mu sync.RWMutex

func NewDB(path string) (*DB, error) {
	err := os.WriteFile(path, []byte(""), 0600)
	if err != nil {
		return &database, err
	}
	mu = sync.RWMutex{}
	database = DB{
		path:   path,
		mu:     &mu,
		postID: 1,
	}
	return &database, database.ensureDB()
}

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Posts map[int]Post `json:"posts"`
}

type DB struct {
	path   string
	mu     *sync.RWMutex
	postID int
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
			return Post{}, err
		}
	} else {
		dbStruct = DBStructure{
			Posts: make(map[int]Post),
		}
	}
	dbStruct.Posts[database.postID] = post
	err = db.writeDB(dbStruct)
	if err != nil {
		return Post{}, err
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

func (db *DB) ensureDB() error {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, err := os.ReadFile(db.path)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	dbStruct := DBStructure{}
	db.mu.RLock()
	err := database.ensureDB()
	if err != nil {
		return dbStruct, err
	}
	data, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	if err != nil {
		return dbStruct, err
	}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return dbStruct, err
	}
	return dbStruct, nil
}

func (db *DB) writeDB(dbStruct DBStructure) error {
	dbJSON, err := json.Marshal(dbStruct)
	if err != nil {
		return err
	}
	db.mu.Lock()
	err = os.WriteFile(db.path, []byte(dbJSON), 0600)
	if err != nil {
		return err
	}
	db.mu.Unlock()
	return nil
}
