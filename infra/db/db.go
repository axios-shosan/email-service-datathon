package database

import (
	"email-serving-datathon/models"
	"email-serving-datathon/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type name struct {
}

func Connect() *gorm.DB {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_ADDRESS") + ")" + "/datathon_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	utils.Panic(err, "Failed to connect to database!!!")
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Mail{})
}
