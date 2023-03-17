package validation

import (
	"SMM-PPOB/package/mysql"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ExistsDatabase(field validator.FieldLevel) bool {

	tagParam := field.GetTag()
	var fieldName string

	countSeparator := strings.Count(tagParam, "_")

	if countSeparator >= 2 {
		split := strings.SplitN(tagParam, "_", 2)
		fieldName = split[1]
	} else if countSeparator == 1 {
		split := strings.Split(tagParam, "_")
		fieldName = split[1]
	} else {
		panic("Failed ! Separator can't be empty")
	}

	value := field.Field().String()
	tableName := field.Param()

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	var total int64

	err := db.Table(tableName).Where(fieldName+"= ?", value).Count(&total).Error

	if err != nil {
		return false
	}

	if total < 1 {
		return false
	}

	return true

}
