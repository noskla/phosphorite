package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/user/"+baseUserUUID.String(), nil)
	if !assert.NoError(t, err) {
		return
	}
	httpEngine.ServeHTTP(res, req)

	response := make(map[string]interface{})
	bodyJSON, err := ioutil.ReadAll(res.Body)
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, json.Unmarshal(bodyJSON, &response))

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, true, response["ok"].(bool))

	userInfo := response["user"].(map[string]interface{})
	assert.Equal(t, baseUserUUID.String(), userInfo["id"].(string))
}

func TestUserListing(t *testing.T) {
	err, answer, userList, userCount := GetUserList(1, 10, "name")
	assert.NoError(t, err)
	if !assert.Equal(t, 1, answer) {
		return
	}
	assert.NotEqual(t, 0, userCount)
	assert.Equal(t, "TestUser", userList[0].Name)
}
