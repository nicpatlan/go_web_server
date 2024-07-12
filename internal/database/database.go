package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
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
		userID: 1,
	}
	return &database, database.ensureDB()
}

type Post struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Posts map[int]Post `json:"posts"`
	Users map[int]User `json:"users"`
}

type DB struct {
	path   string
	mu     *sync.RWMutex
	postID int
	userID int
}

func (db *DB) CreateUser(email, password string) (UserResponse, error) {
	userResponse := UserResponse{
		ID:    db.userID,
		Email: email,
	}
	hashedPass, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		log.Printf("Error generating hashed password: %s", e)
		return userResponse, e
	}

	user := User{
		ID:       db.userID,
		Email:    email,
		Password: hashedPass,
	}

	var dbStruct DBStructure
	var err error
	if db.userID != 1 {
		dbStruct, err = db.loadDB()
		if err != nil {
			return userResponse, err
		}
		_, err = db.locateUser(email, dbStruct)
		if err == nil {
			return userResponse, errors.New("duplicate email")
		}
	} else {
		dbStruct = DBStructure{
			Users: make(map[int]User),
		}
	}
	dbStruct.Users[db.userID] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return userResponse, err
	}
	db.userID++
	return userResponse, nil
}

func (db *DB) ValidateUserPassword(email, password string) (UserResponse, error) {
	userResponse := UserResponse{}
	if db.userID != 1 {
		dbStruct, err := db.loadDB()
		if err != nil {
			return userResponse, err
		}
		user, err := db.locateUser(email, dbStruct)
		if err != nil {
			return userResponse, err
		}
		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		if err != nil {
			return userResponse, errors.New("invalid password")
		}
		userResponse.ID = user.ID
		userResponse.Email = user.Email
		return userResponse, nil
	}
	return userResponse, errors.New("user not found")
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

func (db *DB) locateUser(email string, dbStruct DBStructure) (User, error) {
	for idx := 1; idx <= len(dbStruct.Users); idx++ {
		if dbStruct.Users[idx].Email == email {
			return dbStruct.Users[idx], nil
		}
	}
	return User{}, errors.New("email not found")
}
