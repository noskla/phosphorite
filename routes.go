package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func RouteGetUserList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	sortBy := c.DefaultQuery("sort_by", "name")

	page, pageErr := strconv.Atoi(pageStr)
	pageSize, sizeErr := strconv.Atoi(pageSizeStr)
	if pageErr != nil && sizeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":          false,
			"status_code": 8,
			"users":       []string{},
			"count":       0,
		})
		return
	}

	err, answer, users, count := GetUserList(page, pageSize, sortBy)
	if err != nil || answer != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":          false,
			"status_code": answer,
			"users":       []string{},
			"count":       0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":          true,
		"status_code": answer,
		"users":       users,
		"count":       count,
	})

}
