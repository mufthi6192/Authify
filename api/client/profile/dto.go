package profile

import "github.com/golang-module/carbon/v2"

type GetProfileDto struct {
	UserId          string `json:"user_id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Level           string `json:"level"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	IsVerifiedPhone bool   `json:"is_verified_phone"`
}

type ResendEmailVerificationDto struct {
	Name  string
	Email string
	Code  string
}

type LoginHistoriesDto struct {
	IpAddress string `json:"ip_address"`
	Device    string `json:"device"`
	CreatedAt string `json:"created_at"`
}

func (lH LoginHistoriesDto) Formatter() LoginHistoriesDto {

	createdAt := carbon.Parse(lH.CreatedAt).SetLocale("id").DiffForHumans()

	return LoginHistoriesDto{
		IpAddress: lH.IpAddress,
		Device:    lH.Device,
		CreatedAt: createdAt,
	}

}

type ChangePasswordDto struct {
	CurrentPassword         string `json:"current_password" validate:"required,char=6"`
	NewPassword             string `json:"new_password" validate:"required,char=6,nefield=CurrentPassword"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,char=6,eqfield=NewPassword"`
	UserId                  string
}
