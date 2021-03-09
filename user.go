package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"phosphorite/models"
)

/*
	0 - Hash error           1 - OK
	2 - Username length err  3 - UUID error
*/
func CreateUser(db *pg.DB, name string, password string, language string, IP string) (error, int) {

	if len(name) > 16 {
		return nil, 2
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error during password hash:", err)
		return err, 0
	}

	userUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error during creation of user UUID", err)
		return err, 3
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
		return err, 0
	}
	return nil, 1

}
