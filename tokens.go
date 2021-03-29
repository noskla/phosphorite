package main

import (
	"github.com/google/uuid"
	"phosphorite/models"
	"strconv"
	"time"
)

const (
	TokenLength = 64
)

/*
	0 - Requested non-UUID value	1 - Ok
	2 - Query error
*/
func GenerateAndAddToken(userID string) (error, int, string) {
	ID, err := uuid.Parse(userID)
	if err != nil {
		return err, 0, ""
	}

	token := RandomString(TokenLength) + strconv.Itoa(int(time.Now().UnixNano()))
	model := &models.Token{Token: token, Owner: &models.User{ID: ID}}
	res, err := Database.Model(model).Insert()
	if err != nil || res.RowsAffected() != 1 {
		return err, 2, ""
	}

	return nil, 1, token
}
