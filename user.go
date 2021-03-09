package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"phosphorite/models"
)

/*
	1 - OK           2 - Requested non-UUID value
	3 - Query error
*/
func GetUserByID(db *pg.DB, ID string) (error, int, *models.User) {

	user := new(models.User)
	userUUID, err := uuid.Parse(ID)
	if err != nil {
		return err, 2, user
	}

	err = db.Model(user).Where("id = ?", userUUID).Select()
	if err != nil {
		log.Println("Error requesting user: ", err)
		return err, 3, user
	}

	return nil, 1, user

}

/*
	0 - Hash error           1 - OK
	2 - Username too long    3 - UUID error
    4 - Username too short   5 - Password too long
    6 - Password too short
*/
func CreateUser(db *pg.DB, name string, password string, language string, IP string) (error, int, uuid.UUID) {

	if len(name) > 16 {
		return nil, 2, uuid.UUID{}
	} else if len(name) < 4 {
		return nil, 4, uuid.UUID{}
	}

	if len(password) > 64 {
		return nil, 5, uuid.UUID{}
	} else if len(password) < 6 {
		return nil, 6, uuid.UUID{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error during password hash:", err)
		return err, 0, uuid.UUID{}
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error during creation of user UUID", err)
		return err, 3, uuid.UUID{}
	}

	_, err = db.Model(&models.User{
		ID:           userUUID,
		Name:         name,
		Password:     string(hashedPassword),
		RegisterIP:   IP,
		LastLoginIP:  IP,
		LanguageCode: language,
	}).Insert()

	if err != nil {
		log.Println("Error during user creation:", err)
		return err, 0, uuid.UUID{}
	}
	return nil, 1, userUUID

}
