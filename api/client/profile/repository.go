package profile

import (
	"SMM-PPOB/api/auth"
	responseFormatter "SMM-PPOB/helper/formatter"
	"SMM-PPOB/package/password"
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

func (r *repo) GetProfileRepository(userId string) responseFormatter.QueryData {
	//TODO implement me
	panic("implement me")
}

func (r *repo) GetLoginHistoryRepository(userId string) responseFormatter.QueryData {
	db := r.DB

	var loginHistory []LoginHistoriesDto
	var formattedLoginHistory []LoginHistoriesDto

	err := db.Table("login_histories").
		Select("ip_address,device,created_at").
		Where("user_id = ?", userId).
		Scan(&loginHistory).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) == true {
		return responseFormatter.QueryResponse(404, false, "Gagal ! Tidak ada data login ditemukan", nil)
	}

	for _, d := range loginHistory {
		data := LoginHistoriesDto{
			IpAddress: d.IpAddress,
			Device:    d.Device,
			CreatedAt: carbon.Parse(d.CreatedAt).SetLocale("id").DiffForHumans(),
		}
		formattedLoginHistory = append(formattedLoginHistory, data)
	}

	return responseFormatter.QueryResponse(200, true, "Successfully get data", formattedLoginHistory)
}

func (r *repo) GetLatestLoginHistoryRepository(userId string) responseFormatter.QueryData {

	db := r.DB

	var loginHistory LoginHistoriesDto

	err := db.Table("login_histories").
		Select("ip_address,device,created_at").
		Where("user_id = ?", userId).
		Order("created_at desc").
		First(&loginHistory).
		Error

	data := loginHistory.Formatter()

	if errors.Is(err, gorm.ErrRecordNotFound) == true {
		return responseFormatter.QueryResponse(404, false, "Gagal ! Tidak ada data login ditemukan", nil)
	}

	return responseFormatter.QueryResponse(200, true, "Successfully get data", data)
}

func (r *repo) ChangePasswordRepository(data ChangePasswordDto) responseFormatter.QueryData {

	db := r.DB

	group := &sync.WaitGroup{}
	checkPasswordChan := make(chan param)
	changePasswordChan := make(chan param)

	currentPassword, err := password.Generate(data.CurrentPassword)
	if err != nil {
		return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
	}

	newPassword, err := password.Generate(data.NewPassword)
	if err != nil {
		return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
	}

	tx := db.Begin()
	go func() {
		group.Add(1)
		defer group.Done()

		var total int64

		err := db.Table("users").
			Where("id = ? AND password = ?", data.UserId, currentPassword).
			Count(&total).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			checkPasswordChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal ! Data tidak ditemukan")),
				data:   nil,
			}
		}

		if err != nil {
			checkPasswordChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
				data:   nil,
			}
		} else {
			checkPasswordChan <- param{
				status: true,
				error:  nil,
				data:   total,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := db.Table("users").
			Where("id = ? AND password = ?", data.UserId, currentPassword).
			Update("password", newPassword).
			Error

		if err != nil {
			changePasswordChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal mengubah password ! Silahkan ulangi lagi")),
				data:   nil,
			}
		} else {
			changePasswordChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()

	group.Wait()

	checkPassword := <-checkPasswordChan
	changePassword := <-changePasswordChan

	if checkPassword.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, checkPassword.error.Error(), nil)
	} else if checkPassword.data.(int64) != 1 {
		tx.Rollback()
		return responseFormatter.QueryResponse(400, false, "Gagal ! Password lama anda salah", nil)
	} else if changePassword.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, changePassword.error.Error(), nil)
	} else {
		tx.Commit()
		return responseFormatter.QueryResponse(200, true, "Successfully change password", nil)
	}

}

func (r *repo) UpdateEmailVerificationRepository(code string) responseFormatter.QueryData {

	db := r.DB

	var total int64

	err := db.Table("email_verifications").Where("code = ?", code).Count(&total).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responseFormatter.QueryResponse(404, false, "Gagal ! Data verifikasi tidak ditemukan", nil)
		} else {
			return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
		}
	} else if total < 1 {
		return responseFormatter.QueryResponse(404, false, "Gagal ! Data verifikasi tidak ditemukan", nil)
	} else {
		var email string
		err := db.Table("email_verifications").
			Where("code = ?", code).
			Pluck("user_id", &email).
			Error

		if err != nil {
			return responseFormatter.QueryResponse(500, false, "Gagal ! Data user tidak ditemukan", nil)
		} else {

			var uid string

			err := db.Table("users").Where("email = ?", email).Pluck("id", &uid).Error

			if err != nil {
				return responseFormatter.QueryResponse(404, false, "Gagal ! Id pengguna tidak ditemukan", nil)
			} else {
				deleteChan := make(chan param)
				verifyChan := make(chan param)
				defer close(deleteChan)
				defer close(verifyChan)

				group := &sync.WaitGroup{}

				tx := db.Begin()
				go func() {
					group.Add(1)
					defer group.Done()

					err := tx.Model(&auth.UserVerification{}).
						Where("user_id = ?", uid).
						Update("is_verified_email", true).
						Error

					if err != nil {
						verifyChan <- param{
							status: false,
							error:  nil,
							data:   nil,
						}
					} else {
						verifyChan <- param{
							status: true,
							error:  nil,
							data:   nil,
						}
					}

				}()
				go func() {
					group.Add(1)
					defer group.Done()

					err := tx.Where("user_id = ?", email).Delete(&auth.EmailVerification{}).Error

					if err != nil {
						deleteChan <- param{
							status: false,
							error:  nil,
							data:   nil,
						}
					} else {
						deleteChan <- param{
							status: true,
							error:  nil,
							data:   nil,
						}
					}
				}()

				group.Wait()

				deleted := <-deleteChan
				verify := <-verifyChan

				if deleted.status != true {
					tx.Rollback()
					return responseFormatter.QueryResponse(500, false, "Gagal menghapus data verifikasi, silahkan coba lagi", nil)
				} else if verify.status != true {
					tx.Rollback()
					return responseFormatter.QueryResponse(500, false, "Gagal melakukan verifikasi, silahkan coba lagi", nil)
				} else {
					tx.Commit()
					return responseFormatter.QueryResponse(200, true, "Berhasil melakukan verifikasi, mohon tunguu anda akan dialihkan", nil)
				}
			}
		}
	}

}

func (r *repo) ResendVerificationEmailRepository(data ResendEmailVerificationDto) responseFormatter.QueryData {

	db := r.DB
	newTx := db.Begin()

	err := newTx.Create(&auth.EmailVerification{
		UserId:    data.Email,
		Code:      data.Code,
		CreatedAt: carbon.Now().SetLocale("id"),
		UpdatedAt: carbon.Now().SetLocale("id"),
	}).Error

	if err != nil {
		newTx.Rollback()
		return responseFormatter.QueryResponse(500, false, err.Error(), nil)
	}

	newTx.Commit()
	return responseFormatter.QueryResponse(200, true, "Successfully add resend", nil)
}
