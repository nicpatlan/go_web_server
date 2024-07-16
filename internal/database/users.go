package database

import (
	"errors"
	"log"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     []byte    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpires time.Time `json:"expiration"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type TokenResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (db *DB) CreateUser(email, password string) (UserResponse, error) {
	userResponse := UserResponse{
		ID:    db.userID,
		Email: email,
	}
	hashedPass, e := HashPassword(password)
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

func (db *DB) ValidateUserPassword(email, password string) (TokenResponse, error) {
	tokenResponse := TokenResponse{}
	if db.userID != 1 {
		dbStruct, err := db.loadDB()
		if err != nil {
			return tokenResponse, err
		}
		user, err := db.locateUser(email, dbStruct)
		if err != nil {
			return tokenResponse, err
		}
		err = ValidatePassword(password, user.Password)
		if err != nil {
			return tokenResponse, errors.New("invalid password")
		}
		refreshToken, err := db.UpdateUserRefreshToken(user, dbStruct)
		if err != nil {
			return tokenResponse, nil
		}
		tokenResponse.ID = user.ID
		tokenResponse.Email = user.Email
		tokenResponse.RefreshToken = refreshToken
		return tokenResponse, nil
	}
	return tokenResponse, errors.New("user not found")
}

func (db *DB) UpdateUser(id int, email, password string) (UserResponse, error) {
	userResponse := UserResponse{}
	if db.userID != 1 {
		dbStruct, err := db.loadDB()
		if err != nil {
			return userResponse, nil
		}
		user := dbStruct.Users[id]
		user.Email = email
		hashedPass, err := HashPassword(password)
		if err != nil {
			return userResponse, err
		}
		user.Password = hashedPass
		dbStruct.Users[id] = user
		err = db.writeDB(dbStruct)
		if err != nil {
			return userResponse, err
		}
		userResponse.ID = id
		userResponse.Email = email
		return userResponse, nil
	}
	return userResponse, errors.New("user not found")
}

func (db *DB) UpdateUserRefreshToken(user User, dbStruct DBStructure) (string, error) {
	refreshToken, err := GenerateRandom()
	if err != nil {
		return "", err
	}
	sixtyDays := 60 * 24 * time.Hour
	expiration := time.Now().UTC().Add(sixtyDays)
	user.RefreshToken = refreshToken
	user.TokenExpires = expiration
	dbStruct.Users[user.ID] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (db *DB) ValidateUserRefreshToken(token string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, err := db.locateUserRefreshToken(token, dbStruct)
	if err != nil {
		return User{}, err
	}
	currentTime := time.Now().UTC()
	if currentTime.After(user.TokenExpires) {
		return User{}, errors.New("refresh token expired")
	}
	return user, nil
}

func (db *DB) RevokeUserRefreshToken(token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	user, err := db.locateUserRefreshToken(token, dbStruct)
	if err != nil {
		return err
	}
	user.RefreshToken = ""
	user.TokenExpires = time.Now().UTC()
	dbStruct.Users[user.ID] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
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

func (db *DB) locateUserRefreshToken(token string, dbStruct DBStructure) (User, error) {
	for id := 1; id <= len(dbStruct.Users); id++ {
		if dbStruct.Users[id].RefreshToken == token {
			return dbStruct.Users[id], nil
		}
	}
	return User{}, errors.New("invalid refresh token")
}
