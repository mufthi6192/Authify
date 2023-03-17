package main

import (
	queue "SMM-PPOB/app/queue/email"
	"SMM-PPOB/database/migration"
	"SMM-PPOB/database/seeder"
	"SMM-PPOB/package/mysql"
	"SMM-PPOB/routes"
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	go queue.EmailQueue()
	initApp()
	//migrateAndSeed()
}

func migrateAndSeed() {

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	migration.Fresh(db)
	seeder.Run(db)
}

func initApp() {
	e := echo.New()
	err := routes.Routes(e)

	if err != nil {
		panic("Failed to init routes")
	}

	if err = e.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
