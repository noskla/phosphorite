package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"phosphorite/models"
	"time"
)

const (
	UserNameMaxLen     = 16
	UserNameMinLen     = 4
	UserPasswordMaxLen = 64
	UserPasswordMinLen = 6
)

/*
	1 - OK           2 - Requested non-UUID value
	3 - Query error
*/
func GetUserByID(ID string) (error, int, *models.User) {

	user := new(models.User)
	if _, err := uuid.Parse(ID); err != nil {
		return err, 2, user
	}

	err := Database.Model(user).Where("id = ?", ID).ExcludeColumn("password", "avatar").Select()
	if err != nil {
		log.Println("Error requesting user: ", err)
		return err, 3, user
	}

	return nil, 1, user

}

/*
	0 - Password is not correct		1 - OK
	2 - Incorrect username length	3 - Query error
*/
func ValidateUserPassword(name string, password string, saveDate bool) (error, int, uuid.UUID) {

	if nameLength := len(name); nameLength > UserNameMaxLen || nameLength < UserNameMinLen {
		return nil, 2, uuid.UUID{}
	}
	if passLength := len(password); passLength > UserPasswordMaxLen || passLength < UserPasswordMinLen {
		return nil, 0, uuid.UUID{}
	}

	var userUUID uuid.UUID
	var hashedPassword string

	err := Database.Model((*models.User)(nil)).
		Column("id", "password").
		Where("name = ?", name).
		Select(&userUUID, &hashedPassword)
	if err != nil {
		log.Println("Query error during password validation: ", err)
		return nil, 3, uuid.UUID{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, 0, uuid.UUID{}
	}

	if saveDate {
		if _, err := Database.Model(&models.User{LastLoginDate: time.Now()}).
			Set("last_login_date = ?last_login_date").Where("id = ?", userUUID).Update(); err != nil {
			log.Println("Query error during last login date update on validation: ", err)
			return nil, 3, userUUID
		}
	}

	return nil, 1, userUUID

}

/*
	0 - Hash error           1 - OK
	2 - Username too long    3 - UUID error
    4 - Username too short   5 - Password too long
    6 - Password too short	 7 - Username in use
*/
func CreateUser(name string, password string, language string, IP string) (error, int, uuid.UUID) {

	if len(name) > UserNameMaxLen {
		return nil, 2, uuid.UUID{}
	} else if len(name) < UserNameMinLen {
		return nil, 4, uuid.UUID{}
	}

	if len(password) > UserPasswordMaxLen {
		return nil, 5, uuid.UUID{}
	} else if len(password) < UserPasswordMinLen {
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

	_, err = Database.Model(&models.User{
		ID:           userUUID,
		Name:         name,
		Password:     string(hashedPassword),
		RegisterIP:   IP,
		LastLoginIP:  IP,
		LanguageCode: language,
	}).Insert()

	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return nil, 7, uuid.UUID{}
		}
		log.Println("Error during user creation:", err)
		return err, 0, uuid.UUID{}
	}
	return nil, 1, userUUID

}

/*
	0 - Query error			1 - Ok
*/
func GetUserList(page int, pageSize int, sortBy string) (error, int, []models.User, int) {

	if !SliceContains([]string{"name", "register_date", "last_login_date", "rank_level_id"}, sortBy) {
		sortBy = "name"
	}
	if page < 1 || page > MaxInt32 {
		page = 1
	}

	var users []models.User
	offset := (pageSize * page) - pageSize
	count, err := Database.Model(&users).Order(sortBy).Limit(pageSize).Offset(offset).
		ExcludeColumn("password", "avatar").SelectAndCount()
	if err != nil {
		log.Println("Error during querying user list: ", err)
		return err, 0, []models.User{}, 0
	}

	return nil, 1, users, count

}

/*
	0 - Requested non-UUID value	1 - Ok
	2 - Query error					3 - User does not exist
*/
func DeleteUser(ID string) (error, int) {

	user := new(models.User)
	if _, err := uuid.Parse(ID); err != nil {
		return err, 0
	}

	res, err := Database.Model(user).Where("id = ?", ID).Delete()
	if err != nil {
		log.Println("Query error during user deletion: ", err)
		return err, 2
	} else if res.RowsAffected() == 0 {
		return nil, 3
	}

	return nil, 1

}
