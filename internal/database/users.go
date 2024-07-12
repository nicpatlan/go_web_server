package database

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

func (db *DB) locateUser(email string, dbStruct DBStructure) (User, error) {
	for idx := 1; idx <= len(dbStruct.Users); idx++ {
		if dbStruct.Users[idx].Email == email {
			return dbStruct.Users[idx], nil
		}
	}
	return User{}, errors.New("email not found")
}
