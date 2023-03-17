package user

import (
	"SMM-PPOB/app/validation"
	responseFormatter "SMM-PPOB/helper/formatter"
	"SMM-PPOB/package/mysql"
	password2 "SMM-PPOB/package/password"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func (s *service) AddUserService() responseFormatter.HttpData {
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

	errValidate := validate.Struct(AddUserDto{
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

	data := AddUserDto{
		Name:     ctx.FormValue("name"),
		Username: ctx.FormValue("username"),
		Email:    ctx.FormValue("email"),
		Password: pass,
		Phone:    ctx.FormValue("phone"),
	}

	addUser := Repository(db).AddUserRepository(data)

	if addUser.Status != true {
		return responseFormatter.HttpResponse(addUser.Code, addUser.Message, nil)
	}

	return responseFormatter.HttpResponse(200, "Berhasil melakukan pendaftaran ! Mohon tunggu anda akan dialihkan", nil)
}

func (s *service) GetUserService() responseFormatter.HttpData {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetUserDetailService() responseFormatter.HttpData {
	//TODO implement me
	panic("implement me")
}

func (s *service) UpdateUserService() responseFormatter.HttpData {

	panic("err")

}

func (s *service) DeleteUserService() responseFormatter.HttpData {

	ctx := s.CTX
	userId := ctx.Param("userId")

	db := mysql.Connect()
	newDb, _ := db.DB()
	defer func(newDb *sql.DB) {
		err := newDb.Close()
		if err != nil {
			panic("Failed to close database")
		}
	}(newDb)

	deleteUser := Repository(db).DeleteUserRepository(userId)

	if deleteUser.Status != true {
		return responseFormatter.HttpResponse(deleteUser.Code, deleteUser.Message, nil)
	}

	return responseFormatter.HttpResponse(200, deleteUser.Message, nil)

}
