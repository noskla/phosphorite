package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"phosphorite/models"
)

func DatabaseConnect() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     GetEnvVariable("PHO_PG_USER", "postgres"),
		Password: GetEnvVariable("PHO_PG_PASSWORD", "password"),
		Addr: GetEnvVariable("PHO_PG_HOST", "127.0.0.1") + ":" +
			GetEnvVariable("PHO_PG_PORT", "5432"),
		Database: GetEnvVariable("PHO_PG_DATABASE", "phosphorite"),
	})
}

func DatabaseInitSchema(database *pg.DB) {

	modelSchema := []interface{}{
		(*models.User)(nil),
		(*models.Token)(nil),
		(*models.Notification)(nil),
		(*models.Favourite)(nil),
	}

	for _, model := range modelSchema {
		err := database.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			log.Fatalln(err)
		}
	}

}
