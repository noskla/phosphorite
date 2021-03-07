package main

import (
	"log"
)

func main() {
	r := CreateHTTPEngine()
	database := DatabaseConnect()
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalln("Error closing database: ", err)
		}
	}()
	DatabaseInitSchema(database)

	log.Fatalln(r.Run(":" + GetEnvVariable("PHO_PORT", "8321")))
}
