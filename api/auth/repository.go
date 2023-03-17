package auth

import (
	"SMM-PPOB/helper"
	responseFormatter "SMM-PPOB/helper/formatter"
	"errors"
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

type repo struct {
	DB *gorm.DB
}

func Repository(db *gorm.DB) RepositoryInterface {
	return &repo{
		DB: db,
	}
}

func (r *repo) LoginRepository(data LoginDto) responseFormatter.QueryData {

	db := r.DB

	authCheckChan := make(chan param)
	dataChan := make(chan param)
	loginAttemptChan := make(chan param)
	defer close(authCheckChan)
	defer close(dataChan)
	defer close(loginAttemptChan)

	group := &sync.WaitGroup{}

	tx := db.Begin()

	go func() {
		group.Add(1)
		defer group.Done()

		var total int64

		err := db.Table("users").Where("username = ? AND password = ?", data.Username, data.Password).Count(&total).Error

		if err != nil {
			authCheckChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
				data:   nil,
			}
		} else {
			authCheckChan <- param{
				status: true,
				error:  nil,
				data:   total,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		var uid string

		err := db.Table("users").Where("username = ? AND password = ?", data.Username, data.Password).Limit(1).Pluck("id", &uid).Error

		if err != nil {
			dataChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal ! Username atau password yang anda masukan tidak ditemukan")),
				data:   nil,
			}
		} else {
			dataChan <- param{
				status: true,
				error:  nil,
				data:   uid,
			}
		}
	}()
	go func() {
		group.Add(1)
		defer group.Done()

		err := tx.Create(&LoginAttempt{
			Username:  data.Username,
			Password:  data.Password,
			IpAddress: data.Ipaddress,
			Device:    data.UserDevice,
		}).Error

		if err != nil {
			loginAttemptChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf("Gagal melakukan pengecekan data ! Silahkan coba lagi atau hubungi admin jika anda mengalami masalah secara terus menerus")),
				data:   nil,
			}
		} else {
			loginAttemptChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()

	group.Wait()

	authCheck := <-authCheckChan
	dataId := <-dataChan
	loginAttempt := <-loginAttemptChan
	userId := dataId.data.(string)

	if authCheck.status == false {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, authCheck.error.Error(), nil)
	} else if dataId.status == false {
		return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
	} else if loginAttempt.status == false {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, loginAttempt.error.Error(), nil)
	} else if authCheck.data.(int64) < 1 {
		tx.Commit()
		return responseFormatter.QueryResponse(404, false, "Gagal ! Username atau password yang anda masukan tidak ditemukan", nil)
	} else {

		tx.Commit()

		loginHistoryChan := make(chan param)
		getUserChan := make(chan param)
		defer close(loginHistoryChan)
		defer close(getUserChan)

		newGroup := &sync.WaitGroup{}

		newTx := db.Begin()

		go func() {
			newGroup.Add(1)
			defer newGroup.Done()
			var result UserDataDto

			err := db.Table("users").
				Joins("INNER JOIN user_details on users.id = user_details.user_id").
				Joins("INNER JOIN user_levels on users.id = user_levels.user_id").
				Joins("INNER JOIN user_verifications on users.id = user_verifications.user_id").
				Where("users.id = ?", userId).
				Where("users.username = ?", data.Username).
				Select("users.id as user_id, users.username as username, users.email as email," +
					"user_details.name as name,user_details.phone as phone,user_levels.level as level," +
					"user_verifications.is_verified_email as is_verified_email," +
					"user_verifications.is_verified_phone as is_verified_phone").
				Take(&result).Error

			if err != nil {
				getUserChan <- param{
					status: false,
					error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
					data:   nil,
				}
			} else {
				getUserChan <- param{
					status: true,
					error:  nil,
					data:   result,
				}
			}
		}()
		go func() {
			newGroup.Add(1)
			defer newGroup.Done()

			err := newTx.Create(&LoginHistory{
				UserId:    userId,
				IpAddress: data.Ipaddress,
				Device:    data.UserDevice,
				CreatedAt: carbon.Now().SetLocale("id"),
			}).Error

			if err != nil {
				loginHistoryChan <- param{
					status: false,
					error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
					data:   nil,
				}
			} else {
				loginHistoryChan <- param{
					status: true,
					error:  nil,
					data:   nil,
				}
			}
		}()

		newGroup.Wait()

		loginHistory := <-loginHistoryChan
		getUser := <-getUserChan

		if loginHistory.status != true {
			newTx.Rollback()
			return responseFormatter.QueryResponse(500, false, loginHistory.error.Error(), nil)
		} else if getUser.status != true {
			newTx.Rollback()
			return responseFormatter.QueryResponse(500, false, getUser.error.Error(), nil)
		} else {

			result := getUser.data.(UserDataDto)

			newTx.Commit()
			return responseFormatter.QueryResponse(200, true, "Successfully get user", result)
		}
	}
}

