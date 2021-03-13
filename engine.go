package main

import (
	"github.com/gin-gonic/gin"
)

func CreateHTTPEngine() *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Generic routes
	r.GET("/ping", RoutePing)

	// Route groups
	// -- User
	users := r.Group("/user")
	{
		users.GET("/:id", RouteGetUserByID)
		users.DELETE("/:id", RouteDeleteUser)
	}
	r.GET("/users", RouteGetUserList)

	// -- Album
	r.Group("/album")
	{

	}

	// -- Song
	r.Group("/song")
	{

	}

	// --

	return r

}
