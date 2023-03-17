package mysql

import (
	"SMM-PPOB/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	var username string
	var password string
	var host string
	var port string
	var database string

	switch config.Testing {
	case true:
		username = config.MysqlTestingUsername
		password = config.MysqlTestingPassword
		host = config.MysqlTestingHost
		port = config.MysqlTestingPort
		database = config.MysqlTestingDatabase
	case false:
		username = config.MysqlUsername
		password = config.MysqlPassword
		host = config.MysqlHost
		port = config.MysqlPort
		database = config.MysqlDatabase
	default:
		panic("Failed to read config file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	return db
}
