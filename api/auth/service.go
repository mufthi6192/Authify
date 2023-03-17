package auth

import (
	queue "SMM-PPOB/app/queue/email"
	"SMM-PPOB/app/validation"
	"SMM-PPOB/helper"
	responseFormatter "SMM-PPOB/helper/formatter"
	"SMM-PPOB/package/mysql"
	password2 "SMM-PPOB/package/password"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"time"
)

type service struct {
	CTX echo.Context
}

func Service(ctx echo.Context) ServiceInterface {
	return &service{
		CTX: ctx,
	}
}

func (s *service) LoginService() responseFormatter.HttpData {

	ctx := s.CTX
	validate := validator.New()

	ipAddress := ctx.RealIP()
	userDevice := ctx.Request().UserAgent()
	username := ctx.FormValue("username")
	password, err := password2.Generate(ctx.FormValue("password"))

	if err != nil {
		return responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)
	}

	dataLogin := LoginDto{
		Username:   username,
		Password:   password,
		Ipaddress:  ipAddress,
		UserDevice: userDevice,
	}

	errValidate := validate.Struct(dataLogin)

	if errValidate != nil {
		errMsg := make(map[string]interface{})
		for _, e := range errValidate.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Kolom %s tidak boleh kosong", e.Field())
			}
		}
		return responseFormatter.HttpResponse(400, errMsg, nil)
	}

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	login := Repository(db).LoginRepository(dataLogin)

	if login.Status != true {
		return responseFormatter.HttpResponse(login.Code, login.Message, nil)
	} else {

		result := login.Data.(UserDataDto)

		claims := jwt.MapClaims{
			"user_id":           result.UserId,
			"username":          result.Username,
			"email":             result.Email,
			"name":              result.Name,
			"phone":             result.Phone,
			"level":             result.Level,
			"is_verified_email": result.IsVerifiedEmail,
			"is_verified_phone": result.IsVerifiedPhone,
			"random":            time.Now(),
			"exp":               time.Now().Add(time.Hour * 3).Unix(),
		}

		toJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		token, err := toJwt.SignedString([]byte("ice_dolce_latte"))

		if err != nil {
			log.Print(err.Error())
			return responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)
		}

		data := make(map[string]interface{})
		data["token"] = token

		return responseFormatter.HttpResponse(200, "Berhasil masuk ! Mohon tunggu anda akan dialihkan", data)
	}
}

func (s *service) LogoutService() responseFormatter.HttpData {

	t := s.CTX.Get("user").(*jwt2.Token)
	token := t.Raw

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	logout := Repository(db).LogoutRepository(token)

	if logout.Status != true {
		return responseFormatter.HttpResponse(logout.Code, logout.Message, nil)
	}

	return responseFormatter.HttpResponse(200, logout.Message, nil)

}

