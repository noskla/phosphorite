package main

import (
	"github.com/gin-gonic/gin"
)

func RoutePing(c *gin.Context) {
	c.String(200, "Pong!")
}


// Hej czy wiesz dlaczego Tiemman w 66 ruchu poświęcił swojego gońca w partii z 1979 roku?? 
