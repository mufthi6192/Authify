package profile

import (
	queue "SMM-PPOB/app/queue/email"
	"SMM-PPOB/app/validation"
	"SMM-PPOB/helper"
	responseFormatter "SMM-PPOB/helper/formatter"
	"SMM-PPOB/package/mysql"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
)

type service struct {
	CTX echo.Context
}

func Service(ctx echo.Context) ServiceInterface {
	return &service{
		CTX: ctx,
	}
}

func (s *service) GetProfileService() responseFormatter.HttpData {

	c := s.CTX

	token := c.Get("user").(*jwt2.Token)
	claims := token.Claims.(jwt2.MapClaims)
	userId := claims["user_id"].(string)
	username := claims["username"].(string)
	email := claims["username"].(string)
	name := claims["name"].(string)
	phone := claims["phone"].(string)
	level := claims["level"].(string)
	isVerifiedEmail := claims["is_verified_email"].(bool)
	isVerifiedPhone := claims["is_verified_phone"].(bool)

	data := GetProfileDto{
		UserId:          userId,
		Username:        username,
		Email:           email,
		Name:            name,
		Phone:           phone,
		Level:           level,
		IsVerifiedEmail: isVerifiedEmail,
		IsVerifiedPhone: isVerifiedPhone,
	}

	return responseFormatter.HttpResponse(200, "Successfully get profile", data)
}

func (s *service) GetLoginHistoryService() responseFormatter.HttpData {

	c := s.CTX

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

	loginHistory := Repository(db).GetLoginHistoryRepository(userId)

	if loginHistory.Status != true {
		return responseFormatter.HttpResponse(loginHistory.Code, loginHistory.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Successfully get data", loginHistory.Data)
}

func (s *service) GetLatestLoginHistoryService() responseFormatter.HttpData {
	c := s.CTX

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

	loginHistory := Repository(db).GetLatestLoginHistoryRepository(userId)

	if loginHistory.Status != true {
		return responseFormatter.HttpResponse(loginHistory.Code, loginHistory.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Successfully get data", loginHistory.Data)
}

func (s *service) ChangePasswordService() responseFormatter.HttpData {

	validate := validator.New()
	ctx := s.CTX

	err := validate.RegisterValidation("char", validation.MinimumChar)
	if err != nil {
		return responseFormatter.HttpResponse(500, responseFormatter.InternalServerError, nil)
	}

	token := ctx.Get("user").(*jwt2.Token)
	claims := token.Claims.(jwt2.MapClaims)
	userId := claims["user_id"].(string)

	data := ChangePasswordDto{
		CurrentPassword:         ctx.FormValue("current_password"),
		NewPassword:             ctx.FormValue("new_password"),
		NewPasswordConfirmation: ctx.FormValue("new_password_confirmation"),
		UserId:                  userId,
	}

	errValidate := validate.Struct(data)

	if errValidate != nil {
		errMsg := make(map[string]interface{})

		for _, e := range errValidate.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Kolom %s wajib diisi", e.Field())
			}
			if e.Tag() == "eqfield" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Password tidak sama")
			}
			if e.Tag() == "nefield" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Password lama tidak boleh sama dengan password baru")
			}
			if e.Tag() == "char" {
				errMsg[e.Field()] = fmt.Sprintf("Gagal ! Password harus memiliki minimal 6 karakter tanpa spasi")
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

	changePassword := Repository(db).ChangePasswordRepository(data)

	if changePassword.Status != true {
		return responseFormatter.HttpResponse(changePassword.Code, changePassword.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil mengganti password. Mohon tunggu anda akan dialihkan", nil)

}

func (s *service) UpdateEmailVerificationService() responseFormatter.HttpData {

	code := s.CTX.QueryParam("code")

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	update := Repository(db).UpdateEmailVerificationRepository(code)

	if update.Status != true {
		return responseFormatter.HttpResponse(update.Code, update.Message, nil)
	}

	return responseFormatter.HttpResponse(update.Code, update.Message, nil)

}

func (s *service) ResendVerificationEmailService() responseFormatter.HttpData {

	c := s.CTX

	token := c.Get("user").(*jwt2.Token)
	claims := token.Claims.(jwt2.MapClaims)
	name := claims["name"].(string)
	email := claims["email"].(string)

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	tx := db.Begin()

	t := carbon.Now().String()
	code := fmt.Sprintf("%s-createdAt%s", email, t)

	hashedCode := helper.Hash256String(code)

	mail := queue.InsertEmailQueue(queue.VerificationEmailData{
		FromMail:         "verification@diselesain.my.id",
		ToMail:           email,
		ToName:           name,
		SubjectMail:      fmt.Sprintf("Konfirmasi untuk akun %s", email),
		TypeMail:         queue.VerificationEmailType,
		VerificationCode: hashedCode,
	}, tx)

	if mail != true {
		return responseFormatter.HttpResponse(500, "fail", nil)
	}

	verification := Repository(db).ResendVerificationEmailRepository(ResendEmailVerificationDto{
		Name:  name,
		Email: email,
		Code:  hashedCode,
	})

	if verification.Status != true {
		return responseFormatter.HttpResponse(verification.Code, verification.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil ! Silahkan cek inbox atau spam email anda", nil)
}
