package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var database *pg.DB
var baseUserUUID uuid.UUID

func TestDatabaseInitialization(t *testing.T) {
	database = DatabaseConnect()
	DatabaseInitSchema(database)
}

func TestUserCreation(t *testing.T) {
	var answer int
	var err error
	err, answer, baseUserUUID = CreateUser(database, "TestUser", "12345678", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}

func TestUserSelection(t *testing.T) {
	userUUID := baseUserUUID.String()
	fmt.Println("User ID: ", userUUID)
	err, answer, user := GetUserByID(database, userUUID)
	assert.NoError(t, err)
	if assert.Equal(t, 1, answer) {
		assert.Equal(t, userUUID, user.ID.String())
		assert.Equal(t, "TestUser", user.Name)
		assert.NotEqual(t, "12345678", user.Password)
	}
}

func TestUserFailNoPassword(t *testing.T) {
	err, answer, _ := CreateUser(database, "TestUser2", "", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 6, answer)
}

func TestUserIPv6(t *testing.T) {
	err, answer, _ := CreateUser(database, "TestUser3", "12345678", "en", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}

func TestUserPasswordValidation(t *testing.T) {
	err, answer, userUUID := ValidateUserPassword(database, "TestUser", "12345678", true)
	assert.NoError(t, err)
	if assert.Equal(t, 1, answer) {
		assert.Equal(t, userUUID, baseUserUUID)
	}
}
