package seeder

import (
	"SMM-PPOB/api/auth"
	"SMM-PPOB/helper"
	password2 "SMM-PPOB/package/password"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"sync"
)

type param struct {
	status bool
	error  error
	data   interface{}
}

func UserSeeder(db *gorm.DB) {

	countChan := make(chan param)
	insertUserChan := make(chan param)
	insertUserDetailChan := make(chan param)
	insertUserLevelChan := make(chan param)
	insertUserVerificationChan := make(chan param)
	defer close(countChan)
	defer close(insertUserChan)
	defer close(insertUserLevelChan)
	defer close(insertUserDetailChan)
	defer close(insertUserVerificationChan)

	uid, err := helper.GenerateUid()
	if err != nil {
		panic("Failed to generate UUID")
	}

	password, err := password2.Generate("testing")
	if err != nil {
		panic("Failed to generate Password")
	}

	dataUser := auth.User{
		Id:        uid,
		Email:     "testing@gmail.com",
		Username:  "testing",
		Password:  password,
		CreatedAt: carbon.Now().SetLocale("id"),
		UpdatedAt: carbon.Now().SetLocale("id"),
	}

	dataUserDetails := auth.UserDetail{
		UserId: uid,
		Name:   "Testing",
		Phone:  "085267890987",
	}

	dataUserLevel := auth.UserLevel{
		UserId: uid,
		Level:  "admin",
	}

	dataInsertVerification := auth.UserVerification{
		UserId:          uid,
		IsVerifiedEmail: false,
		IsVerifiedPhone: false,
		CreatedAt:       carbon.Now().SetLocale("id"),
		UpdatedAt:       carbon.Now().SetLocale("id"),
	}

	group := &sync.WaitGroup{}
	tx := db.Begin()

	go func() {
		group.Add(1)
		defer group.Done()

		var total int64

		err := db.Table("users").Where("id = ?", uid).Count(&total).Error

		if err != nil {
			countChan <- param{
				status: false,
				error:  err,
				data:   nil,
			}
		} else {
			countChan <- param{
				status: true,
				error:  nil,
				data:   total,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&dataUser).Error

		if err != nil {
			insertUserChan <- param{
				status: false,
				error:  err,
				data:   nil,
			}
		} else {
			insertUserChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&dataUserDetails).Error

		if err != nil {
			insertUserDetailChan <- param{
				status: false,
				error:  err,
				data:   nil,
			}
		} else {
			insertUserDetailChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&dataUserLevel).Error

		if err != nil {
			insertUserLevelChan <- param{
				status: false,
				error:  err,
				data:   nil,
			}
		} else {
			insertUserLevelChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&dataInsertVerification).Error

		if err != nil {
			insertUserVerificationChan <- param{
				status: false,
				error:  err,
				data:   nil,
			}
		} else {
			insertUserVerificationChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()

	group.Wait()

	count := <-countChan
	user := <-insertUserChan
	userDetail := <-insertUserDetailChan
	userVerification := <-insertUserVerificationChan
	userLevel := <-insertUserLevelChan

	if count.status != true {
		tx.Rollback()
		panic(count.error)
	} else if user.status != true {
		tx.Rollback()
		panic(user.error)
	} else if userDetail.status != true {
		tx.Rollback()
		panic(userDetail.error)
	} else if userVerification.status != true {
		tx.Rollback()
		panic(userVerification.error)
	} else if userLevel.status != true {
		tx.Rollback()
		panic(userLevel.error)
	} else if count.data.(int64) > 0 {
		tx.Rollback()
		panic("Failed ! Data exist")
	} else {
		tx.Commit()
		fmt.Println("Successfully seed user")
	}
}
