package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreation(t *testing.T) {
	db := DatabaseConnect()
	DatabaseInitSchema(db)

	err, answer := CreateUser(db, "TestUser", "12345678", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}
