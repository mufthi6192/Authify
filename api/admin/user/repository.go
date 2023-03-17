package user

import (
	"SMM-PPOB/api/auth"
	"SMM-PPOB/helper"
	responseFormatter "SMM-PPOB/helper/formatter"
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"sync"
)

type repo struct {
	DB *gorm.DB
}

type param struct {
	status bool
	error  error
	data   interface{}
}

func Repository(db *gorm.DB) RepositoryInterface {
	return &repo{
		DB: db,
	}
}

func (r *repo) AddUserRepository(data AddUserDto) responseFormatter.QueryData {
	db := r.DB
	group := &sync.WaitGroup{}

	insertUserChan := make(chan param)
	insertUserDetailChan := make(chan param)
	insertUserLevelChan := make(chan param)
	insertUserVerificationChan := make(chan param)
	defer close(insertUserChan)
	defer close(insertUserLevelChan)
	defer close(insertUserDetailChan)
	defer close(insertUserVerificationChan)

	uid, err := helper.GenerateUid()
	if err != nil {
		return responseFormatter.QueryResponse(500, false, "Gagal melakukan pembuatan User ID, Silahkan coba lagi", nil)
	}

	tx := db.Begin()

	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&auth.User{
			Id:        uid,
			Email:     data.Email,
			Username:  data.Username,
			Password:  data.Password,
			CreatedAt: carbon.Now(),
			UpdatedAt: carbon.Now(),
		}).Error

		if err != nil {
			insertUserChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal menambah data user, silahkan coba lagi")),
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

		err := tx.Create(&auth.UserDetail{
			UserId: uid,
			Name:   data.Name,
			Phone:  data.Phone,
		}).Error

		if err != nil {
			insertUserDetailChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal mendaftarkan detail user, silahkan coba lagi")),
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

		err := tx.Create(&auth.UserLevel{
			UserId: uid,
			Level:  "member",
		}).Error

		if err != nil {
			insertUserLevelChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal membuat level user, silahkan coba lagi")),
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

		err := tx.Create(&auth.UserVerification{
			UserId:          uid,
			IsVerifiedEmail: false,
			IsVerifiedPhone: false,
			CreatedAt:       carbon.Now(),
			UpdatedAt:       carbon.Now(),
		}).Error

		if err != nil {
			insertUserVerificationChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal memverifikasi user, silahkan coba lagi")),
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

	user := <-insertUserChan
	userDetail := <-insertUserDetailChan
	userVerification := <-insertUserVerificationChan
	userLevel := <-insertUserLevelChan

	if user.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, user.error.Error(), nil)
	}
	if userDetail.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, userDetail.error.Error(), nil)
	}
	if userVerification.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, userVerification.error.Error(), nil)
	}
	if userLevel.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, userLevel.error.Error(), nil)
	}

	tx.Commit()
	return responseFormatter.QueryResponse(201, true, "Successfully register", nil)
}

func (r *repo) GetUserRepository() responseFormatter.QueryData {
	//TODO implement me
	panic("implement me")
}

func (r *repo) GetUserDetailRepository() responseFormatter.QueryData {
	//TODO implement me
	panic("implement me")
}

func (r *repo) UpdateUserRepository() responseFormatter.QueryData {
	//TODO implement me
	panic("implement me")
}

func (r *repo) DeleteUserRepository(userId string) responseFormatter.QueryData {

	db := r.DB

	var total int64

	err := db.Table("users").Where("id = ?", userId).Count(&total).Error

	if err != nil {
		return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
	} else if total < 1 {
		return responseFormatter.QueryResponse(404, false, "Gagal ! Data user tidak ditemukan", nil)
	} else {

		tx := db.Begin()

		group := &sync.WaitGroup{}
		deleteUserChan := make(chan param)
		deleteUserLevelChan := make(chan param)
		deleteUserDetailChan := make(chan param)
		deleteUserVerificationChan := make(chan param)
		defer close(deleteUserChan)
		defer close(deleteUserDetailChan)
		defer close(deleteUserLevelChan)
		defer close(deleteUserVerificationChan)

		go func() {
			group.Add(1)
			defer group.Done()

			err := tx.Where("id = ?", userId).Delete(&auth.User{}).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					deleteUserChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal ! Data user tidak ditemukan")),
						data:   nil,
					}
				} else {
					deleteUserChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal menghapus data user, silahkan coba lagi")),
						data:   nil,
					}
				}
			} else {
				deleteUserChan <- param{
					status: true,
					error:  nil,
					data:   nil,
				}
			}
		}()
		go func() {
			group.Add(1)
			defer group.Done()

			err := tx.Where("user_id = ?", userId).Delete(&auth.UserDetail{}).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					deleteUserDetailChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal ! Data user tidak ditemukan")),
						data:   nil,
					}
				} else {
					deleteUserDetailChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal menghapus data detail, silahkan coba lagi")),
						data:   nil,
					}
				}
			} else {
				deleteUserDetailChan <- param{
					status: true,
					error:  nil,
					data:   nil,
				}
			}
		}()
		go func() {
			group.Add(1)
			defer group.Done()

			err := tx.Where("user_id = ?", userId).Delete(&auth.UserLevel{}).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					deleteUserLevelChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal ! Data user tidak ditemukan")),
						data:   nil,
					}
				} else {
					deleteUserLevelChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal menghapus data level, silahkan coba lagi.")),
						data:   nil,
					}
				}
			} else {
				deleteUserLevelChan <- param{
					status: true,
					error:  nil,
					data:   nil,
				}
			}
		}()
		go func() {
			group.Add(1)
			defer group.Done()

			err := tx.Where("user_id = ?", userId).Delete(&auth.UserVerification{}).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					deleteUserVerificationChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal ! Data user tidak ditemukan")),
						data:   nil,
					}
				} else {
					deleteUserVerificationChan <- param{
						status: false,
						error:  errors.New(fmt.Sprintf("Gagal menghapus data verifikasi, silahkan coba lagi")),
						data:   nil,
					}
				}
			} else {
				deleteUserVerificationChan <- param{
					status: true,
					error:  nil,
					data:   nil,
				}
			}
		}()

		group.Wait()

		deleteUser := <-deleteUserChan
		deleteUserLevel := <-deleteUserLevelChan
		deleteUserVerification := <-deleteUserVerificationChan
		deleteUserDetail := <-deleteUserDetailChan

		if deleteUser.status != true {
			tx.Rollback()
			return responseFormatter.QueryResponse(500, false, deleteUser.error.Error(), nil)
		} else if deleteUserLevel.status != true {
			tx.Rollback()
			return responseFormatter.QueryResponse(500, false, deleteUserLevel.error.Error(), nil)
		} else if deleteUserVerification.status != true {
			tx.Rollback()
			return responseFormatter.QueryResponse(500, false, deleteUserVerification.error.Error(), nil)
		} else if deleteUserDetail.status != true {
			tx.Rollback()
			return responseFormatter.QueryResponse(500, false, deleteUserDetail.error.Error(), nil)
		} else {
			tx.Commit()
			return responseFormatter.QueryResponse(200, true, "Berhasil menghapus data user, mohon tunggu anda akan dialihkan", nil)
		}

	}

}
