package auth

type LoginDto struct {
	Username   string `validate:"required"`
	Password   string `validate:"required"`
	Ipaddress  string
	UserDevice string
}

type UserDataDto struct {
	UserId          string `json:"user_id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Level           string `json:"level"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	IsVerifiedPhone bool   `json:"is_verified_phone"`
}

type UserRegisterDto struct {
	Name     string `validate:"required,char=1"`
	Username string `validate:"required,char=6,unique_username=users"`
	Email    string `validate:"required,email,unique_email=users"`
	Password string `validate:"required,char=8"`
	Phone    string `validate:"required,char=10,unique_phone=user_details"`
}

type ResetPasswordDto struct {
	Code                    string `validate:"required,char=6"`
	NewPassword             string `json:"new_password" validate:"required,char=6"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,char=6,eqfield=NewPassword"`
}
