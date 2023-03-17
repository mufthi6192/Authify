package migration

import (
	"SMM-PPOB/api/auth"
	queue "SMM-PPOB/app/queue/email"
	"fmt"
	"gorm.io/gorm"
)

func Fresh(db *gorm.DB) {

	//Drop All Table Section
	var tables []string
	err := db.Raw("SHOW TABLES").Pluck("Tables_in_smm_ppob", &tables).Error

	if err != nil {
		panic("Failed to fetch all table name")
	}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("DROP TABLE %s", table))
	}

	//Register Migration
	auth.Migration(db)
	queue.Migration(db)

	fmt.Println("Successfully Migrate Database")
}
