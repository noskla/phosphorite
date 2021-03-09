package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

var database *pg.DB

func TestDatabaseInitialization(t *testing.T) {
	database = DatabaseConnect()
	DatabaseInitSchema(database)
}

func TestUserCreation(t *testing.T) {
	err, answer := CreateUser(database, "TestUser", "12345678", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}

func TestUserFailNoPassword(t *testing.T) {
	err, answer := CreateUser(database, "TestUser2", "", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 6, answer)
}

func TestUserIPv6(t *testing.T) {
	err, answer := CreateUser(database, "TestUser3", "12345678", "en", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}
