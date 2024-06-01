package database

import (
	"ds/internal/config"
	"ds/internal/database/models"

	mysql_driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	cfg := mysql_driver.Config{
		DBName:               config.Config(config.DB_NAME_KEY),
		User:                 config.Config(config.DB_USER_KEY),
		Passwd:               config.Config(config.DB_PASSWORD_KEY),
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	d := mysql.Open(cfg.FormatDSN())
	return gorm.Open(d, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.HistoryRecord{}, &models.HistoryRecordGeoJSON{})
}
