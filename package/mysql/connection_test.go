package mysql

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {

	db := Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	err := db.Raw("SHOW DATBASES").Error

	if err != nil {
		t.Fatal("Failed to connect MySQL")
	}

	fmt.Println("Successfully Connect MySQL")
}
