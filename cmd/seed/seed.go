package main

import (
	"os"

	"github.com/night1010/everhealth/logger"
	"github.com/night1010/everhealth/migration"
	"github.com/night1010/everhealth/repository"
)

func main() {
	logger.SetLogrusLogger()

	_ = os.Setenv("APP_ENV", "debug")

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	migration.Seed(db)
}
