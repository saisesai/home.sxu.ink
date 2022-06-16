package model

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.WithError(err).Fatalln("failed to connect to db!")
	}
	err = db.AutoMigrate(&User{}, &Relay{})
	if err != nil {
		log.WithError(err).Fatalln("failed to migrate db!")
	}
}
