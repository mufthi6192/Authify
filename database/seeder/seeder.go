package seeder

import (
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {

	//Register Seeder here
	UserSeeder(db)

}
