package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

func SetupDB(cfg *Config) *gorm.DB {
	DBString := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	URL := fmt.Sprintf(DBString, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(cfg.DBType, URL)
	if err != nil {
		logrus.Fatal("cannot connect to the database", err)
	}
	DB = db
	return db
}
