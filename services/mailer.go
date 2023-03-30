package services

import (
	"email-serving-datathon/models"
	"fmt"
	"gorm.io/gorm"
)

func CreateMail(db *gorm.DB, mail models.Mail) error {
	return db.Create(&mail).Error
}

func FindMail(db *gorm.DB, userId uint, email string) (models.Mail, error) {
	var mail = models.Mail{}
	err := db.Where("email = ? AND user_id = ?", email, userId).Find(&mail).Error
	//err := db.Limit(1).Find(&mail).Error
	if err != nil {
		return models.Mail{}, err
	}
	return mail, nil
}
func UpdateMail(db *gorm.DB, mail models.Mail) error {
	return db.Model(models.Mail{
		Model: gorm.Model{ID: mail.ID},
	}).Updates(mail).Error
}

func SaveEmail(db *gorm.DB, msg models.Mail) error {
	mail, err := FindMail(db, msg.UserId, msg.Email)
	if err != nil {
		return err
	}
	fmt.Println("Printing mail")
	fmt.Print(mail)
	if mail != (models.Mail{}) {
		err = UpdateMail(db, mail)
	} else {
		return CreateMail(db, models.Mail{
			Model:    gorm.Model{},
			UserId:   msg.UserId,
			TeamCode: msg.TeamCode,
			Email:    msg.Email,
			Status:   models.Success,
		})
	}
	return nil
}
