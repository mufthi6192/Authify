package queue

import (
	"SMM-PPOB/app/email"
	"SMM-PPOB/package/mysql"
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	VerificationEmailType   = "verification"
	ForgetPasswordEmailType = "forget_password"
)

type EmailQueueModel struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	FromMail    string `gorm:"not null"`
	ToMail      string `gorm:"not null"`
	ToName      string `gorm:"not null"`
	SubjectMail string `gorm:"not null"`
	BodyMail    string `gorm:"not null;type:longtext"`
	TypeMail    string `gorm:"type:enum('verification','forget_password');not null"`
	StatusMail  bool   `gorm:"default:false; not null"`
	CreatedAt   carbon.Carbon
	UpdatedAt   carbon.Carbon
}

type VerificationEmailData struct {
	FromMail         string
	ToMail           string
	ToName           string
	SubjectMail      string
	TypeMail         string
	VerificationCode string
}

func (EmailQueueModel) TableName() string {
	return "email_queues"
}

func Migration(db *gorm.DB) {
	err := db.AutoMigrate(&EmailQueueModel{})
	if err != nil {
		panic("Failed to migrate Email Queue")
	}
}

func InsertEmailQueue(data VerificationEmailData, tx *gorm.DB) bool {

	err := tx.Create(&EmailQueueModel{
		FromMail:    data.FromMail,
		ToMail:      data.ToMail,
		ToName:      data.ToName,
		SubjectMail: data.SubjectMail,
		BodyMail:    data.VerificationCode,
		TypeMail:    data.TypeMail,
		StatusMail:  false,
		CreatedAt:   carbon.Now().SetLocale("id"),
		UpdatedAt:   carbon.Now().SetLocale("id"),
	}).Error

	if err != nil {
		tx.Rollback()
		return false
	}

	tx.Commit()
	return true

}

func EmailQueue() {

	db := mysql.Connect()

	for {
		var total int64
		err := db.Table("email_queues").Where("status_mail = ?", false).Count(&total).Error

		if err != nil {
			log.Fatal("Failed to count email queue")
		} else if total < 1 {
			log.Print("There's no remain email queue")
		} else {
			sendMail(db)
		}
		time.Sleep(1 * time.Second)
		continue
	}

}

func sendMail(db *gorm.DB) {

	var data email.Mail

	err := db.Table("email_queues").
		Select("id,from_mail,to_mail,to_name,subject_mail,body_mail,type_mail").
		Where("status_mail = ?", false).
		First(&data).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print("There's no remain email queue")
		} else {
			log.Print("Failed to get data")
		}
	} else {
		if data.TypeMail == VerificationEmailType {
			err := email.SendVerificationMail(data)

			if err != nil {
				log.Print("Failed to sending mail")
			} else {
				log.Print(fmt.Sprintf("Mail delivered to %s", data.ToMail))
				UpdateEmailQueue(data.Id, db)
			}
		} else if data.TypeMail == ForgetPasswordEmailType {
			err := email.SendForgetPasswordMail(data)

			if err != nil {
				log.Print("Failed to sending mail")
			} else {
				log.Print(fmt.Sprintf("Mail delivered to %s", data.ToMail))
				UpdateEmailQueue(data.Id, db)
			}
		} else {
			log.Print("Fail ! Mail type unknown")
		}
	}

}

func UpdateEmailQueue(id uint, db *gorm.DB) {

	tx := db.Begin()

	err := tx.Model(&EmailQueueModel{}).Where("id = ?", id).Update("status_mail", true).Error

	if err != nil {
		tx.Rollback()
		log.Print("Failed to upload data")
	} else {
		tx.Commit()
	}
}
