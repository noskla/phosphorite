package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouteGetUserByID(c *gin.Context) {
	userID := c.Param("id")
	_, code, user := GetUserByID(userID)
	if code != 1 {

		var httpCode int
		if code == 3 {
			httpCode = http.StatusInternalServerError
		} else {
			httpCode = http.StatusBadRequest
		}

		c.JSON(httpCode, gin.H{
			"ok":          false,
			"status_code": code,
		})
		return

	} else {

		c.JSON(200, gin.H{
			"ok":          true,
			"status_code": code,
			"user":        &user,
		})

	}
}
