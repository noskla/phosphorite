package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"phosphorite/models"
	"testing"
)

var baseUserUUID uuid.UUID
var httpEngine *gin.Engine

func TestDatabaseInitialization(t *testing.T) {
	Database = DatabaseConnect()
	DatabaseInitSchema(Database)
}

func TestUserCreation(t *testing.T) {
	var answer int
	var err error
	err, answer, baseUserUUID = CreateUser("TestUser", "12345678", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}

func TestUserSelection(t *testing.T) {
	userUUID := baseUserUUID.String()
	fmt.Println("User ID: ", userUUID)
	err, answer, user := GetUserByID(userUUID)
	assert.NoError(t, err)
	if assert.Equal(t, 1, answer) {
		assert.Equal(t, userUUID, user.ID.String())
		assert.Equal(t, "TestUser", user.Name)
		assert.NotEqual(t, "12345678", user.Password)
	}
}

func TestUserFailNoPassword(t *testing.T) {
	err, answer, _ := CreateUser("TestUser2", "", "en", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, 6, answer)
}

func TestUserIPv6(t *testing.T) {
	err, answer, _ := CreateUser("TestUser3", "12345678", "en", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	assert.NoError(t, err)
	assert.Equal(t, 1, answer)
}

func TestUserPasswordValidation(t *testing.T) {
	err, answer, userUUID := ValidateUserPassword("TestUser", "12345678", true)
	assert.NoError(t, err)
	if assert.Equal(t, 1, answer) {
		assert.Equal(t, userUUID, baseUserUUID)
	}
}

func TestUserGetByIDRoute(t *testing.T) {
	httpEngine = CreateHTTPEngine()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/user/"+baseUserUUID.String(), nil)
	if !assert.NoError(t, err) {
		return
	}
	httpEngine.ServeHTTP(w, req)

	response := make(map[string]interface{})
	log.Println(req.Body)
	b, err := ioutil.ReadAll(req.Body)
	if !assert.NoError(t, err) {
		return
	}
	defer req.Body.Close()
	assert.NoError(t, json.Unmarshal(b, &response))

	assert.Equal(t, 200, req)
	assert.Equal(t, true, response["ok"].(bool))

	userInfo := response["user"].(*models.User)
	assert.Equal(t, userInfo.ID.String(), baseUserUUID.String())
	log.Println(userInfo.Name)
}
