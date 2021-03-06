package main

import (
	"log"
	"os"
)

func main() {
	r := CreateHTTPEngine()

	port := os.Getenv("PHO_PORT")
	log.Fatalln(r.Run(":" + port))
}
