package user

type AddUserDto struct {
	Name     string `validate:"required,char=1"`
	Username string `validate:"required,char=6,unique_username=users"`
	Email    string `validate:"required,email,unique_email=users"`
	Password string `validate:"required,char=8"`
	Phone    string `validate:"required,char=10,unique_phone=user_details"`
}

type GetUserDto struct {
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type GetUserDetailDto struct {
	UserId          string `json:"user_id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Level           string `json:"level"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	IsVerifiedPhone bool   `json:"is_verified_phone"`
}
