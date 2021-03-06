package main

import (
	"github.com/gin-gonic/gin"
)

func RoutePing(c *gin.Context) {
	c.String(200, "Pong!")
}
