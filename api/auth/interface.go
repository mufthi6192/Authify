package auth

import (
	responseFormatter "SMM-PPOB/helper/formatter"
)

type RepositoryInterface interface {
	LoginRepository(data LoginDto) responseFormatter.QueryData
	LogoutRepository(token string) responseFormatter.QueryData
	RegisterRepository(data UserRegisterDto) responseFormatter.QueryData
	ResetPasswordRepository(data ResetPasswordDto) responseFormatter.QueryData
	SendResetPasswordRepository(email string, code string) responseFormatter.QueryData
	InsertEmailVerificationRepository(userId string, code string) responseFormatter.QueryData
	GetEmailVerificationRepository(userId string) responseFormatter.QueryData
}

type ServiceInterface interface {
	LoginService() responseFormatter.HttpData
	LogoutService() responseFormatter.HttpData
	RegisterService() responseFormatter.HttpData
	ResetPasswordService() responseFormatter.HttpData
	SendResetPasswordService() responseFormatter.HttpData
	//InsertEmailVerificationService(db *gorm.DB, unique string) responseFormatter.HttpData
}
