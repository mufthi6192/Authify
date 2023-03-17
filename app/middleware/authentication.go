package middleware

import (
	responseFormatter "SMM-PPOB/helper/formatter"
	"SMM-PPOB/package/mysql"
	"database/sql"
	"errors"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"sync"
)

type authParam struct {
	status bool
	error  error
	data   interface{}
}

type resChan chan authParam

func jwtAuth(group *sync.WaitGroup, db *gorm.DB, token string, res resChan) {
	group.Add(1)
	defer group.Done()

	var total int64

	err := db.Table("blacklist_tokens").Where("token = ?", token).Count(&total).Error

	if err != nil {
		res <- authParam{
			status: false,
			error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
			data:   nil,
		}
	} else {
		res <- authParam{
			status: true,
			error:  nil,
			data:   total,
		}
	}
}

func checkUser(group *sync.WaitGroup, db *gorm.DB, userId string, res resChan) {
	group.Add(1)
	defer group.Done()

	var total int64

	err := db.Table("users").Where("id = ?", userId).Count(&total).Error

	if err != nil {
		res <- authParam{
			status: false,
			error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
			data:   nil,
		}
	} else {
		res <- authParam{
			status: true,
			error:  nil,
			data:   total,
		}
	}
}

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := c.Get("user").(*jwt2.Token)
		claims := token.Claims.(jwt2.MapClaims)
		userId := claims["user_id"].(string)

		db := mysql.Connect()
		newDb, _ := db.DB()
		defer func(newDb *sql.DB) {
			err := newDb.Close()
			if err != nil {
				panic("Failed to close database")
			}
		}(newDb)

		group := &sync.WaitGroup{}
		jwtAuthChan := make(chan authParam)
		checkAuthChan := make(chan authParam)
		defer close(jwtAuthChan)
		defer close(checkAuthChan)

		go jwtAuth(group, db, token.Raw, jwtAuthChan)
		go checkUser(group, db, userId, checkAuthChan)

		group.Wait()

		jwtAuthRes := <-jwtAuthChan
		checkAuthRes := <-checkAuthChan

		if jwtAuthRes.status != true || checkAuthRes.status != true {
			return c.String(500, fmt.Sprintf(responseFormatter.InternalServerError))
		}
		if jwtAuthRes.data.(int64) >= 1 || checkAuthRes.data.(int64) < 1 {
			return c.String(401, fmt.Sprintf("Failed ! User unauthorized"))
		}

		return next(c)
	}
}
