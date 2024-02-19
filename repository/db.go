package repository

import (
	"fmt"
	"time"

	"github.com/night1010/everhealth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDsn() string {
	dbConfig := config.NewDbConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Pass, dbConfig.Name)
	return dsn

}

func GetConnection() (*gorm.DB, error) {
	dsn := getDsn()
	mode := logger.Silent
	appConfig := config.NewAppConfig()
	dbConfig := config.NewDbConfig()

	if appConfig.IsInDevMode() {
		mode = logger.Error
	}

	if appConfig.IsInDebugMode() {
		mode = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(mode),
	})
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	sqlDb.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db, nil
}