func (s *service) RegisterService() responseFormatter.HttpData {

	ctx := s.CTX
	validate := validator.New()

	failResponse := responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)

	err := validate.RegisterValidation("char", validation.MinimumChar)
	if err != nil {
		return failResponse
	}
	err = validate.RegisterValidation("unique_email", validation.UniqueWithDatabaseParam)
	if err != nil {
		return failResponse
	}
	err = validate.RegisterValidation("unique_username", validation.UniqueWithDatabaseParam)
	if err != nil {
		return failResponse
	}
	err = validate.RegisterValidation("unique_phone", validation.UniqueWithDatabaseParam)
	if err != nil {
		return failResponse
	}

	errValidate := validate.Struct(UserRegisterDto{
		Name:     ctx.FormValue("name"),
		Username: ctx.FormValue("username"),
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
		Phone:    ctx.FormValue("phone"),
	})

	if errValidate != nil {
		errMsg := make(map[string]interface{})

		for _, e := range errValidate.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Kolom %s wajib diisi", e.Field())
			}
			if e.Tag() == "email" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Email harus dalam format yang benar")
			}
			if e.Tag() == "char" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Kolom %s minimal harus %s karakter", e.Field(), e.Param())
			}
			if e.Tag() == "unique_email" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Email sudah terdaftar sebelumnya")
			}
			if e.Tag() == "unique_username" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Username sudah terdaftar sebelumnya")
			}
			if e.Tag() == "unique_phone" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Nomor Telepon sudah terdaftar sebelumnya")
			}
		}

		return responseFormatter.HttpResponse(400, errMsg, nil)
	}

	pass, errPass := password2.Generate(ctx.FormValue("password"))

	if errPass != nil {
		return failResponse
	}

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	data := UserRegisterDto{
		Name:     ctx.FormValue("name"),
		Username: ctx.FormValue("username"),
		Email:    ctx.FormValue("email"),
		Password: pass,
		Phone:    ctx.FormValue("phone"),
	}

	verificationStatus := func() bool {

		tx := db.Begin()

		t := carbon.Now().String()
		code := fmt.Sprintf("%s-createdAt%s", data.Email, t)

		hashedCode := helper.Hash256String(code)

		insertEmailQueue := func() bool {
			insert := queue.InsertEmailQueue(queue.VerificationEmailData{
				FromMail:         "verification@diselesain.my.id",
				ToMail:           data.Email,
				ToName:           data.Name,
				TypeMail:         queue.VerificationEmailType,
				SubjectMail:      fmt.Sprintf("Verifikasi untuk akun %s", data.Email),
				VerificationCode: hashedCode,
			}, tx)

			return insert
		}

		insertEmailVerification := func() bool {
			insert := Repository(db).InsertEmailVerificationRepository(data.Email, hashedCode)
			return insert.Status
		}

		if insertEmailQueue() != true {
			tx.Rollback()
			return false
		} else if insertEmailVerification() != true {
			tx.Rollback()
			return false
		} else {
			tx.Commit()
			return true
		}
	}

	if verificationStatus() != true {
		return responseFormatter.HttpResponse(500, "Gagal melakukan verifikasi email, silahkan coba lagi", nil)
	}

	register := Repository(db).RegisterRepository(data)

	if register.Status != true {
		return responseFormatter.HttpResponse(register.Code, register.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil melakukan pendaftaran ! Mohon tunggu anda akan dialihkan", nil)

}

func (s *service) ResetPasswordService() responseFormatter.HttpData {

	validate := validator.New()
	ctx := s.CTX

	err := validate.RegisterValidation("char", validation.MinimumChar)
	if err != nil {
		return responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)
	}

	data := ResetPasswordDto{
		Code:                    ctx.QueryParam("code"),
		NewPassword:             ctx.FormValue("new_password"),
		NewPasswordConfirmation: ctx.FormValue("new_password_confirmation"),
	}

	errValidate := validate.Struct(data)

	if errValidate != nil {
		errMsg := make(map[string]interface{})

		for _, e := range errValidate.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! %s tidak boleh kosong", e.Field())
			}
			if e.Tag() == "eqfield" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Password tidak sama")
			}
			if e.Tag() == "char" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Password harus memiliki minimal 6 karakter tanpa spasi")
			}
		}

		return responseFormatter.HttpResponse(400, errMsg, nil)
	}

	pass, err := password2.Generate(ctx.FormValue("new_password_confirmation"))

	if err != nil {
		return responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)
	}

	newData := ResetPasswordDto{
		Code:                    ctx.QueryParam("code"),
		NewPassword:             pass,
		NewPasswordConfirmation: pass,
	}

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	resetPassword := Repository(db).ResetPasswordRepository(newData)

	if resetPassword.Status != true {
		return responseFormatter.HttpResponse(resetPassword.Code, resetPassword.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil mereset password, mohon tunggu anda akan dialihkan", nil)

}

func (s *service) SendResetPasswordService() responseFormatter.HttpData {

	mail := s.CTX.FormValue("email")

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	t := carbon.Now().String()
	code := fmt.Sprintf("%s-createdAt%s", mail, t)

	hashedCode := helper.Hash256String(code)

	insertData := Repository(db).SendResetPasswordRepository(mail, hashedCode)

	if insertData.Status != true {
		return responseFormatter.HttpResponse(insertData.Code, insertData.Message, nil)
	}

	sendMail := queue.InsertEmailQueue(queue.VerificationEmailData{
		FromMail:         "verification@diselesain.my.id",
		ToMail:           mail,
		ToName:           mail,
		SubjectMail:      fmt.Sprintf("Permintaan Reset Password untuk akun %s", mail),
		TypeMail:         queue.ForgetPasswordEmailType,
		VerificationCode: hashedCode,
	}, db)

	if sendMail != true {
		return responseFormatter.HttpResponse(500, "Gagal melakukan permintaan reset password, silahkan coba lagi", nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil melakukan permintaan reset password. Silahkan cek email anda untuk melakukan reset password", nil)

}

//func (s *service) InsertEmailVerificationService(db *gorm.DB, unique string) responseFormatter.HttpData {
//
//	t := carbon.Now().String()
//	code := fmt.Sprintf("%s-createdAt%s", unique, t)
//
//	hashedCode := helper.Hash256String(code)
//
//	insert := Repository(db).InsertEmailVerificationRepository(unique, hashedCode)
//
//	if insert.Status != true {
//		return responseFormatter.HttpResponse(insert.Code, insert.Message, nil)
//	}
//
//	return responseFormatter.HttpResponse(201, "Successfully insert email verification data", nil)
//}
