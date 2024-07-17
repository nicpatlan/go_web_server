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
		path:      path,
		mu:        &mu,
		postID:    1,
		userID:    1,
		unusedIDs: make([]int, 0),
	}
	return &database, database.ensureDB()
}

type DBStructure struct {
	Posts map[int]Post `json:"posts"`
	Users map[int]User `json:"users"`
}

type DB struct {
	path      string
	mu        *sync.RWMutex
	postID    int
	userID    int
	unusedIDs []int
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