func (r *repo) LogoutRepository(token string) responseFormatter.QueryData {

	db := r.DB

	logoutChan := make(chan param)
	countChan := make(chan param)
	defer close(logoutChan)
	defer close(countChan)

	group := &sync.WaitGroup{}

	tx := db.Begin()

	go func() {
		group.Add(1)
		defer group.Done()

		var total int64

		err := db.Table("blacklist_tokens").Where("token = ?", token).Count(&total).Error

		if err != nil {
			countChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
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

		err := tx.Create(&BlacklistToken{
			Token:     token,
			CreatedAt: carbon.Now().SetLocale("id"),
		}).Error

		if err != nil {
			logoutChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
				data:   nil,
			}
		} else {
			logoutChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()

	group.Wait()

	count := <-countChan
	logout := <-logoutChan

	if count.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, count.error.Error(), nil)
	}

	if logout.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, logout.error.Error(), nil)
	}

	if count.data.(int64) >= 1 {
		tx.Rollback()
		return responseFormatter.QueryResponse(419, false, "Gagal ! Anda sudah melakukan logout", nil)
	}

	tx.Commit()
	return responseFormatter.QueryResponse(200, true, "Berhasil melakukan logout. Mohon tunggu anda akan dialihkan", nil)

}

func (r *repo) RegisterRepository(data UserRegisterDto) responseFormatter.QueryData {

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

		err := tx.Create(&User{
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

		err := tx.Create(&UserDetail{
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

		err := tx.Create(&UserLevel{
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

		err := tx.Create(&UserVerification{
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

func (r *repo) ResetPasswordRepository(data ResetPasswordDto) responseFormatter.QueryData {

	db := r.DB

	var uid string

	err := db.Table("forget_password_verifications").Where("code = ?", data.Code).Pluck("user_id", &uid).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responseFormatter.QueryResponse(404, false, "Gagal ! Data reset password tidak ditemukan", nil)
		} else {
			return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
		}
	}

	tx := db.Begin()
	err = tx.Model(&User{}).Where("id = ?", uid).Update("password", data.NewPassword).Error

	if err != nil {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
	}

	tx.Commit()
	return responseFormatter.QueryResponse(200, true, "Successfully reset password", nil)
}

func (r *repo) SendResetPasswordRepository(email string, code string) responseFormatter.QueryData {

	db := r.DB

	var total int64

	err := db.Table("users").Where("email = ?", email).Count(&total).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responseFormatter.QueryResponse(404, false, "Gagal ! Email tidak ditemukan", nil)
		} else {
			return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
		}
	}

	if total < 1 {
		return responseFormatter.QueryResponse(404, false, "Gagal ! Email tidak ditemukan", nil)
	} else {

		var uid string

		err := db.Table("users").Where("email = ?", email).Pluck("id", &uid).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return responseFormatter.QueryResponse(404, false, "Gagal ! Email tidak ditemukan", nil)
			} else {
				return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
			}
		} else {
			tx := db.Begin()

			err := tx.Create(&ForgetPasswordVerification{
				UserId:    uid,
				Code:      code,
				CreatedAt: carbon.Now().SetLocale("id"),
				UpdatedAt: carbon.Now().SetLocale("id"),
			}).Error

			if err != nil {
				tx.Rollback()
				return responseFormatter.QueryResponse(500, false, responseFormatter.InternalServerError, nil)
			} else {
				tx.Commit()
				return responseFormatter.QueryResponse(200, true, "Successfully reset password", nil)
			}
		}
	}

}

func (r *repo) InsertEmailVerificationRepository(userId string, code string) responseFormatter.QueryData {

	db := r.DB

	countChan := make(chan param)
	insertChan := make(chan param)
	defer close(countChan)
	defer close(insertChan)

	group := &sync.WaitGroup{}

	tx := db.Begin()

	go func() {
		group.Add(1)
		defer group.Done()

		var total int64

		err := db.Table("email_verifications").Where("user_id = ?", userId).Count(&total).Error

		if err != nil {
			countChan <- param{
				status: false,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
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

		err := tx.Create(&EmailVerification{
			UserId:    userId,
			Code:      code,
			CreatedAt: carbon.Now().SetLocale("id"),
			UpdatedAt: carbon.Now().SetLocale("id"),
		}).Error

		if err != nil {
			insertChan <- param{
				status: true,
				error:  errors.New(fmt.Sprintf(responseFormatter.InternalServerError)),
				data:   nil,
			}
		} else {
			insertChan <- param{
				status: true,
				error:  nil,
				data:   nil,
			}
		}
	}()

	count := <-countChan
	insert := <-insertChan

	if count.status != true {
		return responseFormatter.QueryResponse(500, false, count.error.Error(), nil)
	} else if count.data.(int64) >= 1 {
		tx.Rollback()
		return responseFormatter.QueryResponse(409, false, "Gagal ! Kode verifikasi sudah tersedia", nil)
	} else if insert.status != true {
		tx.Rollback()
		return responseFormatter.QueryResponse(500, false, count.error.Error(), nil)
	} else {
		tx.Commit()
		return responseFormatter.QueryResponse(200, true, "Successfully insert data", nil)
	}

}

func (r *repo) GetEmailVerificationRepository(userId string) responseFormatter.QueryData {

	db := r.DB

	var code string

	err := db.Table("email_verifications").
		Where("user_id = ?", userId).
		Limit(1).
		Pluck("code", &code).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responseFormatter.QueryResponse(404, false, "Gagal ! Kode verifikasi tidak tersedia", nil)
		} else {
			return responseFormatter.QueryResponse(200, false, responseFormatter.InternalServerError, nil)
		}
	} else {
		return responseFormatter.QueryResponse(200, true, "Successfully get data", code)
	}

}
