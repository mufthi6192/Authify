package auth

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

type User struct {
	Id        string `gorm:"primaryKey;not null;unique"`
	Email     string `gorm:"not null;unique"`
	Username  string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	CreatedAt carbon.Carbon
	UpdatedAt carbon.Carbon
}

type UserDetail struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`
	UserId string `gorm:"not null;primaryKey;unique"`
	Name   string `gorm:"not null"`
	Phone  string
}

type UserLevel struct {
	Id     uint   `gorm:"primaryKey;autoIncrement"`
	UserId string `gorm:"not null;primaryKey;unique"`
	Level  string `gorm:"not null"`
}

type UserVerification struct {
	Id              uint   `gorm:"primaryKey;autoIncrement"`
	UserId          string `gorm:"not null;primaryKey;unique"`
	IsVerifiedEmail bool   `gorm:"default=false"`
	IsVerifiedPhone bool   `gorm:"default=false"`
	CreatedAt       carbon.Carbon
	UpdatedAt       carbon.Carbon
}

type Level struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Level     string `gorm:"not null"`
	CreatedAt carbon.Carbon
	UpdatedAt carbon.Carbon
}

type LoginHistory struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	UserId    string `gorm:"not null;primaryKey"`
	IpAddress string `gorm:"not null"`
	Device    string `gorm:"not null"`
	CreatedAt carbon.Carbon
}

type LoginAttempt struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	IpAddress string `gorm:"not null"`
	Device    string `gorm:"not null"`
}

type BlacklistToken struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Token     string `gorm:"not null"`
	CreatedAt carbon.Carbon
}

type EmailVerification struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	UserId    string `gorm:"primaryKey;not null"`
	Code      string `gorm:"not null;type=longtext"`
	CreatedAt carbon.Carbon
	UpdatedAt carbon.Carbon
}

type ForgetPasswordVerification struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	UserId    string `gorm:"primaryKey;not null"`
	Code      string `gorm:"not null;type=longtext"`
	CreatedAt carbon.Carbon
	UpdatedAt carbon.Carbon
}

func Migration(db *gorm.DB) {

	err := db.AutoMigrate(&User{}, &UserDetail{}, &UserLevel{}, UserVerification{}, Level{},
		&LoginHistory{}, &LoginAttempt{}, &BlacklistToken{},
		&EmailVerification{}, &ForgetPasswordVerification{})

	if err != nil {
		panic("Failed to migrate user")
	}

}
